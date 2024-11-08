package dirsearch

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
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

// 添加全局变量
var (
	actualScanned int32 // 添加这个全局变量用于跟踪实际扫描进度
)

func ScanDir(ctx context.Context, target string, dictPath string, maxThreads int, pathCallback PathCallback, progressCallback ProgressCallback) error {
	// 重置计数器
	atomic.StoreInt32(&actualScanned, 0)

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

	// 立即通知前端总路径数
	progressCallback(0, totalPaths)

	// 创建全局选项
	globalopts := &libgobuster.Options{
		Threads: 1, // 每个工作协程使用单独的插件实例
	}

	// 创建目录扫描选项
	opts := gobusterdir.NewOptionsDir()
	opts.URL = target
	opts.NoTLSValidation = true
	opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	opts.StatusCodes = "200,201,202,203,204,301,302,307,308,401,403,405"
	opts.StatusCodesParsed.AddRange([]int{200, 201, 202, 203, 204, 301, 302, 307, 308, 401, 403, 405})
	opts.Timeout = time.Second * 10

	// 增加缓冲区大小以减少阻塞
	results := make(chan libgobuster.Result, maxThreads*20)
	errorChan := make(chan error, maxThreads*10)
	pathChan := make(chan string, maxThreads*50)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	var isStopped atomic.Value
	isStopped.Store(false)

	// 增加一个关闭标志
	var closeOnce sync.Once

	// 创建插件实例池
	plugins := make([]*gobusterdir.GobusterDir, maxThreads)
	for i := 0; i < maxThreads; i++ {
		plugin, err := gobusterdir.NewGobusterDir(globalopts, opts)
		if err != nil {
			return fmt.Errorf("创建扫描插件失败: %w", err)
		}
		plugins[i] = plugin
	}

	// 清理函数
	cleanup := func() {
		isStopped.Store(true)
		closeOnce.Do(func() {
			close(pathChan)
		})
		wg.Wait()
		closeOnce.Do(func() {
			close(results)
			close(errorChan)
		})
	}

	// 启动工作协程池
	for i := 0; i < maxThreads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			plugin := plugins[workerID]

			for path := range pathChan {
				if isStopped.Load().(bool) {
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					if err := processPath(ctx, plugin, path, results, errorChan); err != nil {
						if !strings.Contains(err.Error(), "context canceled") {
							select {
							case errorChan <- err:
							default:
							}
						}
					}
				}
			}
		}(i)
	}

	// 发送路径到工作通道
	go func() {
		defer func() {
			closeOnce.Do(func() {
				close(pathChan)
			})
		}()

		batchSize := 100
		batch := make([]string, 0, batchSize)

		for _, path := range paths {
			if isStopped.Load().(bool) {
				return
			}

			batch = append(batch, path)
			if len(batch) >= batchSize {
				for _, p := range batch {
					select {
					case <-ctx.Done():
						return
					case pathChan <- p:
					}
				}
				batch = batch[:0]
			}
		}

		// 处理剩余的路径
		for _, p := range batch {
			select {
			case <-ctx.Done():
				return
			case pathChan <- p:
			}
		}
	}()

	// 进度更新
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()

		var lastReported int32 = -1
		for {
			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case <-ticker.C:
				if !isStopped.Load().(bool) {
					current := atomic.LoadInt32(&actualScanned)
					if current != lastReported {
						progressCallback(int(current), totalPaths)
						lastReported = current
					}
				}
			}
		}
	}()

	// 等待所有工作完成
	go func() {
		wg.Wait()
		if !isStopped.Load().(bool) {
			select {
			case doneChan <- struct{}{}:
			default:
			}
		}
	}()

	// 处理结果
	for {
		select {
		case <-ctx.Done():
			cleanup()
			return context.Canceled
		case err := <-errorChan:
			if err != nil && !isStopped.Load().(bool) {
				fmt.Printf("扫描错误: %v\n", err)
			}
		case result, ok := <-results:
			if !ok {
				return nil
			}
			if result == nil || isStopped.Load().(bool) {
				continue
			}
			handleResult(result, pathCallback)
		}
	}
}

// 优化处理单个路径的函数
func processPath(ctx context.Context, plugin *gobusterdir.GobusterDir, path string, results chan libgobuster.Result, errorChan chan error) error {
	defer atomic.AddInt32(&actualScanned, 1) // 确保在处理完路径后增加计数

	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		err := plugin.ProcessWord(ctx, path, &libgobuster.Progress{
			ResultChan: results,
			ErrorChan:  errorChan,
		})
		return err
	}
}

// 处理结果
func handleResult(result libgobuster.Result, pathCallback PathCallback) {
	found := result.(gobusterdir.Result)
	pathInfo := PathInfo{
		URL:           found.URL,
		Path:          found.Path,
		StatusCode:    found.StatusCode,
		ContentType:   found.Header.Get("Content-Type"),
		ContentLength: found.Size,
		Header:        found.Header,
	}
	pathCallback(pathInfo)
}
