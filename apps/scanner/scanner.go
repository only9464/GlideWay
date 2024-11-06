package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanPorts 扫描端口并返回开放的端口列表，直接返回 []interface{}
func ScanPorts(target string) []interface{} {
	var openPorts []interface{}
	var mutex sync.Mutex
	var wg sync.WaitGroup

	const maxConcurrentGoroutines = 1000
	semaphore := make(chan struct{}, maxConcurrentGoroutines)

	// 扫描端口范围 1-65535
	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		semaphore <- struct{}{} // 向信号量发送请求

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 工作完成后释放信号量

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
