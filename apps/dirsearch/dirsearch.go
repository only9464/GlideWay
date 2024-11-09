package dirsearch

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OJ/gobuster/v3/gobusterdir"
	"github.com/OJ/gobuster/v3/libgobuster"
)

var (
	actualScanned int32
)

func ScanDir(ctx context.Context, target string, dictPath string, maxThreads int, pathCallback PathCallback, progressCallback ProgressCallback) error {
	atomic.StoreInt32(&actualScanned, 0)

	content, err := os.ReadFile(dictPath)
	if err != nil {
		return fmt.Errorf("读取字典文件失败: %w", err)
	}

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

	progressCallback(0, totalPaths)

	// 自动计算插件实例数和线程数以更接近maxThreads
	numPlugins := maxThreads / 10
	if numPlugins < 1 {
		numPlugins = 1
	}
	threadsPerPlugin := maxThreads / numPlugins
	if threadsPerPlugin < 1 {
		threadsPerPlugin = 1
	}

	globalopts := &libgobuster.Options{
		Threads: threadsPerPlugin,
	}

	opts := gobusterdir.NewOptionsDir()
	opts.URL = target
	opts.NoTLSValidation = true
	opts.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	opts.StatusCodes = "200,201,202,203,204,301,302,307,308,401,403,405"
	opts.StatusCodesParsed.AddRange([]int{200, 201, 202, 203, 204, 301, 302, 307, 308, 401, 403, 405})
	opts.Timeout = time.Second * 10

	bufferSize := maxThreads * 20
	results := make(chan libgobuster.Result, bufferSize)
	errorChan := make(chan error, bufferSize)
	pathChan := make(chan string, bufferSize*2)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	var isStopped atomic.Value
	isStopped.Store(false)

	var closeOnce sync.Once

	plugins := make([]*gobusterdir.GobusterDir, numPlugins)
	for i := 0; i < numPlugins; i++ {
		plugin, err := gobusterdir.NewGobusterDir(globalopts, opts)
		if err != nil {
			return fmt.Errorf("创建扫描插件失败: %w", err)
		}
		plugins[i] = plugin
	}

	cleanup := func() {
		isStopped.Store(true)
		closeOnce.Do(func() {
			close(pathChan)
			close(results)
			close(errorChan)
		})
		wg.Wait()
	}

	for i := 0; i < numPlugins; i++ {
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

	go func() {
		defer func() {
			closeOnce.Do(func() {
				close(pathChan)
			})
		}()

		batchSize := maxThreads * 10
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

		for _, p := range batch {
			select {
			case <-ctx.Done():
				return
			case pathChan <- p:
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case <-ticker.C:
				if !isStopped.Load().(bool) {
					current := atomic.LoadInt32(&actualScanned)
					progressCallback(int(current), totalPaths)
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		if !isStopped.Load().(bool) {
			select {
			case doneChan <- struct{}{}:
				// 扫描完成时发送最后一次进度更新
				progressCallback(totalPaths, totalPaths) // 确保显示100%完成
				// 关闭所有通道
				cleanup()
			default:
			}
		}
	}()

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
				// 通道关闭时发送最后一次进度更新
				progressCallback(totalPaths, totalPaths)
				return nil
			}
			if result == nil || isStopped.Load().(bool) {
				continue
			}
			handleResult(result, pathCallback)
		}
	}
}
func processPath(ctx context.Context, plugin *gobusterdir.GobusterDir, path string, results chan libgobuster.Result, errorChan chan error) error {
	// 立即增加计数，而不是等待处理完成
	atomic.AddInt32(&actualScanned, 1)

	// 添加URL编码处理
	path = strings.ReplaceAll(path, "%", "%25") // 首先处理%符号
	path = strings.Map(func(r rune) rune {
		switch r {
		case '#', '&', '=', '+', '!', '@', '$', '^', '~':
			return -1 // 移除这些特殊字符
		default:
			return r
		}
	}, path)

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
