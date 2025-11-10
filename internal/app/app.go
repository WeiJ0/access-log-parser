package app

import (
	"access-log-analyzer/pkg/logger"
	"context"
	"fmt"
)

// App 結構表示主應用程式
// 包含應用程式狀態和 Wails runtime 上下文
type App struct {
	ctx   context.Context
	state *State
	log   *logger.Logger
}

// NewApp 建立新的 App 實例
// 初始化應用程式狀態和日誌記錄器
func NewApp() *App {
	return &App{
		state: NewState(),
		log:   logger.Get(),
	}
}

// Startup 在應用程式啟動時調用
// 初始化 Wails runtime 上下文和應用程式狀態
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.log.Info().
		Bool("ctx_nil", ctx == nil).
		Str("ctx_type", fmt.Sprintf("%T", ctx)).
		Msg("應用程式 startup 完成 - context 已儲存")
}

// Shutdown 在應用程式關閉時調用
// 執行清理操作和資源釋放
func (a *App) Shutdown(ctx context.Context) {
	a.log.Info().Msg("應用程式正在關閉...")

	// 清理應用程式狀態
	a.state.Cleanup()

	a.log.Info().Msg("應用程式 shutdown 完成")
}

// DomReady 在前端 DOM 準備就緒時調用
// 可用於執行需要 DOM 的初始化操作
func (a *App) DomReady(ctx context.Context) {
	a.log.Debug().Msg("前端 DOM 已準備就緒")
}
