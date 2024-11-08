package main

import (
	"GlideWay/apps/dirsearch"
	"GlideWay/apps/scanner"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type scanControl struct {
	cancel     context.CancelFunc
	totalPorts int32
	scanned    int32
}

var (
	currentScan *scanControl
	scanMutex   sync.Mutex
)

type ScanProgress struct {
	CurrentPort int32  `json:"current_port"`
	TotalPorts  int32  `json:"total_ports"`
	Status      string `json:"status"`
}

func (a *App) ScanPorts(IP string, startPort int, endPort int, maxThreads int) error {
	if a == nil || a.ctx == nil {
		return fmt.Errorf("app context is not initialized")
	}

	scanMutex.Lock()
	defer scanMutex.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	totalPorts := int32(endPort - startPort + 1)

	// 创建新的 scanControl
	newScan := &scanControl{
		cancel:     cancel,
		totalPorts: totalPorts,
		scanned:    0,
	}

	// 原子性地替换 currentScan
	currentScan = newScan

	config := scanner.ScanConfig{
		Target:     IP,
		StartPort:  startPort,
		EndPort:    endPort,
		MaxThreads: maxThreads,
		Timeout:    time.Second * 2,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				runtime.EventsEmit(a.ctx, "scan-error", "Internal error occurred")
			}
			scanMutex.Lock()
			currentScan = nil
			scanMutex.Unlock()
			runtime.EventsEmit(a.ctx, "scan-status", "idle")
		}()

		// 发送初始状态
		runtime.EventsEmit(a.ctx, "scan-status", "running")
		runtime.EventsEmit(a.ctx, "scan-progress", map[string]interface{}{
			"current_port": startPort,
			"total_ports":  totalPorts,
			"status":       "scanning",
		})

		err := scanner.ScanPortsCombined(ctx, config, func(portInfo scanner.PortInfo) {
			scanMutex.Lock()
			if currentScan == nil {
				scanMutex.Unlock()
				return
			}
			scanMutex.Unlock()

			if portInfo.Protocol == "progress" {
				scanned := atomic.AddInt32(&currentScan.scanned, 1)
				runtime.EventsEmit(a.ctx, "scan-progress", map[string]interface{}{
					"current_port": portInfo.Port,
					"total_ports":  totalPorts,
					"scanned":      scanned,
					"status":       "scanning",
				})
			} else {
				// 发送完整的端口信息，包括指纹识别结果
				runtime.EventsEmit(a.ctx, "port-found", map[string]interface{}{
					"port":             portInfo.Port,
					"protocol":         portInfo.Protocol,
					"service":          portInfo.Service,
					"product_name":     portInfo.ProductName,
					"version":          portInfo.Version,
					"info":             portInfo.Info,
					"hostname":         portInfo.Hostname,
					"operating_system": portInfo.OperatingSystem,
					"device_type":      portInfo.DeviceType,
					"probe_name":       portInfo.ProbeName,
					"tls":              portInfo.TLS,
				})
			}
		})

		scanMutex.Lock()
		defer scanMutex.Unlock()

		if currentScan == nil {
			return
		}

		if err != nil {
			if err == context.Canceled {
				runtime.EventsEmit(a.ctx, "scan-status", "cancelled")
				runtime.EventsEmit(a.ctx, "scan-progress", map[string]interface{}{
					"current_port": atomic.LoadInt32(&currentScan.scanned),
					"total_ports":  totalPorts,
					"status":       "cancelled",
				})
			} else {
				runtime.EventsEmit(a.ctx, "scan-error", err.Error())
				runtime.EventsEmit(a.ctx, "scan-status", "error")
				runtime.EventsEmit(a.ctx, "scan-progress", map[string]interface{}{
					"current_port": atomic.LoadInt32(&currentScan.scanned),
					"total_ports":  totalPorts,
					"status":       "error",
				})
			}
		} else {
			runtime.EventsEmit(a.ctx, "scan-complete", map[string]interface{}{
				"total_ports": totalPorts,
				"scanned":     atomic.LoadInt32(&currentScan.scanned),
			})
			runtime.EventsEmit(a.ctx, "scan-status", "completed")
			runtime.EventsEmit(a.ctx, "scan-progress", map[string]interface{}{
				"current_port": endPort,
				"total_ports":  totalPorts,
				"status":       "completed",
			})
		}
	}()

	return nil
}

func (a *App) StopScan() error {
	scanMutex.Lock()
	defer scanMutex.Unlock()

	if currentScan != nil && currentScan.cancel != nil {
		currentScan.cancel()
		runtime.EventsEmit(a.ctx, "scan-status", "stopping")
		runtime.EventsEmit(a.ctx, "scan-progress", map[string]interface{}{
			"current_port": atomic.LoadInt32(&currentScan.scanned),
			"total_ports":  currentScan.totalPorts,
			"status":       "stopping",
		})
	}
	return nil
}

func (a *App) GetScanStatus() string {
	scanMutex.Lock()
	defer scanMutex.Unlock()

	if currentScan != nil {
		return "running"
	}
	return "idle"
}

func (a *App) GetScanProgress() ScanProgress {
	scanMutex.Lock()
	defer scanMutex.Unlock()

	if currentScan == nil {
		return ScanProgress{
			Status: "idle",
		}
	}

	return ScanProgress{
		CurrentPort: atomic.LoadInt32(&currentScan.scanned),
		TotalPorts:  currentScan.totalPorts,
		Status:      "running",
	}
}

// ... 前面的代码保持不变 ...

// 目录扫描相关结构体和变量
type DirsearchProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
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
				atomic.StoreInt32(&currentDirsearch.scanned, int32(current))
				atomic.StoreInt32(&currentDirsearch.totalPaths, int32(total))
				dirsearchMutex.Unlock()

				progress := DirsearchProgress{
					Current: current,
					Total:   total,
				}
				runtime.EventsEmit(a.ctx, "dirsearch-progress", progress) // 这里发送进度事件
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
		}
	}

	return DirsearchProgress{
		Current: int(atomic.LoadInt32(&currentDirsearch.scanned)),
		Total:   int(atomic.LoadInt32(&currentDirsearch.totalPaths)),
	}
}
