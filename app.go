package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// /////////////////////////////////////////////////////////////////////////////
// ScanPorts 扫描端口并返回开放的端口列表
func (a *App) ScanPorts(target string) []int {
	var openPorts []int
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 扫描常用端口范围 1-1024
	for port := 1; port <= 1024; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			address := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", address, time.Second)

			if err == nil {
				conn.Close()
				mutex.Lock()
				openPorts = append(openPorts, p)
				mutex.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return openPorts
}
