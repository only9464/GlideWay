package util

import (
	"GlideWay/apps/scanner"
	"log"
)

// scanner 包初始化
func ScannerInit() {
	// 自动注册 ScanPorts 函数
	RegisterFunction("ScanPorts", scanner.ScanPorts)
	log.Println("scanner包初始化完成")
}

// init 函数会在包被导入时自动执行，用来注册函数

func init() {
	ScannerInit()
	// 逐个输出已注册的函数
	for name := range functionRegistry {
		log.Printf("已注册函数: %s", name)
	}
}
