package main

import (
	"GlideWay/util"
	"context"
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

// ExecuteFunction 调用动态注册的函数并返回结果
func (a *App) ExecuteFunction(funcName string, params []interface{}) ([]interface{}, error) {
	result, err := util.CallFunction(funcName, params...)
	if err != nil {
		return nil, err
	}

	// 检查 result 是否为切片类型
	if results, ok := result.([]interface{}); ok {
		return results, nil
	}

	// 如果不是切片，返回单个结果
	return []interface{}{result}, nil
}
