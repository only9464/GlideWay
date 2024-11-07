package scanner

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type PortInfo struct {
	Port     int    `json:"port"`
	Service  string `json:"service"`
	Banner   string `json:"banner"`
	Protocol string `json:"protocol"`
}

// probeService 尝试识别服务，返回服务名称、banner信息和协议类型
func probeService(address string, port int) (service string, banner string, protocol string) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), time.Second*3)
	if err != nil {
		return "unknown", "", "unknown"
	}
	defer conn.Close()

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// 对于HTTP服务，发送HTTP请求
	if port == 80 || port == 443 || port == 8080 || port == 8443 {
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	}

	// 读取响应
	reader := bufio.NewReader(conn)
	banner, err = reader.ReadString('\n')
	if err != nil {
		return "unknown", "", "unknown"
	}

	banner = strings.TrimSpace(banner)
	protocol = "unknown"
	service = "unknown"

	// 基于 banner 识别服务
	banner = strings.ToLower(banner)
	switch {
	case strings.Contains(banner, "ssh"):
		service = "SSH"
		protocol = "SSH"
	case strings.Contains(banner, "ftp"):
		service = "FTP"
		protocol = "FTP"
	case strings.Contains(banner, "http"):
		service = "HTTP"
		protocol = "HTTP/1.1"
	case strings.Contains(banner, "smtp"):
		service = "SMTP"
		protocol = "SMTP"
	case strings.Contains(banner, "mysql"):
		service = "MySQL"
		protocol = "MySQL"
	case strings.Contains(banner, "redis"):
		service = "Redis"
		protocol = "Redis"
	case strings.Contains(banner, "mongodb"):
		service = "MongoDB"
		protocol = "MongoDB"
	case strings.Contains(banner, "postgresql"):
		service = "PostgreSQL"
		protocol = "PostgreSQL"
		// ... 可以添加更多服务识别规则
	}

	return service, banner, protocol
}

// 添加一个回调函数类型
type PortCallback func(PortInfo)

func ScanPorts(target string, callback PortCallback) {
	var wg sync.WaitGroup
	const maxConcurrentGoroutines = 500
	semaphore := make(chan struct{}, maxConcurrentGoroutines)

	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			address := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", address, time.Second)

			if err == nil {
				conn.Close()

				// 尝试识别服务
				service, banner, protocol := probeService(target, p)

				portInfo := PortInfo{
					Port:     p,
					Service:  service,
					Banner:   banner,
					Protocol: protocol,
				}

				// 直接回调，不需要互斥锁
				callback(portInfo)
			}
		}(port)
	}

	wg.Wait()
}
