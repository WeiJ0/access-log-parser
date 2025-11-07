package main

import (
	"embed"
	
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"access-log-analyzer/internal/app"
	"access-log-analyzer/pkg/logger"
)

//go:embed all:frontend/dist
var assets embed.FS

// main 是應用程式的入口點
// 初始化 Wails runtime 並啟動應用程式
func main() {
	// 初始化全域 logger
	logger.Init()
	log := logger.Get()
	
	log.Info().Msg("Apache Log Analyzer 正在啟動...")

	// 建立應用程式實例
	appInstance := app.NewApp()

	// 建立應用程式配置
	err := wails.Run(&options.App{
		Title:  "Apache Log 工具",
		Width:  1280,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        appInstance.Startup,
		OnDomReady:       appInstance.DomReady,
		OnShutdown:       appInstance.Shutdown,
		Bind: []interface{}{
			appInstance,
		},
		// 啟用開發者工具（方便除錯）
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		// Windows 特定配置
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal().Err(err).Msg("應用程式啟動失敗")
	}
}
