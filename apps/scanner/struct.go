package scanner

import (
	"context"
	"sync"
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
