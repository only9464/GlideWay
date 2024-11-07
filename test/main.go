package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lcvvvv/gonmap"
)

// 定义扫描结果结构
type ScanResult struct {
	Port     int
	Status   gonmap.Status
	Response *gonmap.Response
}

// 扫描单个端口
func scanPort(ip string, port int, timeout time.Duration) (bool, error) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}

// 扫描IP的所有端口并进行指纹识别
func scanIP(ip string) []ScanResult {
	var results []ScanResult
	var wg sync.WaitGroup
	resultChan := make(chan ScanResult)

	// 创建gonmap实例
	scanner := gonmap.New()
	scanner.SetTimeout(2 * time.Second)

	// 设置并发数
	maxConcurrent := 200
	guard := make(chan struct{}, maxConcurrent)

	// 扫描端口范围(1-65535)
	go func() {
		for port := 1; port <= 65535; port++ {
			wg.Add(1)
			guard <- struct{}{} // 限制并发

			go func(port int) {
				defer wg.Done()
				defer func() { <-guard }()

				// 先进行端口开放性检测
				isOpen, err := scanPort(ip, port, 1*time.Second)
				if err != nil || !isOpen {
					return
				}

				// 对开放端口进行指纹识别
				status, response := scanner.ScanTimeout(ip, port, 2*time.Second)

				resultChan <- ScanResult{
					Port:     port,
					Status:   status,
					Response: response,
				}
			}(port)
		}
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func main() {
	ip := "127.0.0.1" // 要扫描的IP地址

	fmt.Printf("开始扫描 %s ...\n", ip)
	startTime := time.Now()

	results := scanIP(ip)

	// 输出结果
	fmt.Printf("\n扫描完成! 用时: %v\n", time.Since(startTime))
	fmt.Printf("发现开放端口数量: %d\n\n", len(results))

	for _, result := range results {
		fmt.Printf("端口: %d\n", result.Port)
		fmt.Printf("状态: %s\n", result.Status)

		if result.Status == gonmap.Matched && result.Response != nil {
			fp := result.Response.FingerPrint

			// 基本信息
			fmt.Printf("服务: %s\n", fp.Service)
			if fp.ProductName != "" {
				fmt.Printf("产品名称: %s\n", fp.ProductName)
			}
			if fp.Version != "" {
				fmt.Printf("版本: %s\n", fp.Version)
			}

			// 额外信息
			if fp.Info != "" {
				fmt.Printf("详细信息: %s\n", fp.Info)
			}
			if fp.Hostname != "" {
				fmt.Printf("主机名: %s\n", fp.Hostname)
			}
			if fp.OperatingSystem != "" {
				fmt.Printf("操作系统: %s\n", fp.OperatingSystem)
			}
			if fp.DeviceType != "" {
				fmt.Printf("设备类型: %s\n", fp.DeviceType)
			}

			// 调试信息
			if fp.ProbeName != "" {
				fmt.Printf("探针名称: %s\n", fp.ProbeName)
			}
			if fp.MatchRegexString != "" {
				fmt.Printf("匹配规则: %s\n", fp.MatchRegexString)
			}

			// TLS信息
			if result.Response.TLS {
				fmt.Printf("TLS: 是\n")
			}
		}
		fmt.Println("------------------------")
	}
}
