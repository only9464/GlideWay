package gitdorker

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

func (s *App) GitdorkerCalculate(a, b int) int {
	return GitdorkerCalculate(a, b)
}
func (a *App) Gitdorker(mainKeyword string, subKeyword string, token string) *GithubResult {
	ctx := context.Background()
	return SearchGithub(ctx, mainKeyword, subKeyword, token)
}
