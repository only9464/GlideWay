package main

import (
	"GlideWay/apps/dirsearch"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 目录扫描相关结构体和变量
type DirsearchProgress struct {
	Current int     `json:"current"`
	Total   int     `json:"total"`
	Speed   float64 `json:"speed"` // 添加速度字段
}

type PathResult struct {
	Path          string `json:"path"`
	FullUrl       string `json:"fullUrl"`
	StatusCode    int    `json:"statusCode"`
	ContentType   string `json:"contentType"`
	ContentLength int64  `json:"contentLength"`
}

type DirsearchControl struct {
	cancel     context.CancelFunc
	scanned    int32
	totalPaths int32
}

var (
	currentDirsearch *DirsearchControl
	dirsearchMutex   sync.Mutex
)

// StartDirsearch 启动目录扫描
func (a *App) StartDirsearch(target string, dictPath string, maxThreads int) error {
	if a == nil || a.ctx == nil {
		return fmt.Errorf("app context is not initialized")
	}

	dirsearchMutex.Lock()
	defer dirsearchMutex.Unlock()

	if currentDirsearch != nil {
		return fmt.Errorf("dirsearch is already running")
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

		err := dirsearch.ScanDir(
			ctx,
			target,
			dictPath,
			maxThreads,
			// 路径发现回调
			func(pathInfo dirsearch.PathInfo) {
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

// GetDirsearchStatus 获取目录扫描状态
func (a *App) GetDirsearchStatus() string {
	dirsearchMutex.Lock()
	defer dirsearchMutex.Unlock()

	if currentDirsearch != nil {
		return "running"
	}
	return "idle"
}

// GetDirsearchProgress 获取目录扫描进度
func (a *App) GetDirsearchProgress() DirsearchProgress {
	dirsearchMutex.Lock()
	defer dirsearchMutex.Unlock()

	if currentDirsearch == nil {
		return DirsearchProgress{
			Current: 0,
			Total:   0,
			Speed:   0,
		}
	}

	return DirsearchProgress{
		Current: int(atomic.LoadInt32(&currentDirsearch.scanned)),
		Total:   int(atomic.LoadInt32(&currentDirsearch.totalPaths)),
		Speed:   0, // 这里可以添加实时速度计算，但需要维护额外的状态
	}
}
