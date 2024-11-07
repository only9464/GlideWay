package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lcvvvv/gonmap"
)

type ScanConfig struct {
	Target     string
	StartPort  int
	EndPort    int
	MaxThreads int
	Timeout    time.Duration
}

type PortInfo struct {
	Port            int    `json:"port"`
	Protocol        string `json:"protocol"`
	Service         string `json:"service"`
	ProductName     string `json:"product_name"`
	Version         string `json:"version"`
	Info            string `json:"info"`
	Hostname        string `json:"hostname"`
	OperatingSystem string `json:"operating_system"`
	DeviceType      string `json:"device_type"`
	ProbeName       string `json:"probe_name"`
	TLS             bool   `json:"tls"`
}

type PortCallback func(PortInfo)

func ScanPortsCombined(ctx context.Context, config ScanConfig, callback PortCallback) error {
	if callback == nil {
		return fmt.Errorf("callback function cannot be nil")
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxThreads)
	var scanned int32

	// 创建gonmap实例
	scanner := gonmap.New()
	scanner.SetTimeout(config.Timeout)

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

					// 对开放端口进行指纹识别
					status, response := scanner.ScanTimeout(config.Target, p, config.Timeout)

					portInfo := PortInfo{
						Port:     p,
						Protocol: "tcp",
					}

					if status == gonmap.Matched && response != nil {
						fp := response.FingerPrint
						portInfo.Service = fp.Service
						portInfo.ProductName = fp.ProductName
						portInfo.Version = fp.Version
						portInfo.Info = fp.Info
						portInfo.Hostname = fp.Hostname
						portInfo.OperatingSystem = fp.OperatingSystem
						portInfo.DeviceType = fp.DeviceType
						portInfo.ProbeName = fp.ProbeName
						portInfo.TLS = response.TLS
					}

					select {
					case <-ctx.Done():
						return
					default:
						callback(portInfo)
					}
				}
			}(port)
		}
	}

	wg.Wait()
	return nil
}
