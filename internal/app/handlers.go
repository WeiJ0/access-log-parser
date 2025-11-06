package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/wailsapp/wails/v2/pkg/runtime"
	
	"access-log-analyzer/internal/models"
	"access-log-analyzer/internal/parser"
)

// ParseFileRequest 解析檔案的請求參數
type ParseFileRequest struct {
	FilePath string `json:"filePath"` // 檔案路徑
}

// ParseFileResponse 解析檔案的回應
type ParseFileResponse struct {
	Success      bool                 `json:"success"`      // 是否成功
	LogFile      *models.LogFile      `json:"logFile"`      // 日誌檔案資料
	ErrorMessage string               `json:"errorMessage"` // 錯誤訊息
	ErrorSamples []parser.ParseError  `json:"errorSamples"` // 錯誤樣本
}

// SelectFileResponse 選擇檔案的回應
type SelectFileResponse struct {
	Success      bool   `json:"success"`      // 是否成功
	FilePath     string `json:"filePath"`     // 選擇的檔案路徑
	ErrorMessage string `json:"errorMessage"` // 錯誤訊息
}

// ValidateFormatRequest 驗證格式的請求參數
type ValidateFormatRequest struct {
	FilePath string `json:"filePath"` // 檔案路徑
}

// ValidateFormatResponse 驗證格式的回應
type ValidateFormatResponse struct {
	Success      bool   `json:"success"`      // 是否成功
	Valid        bool   `json:"valid"`        // 格式是否有效
	ErrorMessage string `json:"errorMessage"` // 錯誤訊息
}

// SelectFile 開啟檔案選擇對話框
// 讓使用者選擇要解析的 log 檔案
func (a *App) SelectFile(ctx context.Context) SelectFileResponse {
	a.log.Info().Msg("開啟檔案選擇對話框")
	
	// 開啟檔案選擇對話框
	filePath, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "選擇 Apache Log 檔案",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Log 檔案 (*.log)",
				Pattern:     "*.log",
			},
			{
				DisplayName: "所有檔案 (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	
	if err != nil {
		a.log.Error().Err(err).Msg("檔案選擇對話框錯誤")
		return SelectFileResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("無法開啟檔案選擇對話框: %v", err),
		}
	}
	
	// 使用者取消選擇
	if filePath == "" {
		a.log.Debug().Msg("使用者取消檔案選擇")
		return SelectFileResponse{
			Success:      false,
			ErrorMessage: "使用者取消選擇",
		}
	}
	
	a.log.Info().Str("path", filePath).Msg("使用者選擇檔案")
	
	return SelectFileResponse{
		Success:  true,
		FilePath: filePath,
	}
}

// ParseFile 解析指定的 log 檔案
// 返回解析結果和統計資訊
func (a *App) ParseFile(ctx context.Context, req ParseFileRequest) ParseFileResponse {
	a.log.Info().Str("file", req.FilePath).Msg("開始解析檔案")
	
	// 檢查檔案是否存在
	fileInfo, err := os.Stat(req.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ParseFileResponse{
				Success:      false,
				ErrorMessage: "檔案不存在",
			}
		}
		return ParseFileResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("無法讀取檔案資訊: %v", err),
		}
	}
	
	// 檢查檔案大小限制（10GB）
	const maxFileSize = 10 * 1024 * 1024 * 1024 // 10GB
	if fileInfo.Size() > maxFileSize {
		return ParseFileResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("檔案過大: %.2f GB（限制 10 GB）", float64(fileInfo.Size())/(1024*1024*1024)),
		}
	}
	
	// 建立解析器（自動使用所有 CPU 核心）
	logParser := parser.NewParser(parser.FormatCombined, 0)
	
	// 解析檔案
	result, err := logParser.ParseFile(req.FilePath, fileInfo.Size())
	if err != nil {
		return ParseFileResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("解析失敗: %v", err),
		}
	}
	
	// 建立 LogFile 物件
	logFile := &models.LogFile{
		Path:        req.FilePath,
		Name:        filepath.Base(req.FilePath),
		Size:        fileInfo.Size(),
		LoadedAt:    fileInfo.ModTime(),
		TotalLines:  result.TotalLines,
		ParsedLines: result.ParsedLines,
		ErrorLines:  result.ErrorLines,
		Entries:     result.Entries,
		ParseTime:   result.ParseTime.Milliseconds(),
		MemoryUsed:  result.MemoryUsed,
	}
	
	// 將檔案新增到應用程式狀態
	a.state.AddFile(req.FilePath, logFile)
	
	a.log.Info().
		Str("file", req.FilePath).
		Int("total", result.TotalLines).
		Int("parsed", result.ParsedLines).
		Int("errors", result.ErrorLines).
		Float64("throughput_mb_s", result.ThroughputMB).
		Msg("檔案解析完成")
	
	return ParseFileResponse{
		Success:      true,
		LogFile:      logFile,
		ErrorSamples: result.ErrorSamples,
	}
}

// ValidateLogFormat 快速驗證 log 檔案格式
// 讀取前 100 行並檢查是否符合 Apache log 格式
func (a *App) ValidateLogFormat(ctx context.Context, req ValidateFormatRequest) ValidateFormatResponse {
	a.log.Info().Str("file", req.FilePath).Msg("驗證檔案格式")
	
	// 檢查檔案是否存在
	if _, err := os.Stat(req.FilePath); err != nil {
		if os.IsNotExist(err) {
			return ValidateFormatResponse{
				Success:      false,
				ErrorMessage: "檔案不存在",
			}
		}
		return ValidateFormatResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("無法讀取檔案: %v", err),
		}
	}
	
	// 建立解析器並驗證格式
	logParser := parser.NewParser(parser.FormatCombined, 1)
	valid, err := logParser.ValidateFormat(req.FilePath, 100)
	
	if err != nil {
		return ValidateFormatResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("驗證失敗: %v", err),
		}
	}
	
	a.log.Info().
		Str("file", req.FilePath).
		Bool("valid", valid).
		Msg("格式驗證完成")
	
	return ValidateFormatResponse{
		Success: true,
		Valid:   valid,
	}
}

// GetOpenFiles 取得所有已開啟的檔案列表
// 返回檔案路徑和基本資訊
func (a *App) GetOpenFiles(ctx context.Context) []string {
	return a.state.GetTabs()
}

// CloseFile 關閉指定的檔案
// 從應用程式狀態中移除檔案資料
func (a *App) CloseFile(ctx context.Context, filePath string) bool {
	a.log.Info().Str("file", filePath).Msg("關閉檔案")
	
	a.state.RemoveFile(filePath)
	return true
}

// GetFileData 取得指定檔案的資料
// 用於切換分頁時重新載入資料
func (a *App) GetFileData(ctx context.Context, filePath string) *models.LogFile {
	logFile, exists := a.state.GetFile(filePath)
	if !exists {
		return nil
	}
	return logFile
}

// SetActiveFile 設定當前活動的檔案
// 用於追蹤使用者正在檢視的分頁
func (a *App) SetActiveFile(ctx context.Context, filePath string) bool {
	return a.state.SetActiveTab(filePath)
}

// GetActiveFile 取得當前活動的檔案路徑
// 返回使用者正在檢視的分頁
func (a *App) GetActiveFile(ctx context.Context) string {
	return a.state.GetActiveTab()
}
