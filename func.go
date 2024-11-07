package main

import (
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
	mutex      sync.Mutex
	totalPorts int32
	scanned    int32
}

var (
	currentScan *scanControl
	scanMutex   sync.Mutex
)

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
				runtime.EventsEmit(a.ctx, "port-found", portInfo)
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

type ScanProgress struct {
	CurrentPort int32  `json:"current_port"`
	TotalPorts  int32  `json:"total_ports"`
	Status      string `json:"status"`
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
