package main

import (
	"GlideWay/apps/scanner"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ScanPorts(IP string) error {
	// 在新的 goroutine 中执行扫描
	go func() {
		scanner.ScanPorts(IP, func(portInfo scanner.PortInfo) {
			// 发送端口信息到前端
			runtime.EventsEmit(a.ctx, "port-found", portInfo)
		})
		// 扫描完成后发送完成事件
		runtime.EventsEmit(a.ctx, "scan-complete", nil)
	}()

	return nil
}
