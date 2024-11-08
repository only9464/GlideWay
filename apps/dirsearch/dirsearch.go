package dirsearch

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/OJ/gobuster/v3/gobusterdir"
	"github.com/OJ/gobuster/v3/libgobuster"
)

type PathInfo struct {
	URL           string      `json:"url"`
	Path          string      `json:"path"`
	StatusCode    int         `json:"statusCode"`
	ContentType   string      `json:"contentType"`
	ContentLength int64       `json:"contentLength"`
	Header        http.Header `json:"header"`
}

type PathCallback func(PathInfo)
type ProgressCallback func(current, total int)

func ScanDir(ctx context.Context, target string, dictPath string, maxThreads int, pathCallback PathCallback, progressCallback ProgressCallback) error {
	// 读取字典文件
	content, err := os.ReadFile(dictPath)
	if err != nil {
		return fmt.Errorf("读取字典文件失败: %w", err)
	}

	// 按行分割并过滤空行
	paths := make([]string, 0)
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			paths = append(paths, line)
		}
	}

	totalPaths := len(paths)
	if totalPaths == 0 {
		return fmt.Errorf("字典文件为空")
	}

	// 创建全局选项
	globalopts := &libgobuster.Options{
		Threads: maxThreads,
	}

	// 创建目录扫描选项
	opts := gobusterdir.NewOptionsDir()
	opts.URL = target
	opts.NoTLSValidation = true
	opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	opts.StatusCodes = "200,201,202,203,204,301,302,307,308,401,403,405"
	opts.StatusCodesParsed.AddRange([]int{200, 201, 202, 203, 204, 301, 302, 307, 308, 401, 403, 405})
	opts.Timeout = time.Second * 10 // 增加超时时间到10秒

	// 创建目录扫描插件
	plugin, err := gobusterdir.NewGobusterDir(globalopts, opts)
	if err != nil {
		return fmt.Errorf("创建扫描插件失败: %w", err)
	}

	// 创建所有通道
	results := make(chan libgobuster.Result)
	errorChan := make(chan error)
	progressChan := make(chan struct{})
	pathChan := make(chan string, maxThreads)

	var wg sync.WaitGroup
	var cleanupOnce sync.Once

	// 清理函数，用于关闭所有通道
	cleanup := func() {
		cleanupOnce.Do(func() {
			close(pathChan)
			close(progressChan)
			close(results)
			close(errorChan)
		})
	}

	// 确保在函数返回时清理资源
	defer cleanup()

	// 启动工作协程
	for i := 0; i < maxThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				select {
				case <-ctx.Done():
					return
				default:
					// 添加重试逻辑
					maxRetries := 3
					for retry := 0; retry < maxRetries; retry++ {
						err := plugin.ProcessWord(ctx, path, &libgobuster.Progress{
							ResultChan: results,
							ErrorChan:  errorChan,
						})
						if err == nil {
							break
						}
						if retry < maxRetries-1 {
							time.Sleep(time.Second * time.Duration(retry+1))
						}
					}
					time.Sleep(time.Millisecond * 100)
				}
			}
		}()
	}

	// 启动进度跟踪协程
	go func() {
		current := 0
		for range paths {
			select {
			case <-ctx.Done():
				return
			case progressChan <- struct{}{}:
				current++
				select {
				case <-ctx.Done():
					return
				default:
					progressCallback(current, totalPaths) // 这里调用进度回调
				}
			}
		}
	}()

	// 发送路径到工作通道
	go func() {
		for _, path := range paths {
			select {
			case <-ctx.Done():
				return
			case pathChan <- path:
			}
		}
	}()

	// 等待所有工作完成
	go func() {
		wg.Wait()
		cleanup()
	}()

	// 处理结果
	for {
		select {
		case <-ctx.Done():
			fmt.Println("扫描已取消")
			return context.Canceled
		case err := <-errorChan:
			if err != nil {
				fmt.Printf("扫描错误: %v\n", err)
			}
			continue
		case result, ok := <-results:
			if !ok {
				return nil
			}
			if result == nil {
				continue
			}

			found := result.(gobusterdir.Result)
			pathInfo := PathInfo{
				URL:           found.URL,
				Path:          found.Path,
				StatusCode:    found.StatusCode,
				ContentType:   found.Header.Get("Content-Type"),
				ContentLength: found.Size,
				Header:        found.Header,
			}

			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				pathCallback(pathInfo)
			}
		}
	}
}
