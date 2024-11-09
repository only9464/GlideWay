package main

import (
	"GlideWay/apps/dirsearch"
	"GlideWay/apps/gitdorker"
	"GlideWay/apps/scanner"

	"context"
	"embed"

	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	PorotsScanner := scanner.NewApp()
	DirsearchApp := dirsearch.NewApp()
	GitdorkerApp := gitdorker.NewApp()
	// Create application with options
	err := wails.Run(&options.App{
		Title:            "GlideWay",
		Width:            1024,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			PorotsScanner.Startup(ctx)
			DirsearchApp.Startup(ctx)
			GitdorkerApp.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			PorotsScanner,
			DirsearchApp,
			GitdorkerApp,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			Theme:                0,
			BackdropType:         windows.Acrylic, // 使用亚克力效果
		},
		WindowStartState: options.Maximised,
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
