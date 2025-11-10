package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger 包裝 zerolog.Logger
// 提供應用程式級別的日誌功能
type Logger struct {
	*zerolog.Logger
}

var globalLogger *Logger

// Init 初始化全局日誌記錄器
// 設定日誌格式、輸出目標和日誌級別
func Init() {
	// 設定人類可讀的格式（開發模式）
	// 生產環境可切換為 JSON 格式
	zerolog.TimeFieldFormat = time.RFC3339

	// 創建多輸出：控制台 + 檔案
	writers := []io.Writer{
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		},
	}

	// 嘗試創建日誌檔案
	logFile, err := os.OpenFile(
		filepath.Join(".", "app.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err == nil {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        logFile,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    true,
		})
	}

	output := io.MultiWriter(writers...)

	// 建立日誌記錄器
	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	// 設定全局日誌級別（可從環境變數讀取）
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // 預設為 info
	}

	globalLogger = &Logger{&logger}
	log.Logger = logger // 更新 zerolog 的全局 logger

	globalLogger.Info().Msg("日誌系統初始化完成")
}

// Get 取得全局日誌記錄器實例
// 如果尚未初始化則自動初始化
func Get() *Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

// WithContext 建立帶有上下文欄位的新日誌記錄器
// 用於在特定範圍內添加額外的日誌欄位
func (l *Logger) WithContext(key, value string) *Logger {
	newLogger := l.With().Str(key, value).Logger()
	return &Logger{&newLogger}
}

// WithModule 建立帶有模組名稱的日誌記錄器
// 便捷方法，用於模組化日誌記錄
func (l *Logger) WithModule(module string) *Logger {
	return l.WithContext("module", module)
}

// WithFile 建立帶有檔案路徑的日誌記錄器
// 便捷方法，用於檔案操作相關的日誌
func (l *Logger) WithFile(filepath string) *Logger {
	return l.WithContext("file", filepath)
}

// SetLevel 動態設定日誌級別
// 可在運行時調整日誌輸出詳細程度
func SetLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	Get().Info().Str("level", level).Msg("日誌級別已更新")
}
