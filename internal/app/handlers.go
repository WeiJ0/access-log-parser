package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	
	"github.com/wailsapp/wails/v2/pkg/runtime"
	
	"access-log-analyzer/internal/exporter"
	"access-log-analyzer/internal/models"
	"access-log-analyzer/internal/parser"
	"access-log-analyzer/internal/stats"
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
func (a *App) SelectFile() SelectFileResponse {
	a.log.Info().
		Bool("a_ctx_nil", a.ctx == nil).
		Msg("SelectFile 被調用")
	
	// 檢查 a.ctx 是否為 nil
	if a.ctx == nil {
		a.log.Error().Msg("a.ctx 為 nil - Startup 可能未被調用或 context 未正確設置")
		return SelectFileResponse{
			Success:      false,
			ErrorMessage: "應用程式 context 未初始化。請重新啟動應用程式。",
		}
	}
	
	a.log.Info().
		Str("a_ctx_type", fmt.Sprintf("%T", a.ctx)).
		Msg("準備呼叫 runtime.OpenFileDialog")
	
	// 使用 a.ctx（Wails runtime context）而非參數 ctx
	// 這確保在 Wails 桌面環境中正確開啟原生對話框
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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
	
	a.log.Info().
		Str("filePath", filePath).
		Bool("has_error", err != nil).
		Msg("OpenFileDialog 返回")
	
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
func (a *App) ParseFile(req ParseFileRequest) ParseFileResponse {
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
	
	// 計算統計資訊（T070）
	a.log.Info().Msg("開始計算統計資訊")
	statStart := time.Now()
	
	calculator := stats.NewCalculator()
	statistics := calculator.Calculate(result.Entries)
	
	statTime := time.Since(statStart)
	
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
		Statistics:  statistics,  // 加入統計資料
		ParseTime:   result.ParseTime.Milliseconds(),
		StatTime:    statTime.Milliseconds(),  // 統計耗時（T071）
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
		Int64("stat_time_ms", statTime.Milliseconds()).  // 記錄統計耗時（T071）
		Msg("檔案解析完成")
	
	return ParseFileResponse{
		Success:      true,
		LogFile:      logFile,
		ErrorSamples: result.ErrorSamples,
	}
}

// ValidateLogFormat 快速驗證 log 檔案格式
// 讀取前 100 行並檢查是否符合 Apache log 格式
func (a *App) ValidateLogFormat(req ValidateFormatRequest) ValidateFormatResponse {
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
func (a *App) GetOpenFiles() []string {
	return a.state.GetTabs()
}

// CloseFile 關閉指定的檔案
// 從應用程式狀態中移除檔案資料
func (a *App) CloseFile(filePath string) bool {
	a.log.Info().Str("file", filePath).Msg("關閉檔案")
	
	a.state.RemoveFile(filePath)
	return true
}

// GetFileData 取得指定檔案的資料
// 用於切換分頁時重新載入資料
func (a *App) GetFileData(filePath string) *models.LogFile {
	logFile, exists := a.state.GetFile(filePath)
	if !exists {
		return nil
	}
	return logFile
}

// SetActiveFile 設定當前活動的檔案
// 用於追蹤使用者正在檢視的分頁
func (a *App) SetActiveFile(filePath string) bool {
	return a.state.SetActiveTab(filePath)
}

// GetActiveFile 取得當前活動的檔案路徑
// 返回使用者正在檢視的分頁
func (a *App) GetActiveFile() string {
	return a.state.GetActiveTab()
}

// ExportToExcelRequest 匯出 Excel 的請求參數
type ExportToExcelRequest struct {
	FilePath   string `json:"filePath"`   // 要匯出的 log 檔案路徑
	SavePath   string `json:"savePath"`   // Excel 檔案儲存路徑
}

// ExportToExcelResponse 匯出 Excel 的回應
type ExportToExcelResponse struct {
	Success       bool     `json:"success"`       // 是否成功
	ExportPath    string   `json:"exportPath"`    // 匯出檔案路徑
	FileSize      int64    `json:"fileSize"`      // 檔案大小
	TotalRecords  int64    `json:"totalRecords"`  // 總記錄數
	TruncatedRows int64    `json:"truncatedRows"` // 被截斷的行數
	Duration      string   `json:"duration"`      // 匯出耗時
	Warnings      []string `json:"warnings"`      // 警告訊息
	ErrorMessage  string   `json:"errorMessage"`  // 錯誤訊息
}

// SelectSaveLocationResponse 選擇儲存位置的回應
type SelectSaveLocationResponse struct {
	Success      bool   `json:"success"`      // 是否成功
	SavePath     string `json:"savePath"`     // 選擇的儲存路徑
	ErrorMessage string `json:"errorMessage"` // 錯誤訊息
}

// SelectSaveLocation 開啟儲存檔案對話框（T096）
// 讓使用者選擇 Excel 檔案的儲存位置
func (a *App) SelectSaveLocation(defaultName string) SelectSaveLocationResponse {
	a.log.Info().
		Str("defaultName", defaultName).
		Msg("SelectSaveLocation 被調用")
	
	// 檢查 context
	if a.ctx == nil {
		a.log.Error().Msg("context 未初始化")
		return SelectSaveLocationResponse{
			Success:      false,
			ErrorMessage: "應用程式 context 未初始化",
		}
	}
	
	// 如果沒有提供預設名稱，生成一個
	if defaultName == "" {
		defaultName = fmt.Sprintf("log_export_%s.xlsx", time.Now().Format("20060102_150405"))
	}
	
	// 確保副檔名是 .xlsx
	if filepath.Ext(defaultName) != ".xlsx" {
		defaultName = defaultName + ".xlsx"
	}
	
	// 開啟儲存對話框
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "選擇儲存位置",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel 檔案 (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
	})
	
	a.log.Info().
		Str("savePath", savePath).
		Bool("has_error", err != nil).
		Msg("SaveFileDialog 返回")
	
	if err != nil {
		a.log.Error().Err(err).Msg("儲存對話框錯誤")
		return SelectSaveLocationResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("無法開啟儲存對話框: %v", err),
		}
	}
	
	// 使用者取消選擇
	if savePath == "" {
		a.log.Debug().Msg("使用者取消儲存")
		return SelectSaveLocationResponse{
			Success:      false,
			ErrorMessage: "使用者取消儲存",
		}
	}
	
	a.log.Info().Str("path", savePath).Msg("使用者選擇儲存路徑")
	
	return SelectSaveLocationResponse{
		Success:  true,
		SavePath: savePath,
	}
}

