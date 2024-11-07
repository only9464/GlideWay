package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ScanConfig struct {
	Target     string
	StartPort  int
	EndPort    int
	MaxThreads int
	Timeout    time.Duration
}

type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type PortCallback func(PortInfo)

func ScanPortsCombined(ctx context.Context, config ScanConfig, callback PortCallback) error {
	if callback == nil {
		return fmt.Errorf("callback function cannot be nil")
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxThreads)
	var scanned int32

	for port := config.StartPort; port <= config.EndPort; port++ {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			wg.Add(1)
			semaphore <- struct{}{}

			go func(p int) {
				defer func() {
					wg.Done()
					<-semaphore
					// 恢复任何可能的panic
					if r := recover(); r != nil {
						fmt.Printf("Recovered from panic in port scan goroutine: %v\n", r)
					}
				}()

				// 更新进度
				atomic.AddInt32(&scanned, 1)
				select {
				case <-ctx.Done():
					return
				default:
					callback(PortInfo{
						Port:     p,
						Protocol: "progress",
					})
				}

				address := fmt.Sprintf("%s:%d", config.Target, p)
				conn, err := net.DialTimeout("tcp", address, config.Timeout)

				if err == nil && conn != nil {
					conn.Close()
					select {
					case <-ctx.Done():
						return
					default:
						callback(PortInfo{
							Port:     p,
							Protocol: "tcp",
						})
					}
				}
			}(port)
		}
	}

	wg.Wait()
	return nil
}
