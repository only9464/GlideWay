package scanner

import "context"

type App struct {
	ctx context.Context
}

// NewApp 创建新的 App 实例
func NewApp() *App {
	return &App{}
}

// Startup 在应用启动时初始化上下文
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}