// ExportToExcel 將日誌資料匯出為 Excel 檔案（T095）
// 包含日誌條目、統計資料和機器人偵測三個工作表
func (a *App) ExportToExcel(req ExportToExcelRequest) ExportToExcelResponse {
	a.log.Info().
		Str("sourceFile", req.FilePath).
		Str("savePath", req.SavePath).
		Msg("開始匯出 Excel")
	
	// 獲取檔案資料
	logFile, exists := a.state.GetFile(req.FilePath)
	if !exists {
		a.log.Error().Str("file", req.FilePath).Msg("檔案不存在於應用程式狀態中")
		return ExportToExcelResponse{
			Success:      false,
			ErrorMessage: "找不到檔案資料，請先載入檔案",
		}
	}
	
	// 檢查是否有資料
	if len(logFile.Entries) == 0 {
		a.log.Warn().Msg("檔案沒有資料可匯出")
		return ExportToExcelResponse{
			Success:      false,
			ErrorMessage: "檔案沒有資料可匯出",
		}
	}
	
	// 檢查統計資料
	if logFile.Statistics == nil {
		a.log.Error().Msg("統計資料不存在")
		return ExportToExcelResponse{
			Success:      false,
			ErrorMessage: "統計資料不存在，請先載入檔案",
		}
	}
	
	// 類型斷言統計資料
	// 目前系統使用 stats.Statistics，未來可能支援 models.Statistics
	statsData, ok := logFile.Statistics.(stats.Statistics)
	if !ok {
		// 嘗試指標類型
		statsDataPtr, ok := logFile.Statistics.(*stats.Statistics)
		if !ok {
			a.log.Error().
				Str("type", fmt.Sprintf("%T", logFile.Statistics)).
				Msg("統計資料類型錯誤")
			return ExportToExcelResponse{
				Success:      false,
				ErrorMessage: fmt.Sprintf("統計資料格式錯誤（類型: %T），請重新載入檔案", logFile.Statistics),
			}
		}
		statsData = *statsDataPtr
	}
	
	// 建立匯出器（T097：追蹤進度和日誌）
	xlsxExporter := exporter.NewXLSXExporter()
	
	a.log.Info().
		Int("entries", len(logFile.Entries)).
		Msg("準備匯出資料")
	
	// 轉換 Entries 為指標陣列
	entriesPtr := make([]*models.LogEntry, len(logFile.Entries))
	for i := range logFile.Entries {
		entriesPtr[i] = &logFile.Entries[i]
	}
	
	// 執行匯出
	startTime := time.Now()
	result, err := xlsxExporter.ExportWithStatsStatistics(entriesPtr, &statsData, req.SavePath)
	
	if err != nil {
		a.log.Error().
			Err(err).
			Str("savePath", req.SavePath).
			Msg("匯出失敗")
		
		return ExportToExcelResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("匯出失敗: %v", err),
		}
	}
	
	duration := time.Since(startTime)
	
	a.log.Info().
		Str("savePath", req.SavePath).
		Int64("fileSize", result.FileSize).
		Int64("records", result.TotalRecords).
		Int64("truncated", result.TruncatedRows).
		Str("duration", duration.String()).
		Int("warnings", len(result.Warnings)).
		Msg("Excel 匯出完成")
	
	return ExportToExcelResponse{
		Success:       true,
		ExportPath:    result.FilePath,
		FileSize:      result.FileSize,
		TotalRecords:  result.TotalRecords,
		TruncatedRows: result.TruncatedRows,
		Duration:      result.Duration,
		Warnings:      result.Warnings,
	}
}
