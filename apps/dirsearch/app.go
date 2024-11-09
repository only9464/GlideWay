package dirsearch

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

// NewApp 创建新的 App 实例
func NewApp() *App {
	return &App{}
}

// Startup 在应用启动时初始化上下文
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

var (
	currentDirsearch *DirsearchControl
	dirsearchMutex   sync.Mutex
)

// OpenFileDialog 打开文件选择对话框
func (a *App) OpenFileDialog() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "选择字典文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文本文件 (*.txt)",
				Pattern:     "*.txt",
			},
		},
	}

	return runtime.OpenFileDialog(a.ctx, options)
}

// StartDirsearch 启动目录扫描
func (a *App) StartDirsearch(target string, dictPath string, maxThreads int) error {
	if a == nil || a.ctx == nil {
		return fmt.Errorf("app context is not initialized")
	}

	dirsearchMutex.Lock()
	defer dirsearchMutex.Unlock()

	if currentDirsearch != nil {
		currentDirsearch = nil
		// return fmt.Errorf("dirsearch is already running")
	}

	// 创建新的上下文和控制器
	ctx, cancel := context.WithCancel(context.Background())
	currentDirsearch = &DirsearchControl{
		cancel:     cancel,
		scanned:    0,
		totalPaths: 0,
	}

	var (
		lastScanned   int32
		lastTimestamp = time.Now()
	)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("扫描发生panic: %v\n", r)
				runtime.EventsEmit(a.ctx, "dirsearch-status", "error")
			}
			dirsearchMutex.Lock()
			currentDirsearch = nil
			dirsearchMutex.Unlock()
			runtime.EventsEmit(a.ctx, "dirsearch-status", "idle")
		}()

		fmt.Printf("开始扫描: target=%s, dictPath=%s, maxThreads=%d\n", target, dictPath, maxThreads)

		err := ScanDir(
			ctx,
			target,
			dictPath,
			maxThreads,
			// 路径发现回调
			func(pathInfo PathInfo) {
				dirsearchMutex.Lock()
				if currentDirsearch == nil {
					dirsearchMutex.Unlock()
					return
				}
				dirsearchMutex.Unlock()

				result := PathResult{
					Path:          pathInfo.Path,
					FullUrl:       pathInfo.URL + pathInfo.Path,
					StatusCode:    pathInfo.StatusCode,
					ContentType:   pathInfo.ContentType,
					ContentLength: pathInfo.ContentLength,
				}
				runtime.EventsEmit(a.ctx, "path-found", result)
			},
			// 进度更新回调
			func(current, total int) {
				dirsearchMutex.Lock()
				if currentDirsearch == nil {
					dirsearchMutex.Unlock()
					return
				}

				// 计算扫描速度
				now := time.Now()
				speed := float64(current-int(lastScanned)) / now.Sub(lastTimestamp).Seconds()
				lastScanned = int32(current)
				lastTimestamp = now

				atomic.StoreInt32(&currentDirsearch.scanned, int32(current))
				atomic.StoreInt32(&currentDirsearch.totalPaths, int32(total))
				dirsearchMutex.Unlock()

				progress := DirsearchProgress{
					Current: current,
					Total:   total,
					Speed:   speed,
				}
				runtime.EventsEmit(a.ctx, "dirsearch-progress", progress)
			},
		)

		if err != nil {
			fmt.Printf("扫描出错: %v\n", err)
			if err == context.Canceled {
				runtime.EventsEmit(a.ctx, "dirsearch-status", "cancelled")
			} else {
				runtime.EventsEmit(a.ctx, "dirsearch-status", "error")
				runtime.EventsEmit(a.ctx, "dirsearch-error", err.Error())
			}
			return
		}

		fmt.Println("扫描完成")
		runtime.EventsEmit(a.ctx, "dirsearch-status", "completed")
	}()

	return nil
}

// StopDirsearch 停止目录扫描
func (a *App) StopDirsearch() error {
	dirsearchMutex.Lock()
	defer dirsearchMutex.Unlock()

	if currentDirsearch != nil && currentDirsearch.cancel != nil {
		// 调用 context 的 cancel 函数停止扫描
		currentDirsearch.cancel()

		// 发送状态更新事件
		runtime.EventsEmit(a.ctx, "dirsearch-status", "stopping")

		// 发送最终进度事件
		runtime.EventsEmit(a.ctx, "dirsearch-progress", DirsearchProgress{
			Current: int(atomic.LoadInt32(&currentDirsearch.scanned)),
			Total:   int(atomic.LoadInt32(&currentDirsearch.totalPaths)),
			Speed:   0, // 停止时速度为0
		})

		fmt.Println("正在停止目录扫描...")
		return nil
	}

	return fmt.Errorf("no dirsearch is running")
}
