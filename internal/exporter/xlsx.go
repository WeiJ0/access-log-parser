package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"access-log-analyzer/internal/models"
	"access-log-analyzer/pkg/logger"

	"github.com/xuri/excelize/v2"
)

// XLSXExporter 提供 Excel XLSX 格式的匯出功能
// 使用 excelize 庫進行高效能的 Excel 檔案生成
type XLSXExporter struct {
	formatter     *Formatter // 資料格式化器
	streamingMode bool       // 是否使用串流模式
	logger        *logger.Logger
}

// ExportResult 匯出操作的結果資訊
type ExportResult struct {
	FilePath       string    `json:"filePath"`       // 匯出檔案路徑
	TotalRecords   int64     `json:"totalRecords"`   // 總記錄數
	FileSize       int64     `json:"fileSize"`       // 檔案大小（位元組）
	TruncatedRows  int64     `json:"truncatedRows"`  // 被截斷的行數（因 Excel 限制）
	Duration       string    `json:"duration"`       // 匯出耗時
	Warnings       []string  `json:"warnings"`       // 警告訊息
	CreatedAt      time.Time `json:"createdAt"`      // 建立時間
}

// Excel 相關常數
const (
	MaxExcelRows     = 1048576 // Excel 最大行數限制
	MaxExcelCols     = 16384   // Excel 最大列數限制
	DefaultSheetName = "Sheet1"
)

// NewXLSXExporter 建立新的 XLSX 匯出器
func NewXLSXExporter() *XLSXExporter {
	return &XLSXExporter{
		formatter:     NewFormatter(),
		streamingMode: true, // 預設使用串流模式以節省記憶體
		logger:        logger.Get(),
	}
}

// Export 執行完整的 Excel 檔案匯出
// 包含日誌條目、統計資料和機器人偵測三個工作表
func (e *XLSXExporter) Export(logs []*models.LogEntry, stats *models.Statistics, filePath string) (*ExportResult, error) {
	startTime := time.Now()
	
	e.logger.Info().
		Str("filePath", filePath).
		Int("logCount", len(logs)).
		Bool("streamingMode", e.streamingMode).
		Msg("開始 Excel 匯出")
	
	// 驗證輸入參數
	if err := e.validateInputs(logs, stats, filePath); err != nil {
		return nil, fmt.Errorf("輸入驗證失敗: %w", err)
	}
	
	// 資料驗證
	warnings := e.formatter.ValidateData(logs, stats)
	
	// 檢查 Excel 行數限制
	totalRows := int64(len(logs)) + 1 // +1 for header
	truncatedRows := int64(0)
	if totalRows > MaxExcelRows {
		truncatedRows = totalRows - MaxExcelRows
		logs = logs[:MaxExcelRows-1] // 保留標題行的空間
		warnings = append(warnings, fmt.Sprintf("資料超過 Excel 限制，截斷了 %d 行", truncatedRows))
	}
	
	// 創建 Excel 檔案
	f := excelize.NewFile()
	defer f.Close()
	
	// 刪除預設工作表
	f.DeleteSheet(DefaultSheetName)
	
	// 創建三個工作表
	if err := e.createLogEntriesWorksheet(f, logs); err != nil {
		return nil, fmt.Errorf("建立日誌條目工作表失敗: %w", err)
	}
	
	if err := e.createStatisticsWorksheet(f, stats); err != nil {
		return nil, fmt.Errorf("建立統計資料工作表失敗: %w", err)
	}
	
	if err := e.createBotDetectionWorksheet(f, logs); err != nil {
		return nil, fmt.Errorf("建立機器人偵測工作表失敗: %w", err)
	}
	
	// 設定預設工作表為日誌條目
	f.SetActiveSheet(0)
	
	// 確保目錄存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("建立目錄失敗: %w", err)
	}
	
	// 儲存檔案
	if err := f.SaveAs(filePath); err != nil {
		return nil, fmt.Errorf("儲存檔案失敗: %w", err)
	}
	
	// 取得檔案大小
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("取得檔案資訊失敗: %w", err)
	}
	
	duration := time.Since(startTime)
	
	result := &ExportResult{
		FilePath:      filePath,
		TotalRecords:  int64(len(logs)),
		FileSize:      fileInfo.Size(),
		TruncatedRows: truncatedRows,
		Duration:      duration.String(),
		Warnings:      warnings,
		CreatedAt:     time.Now(),
	}
	
	e.logger.Info().
		Str("filePath", filePath).
		Int64("records", result.TotalRecords).
		Int64("fileSize", result.FileSize).
		Str("duration", result.Duration).
		Int("warnings", len(warnings)).
		Msg("Excel 匯出完成")
	
	return result, nil
}

// createLogEntriesWorksheet 建立日誌條目工作表
func (e *XLSXExporter) createLogEntriesWorksheet(f *excelize.File, logs []*models.LogEntry) error {
	sheetName := "日誌條目"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("建立工作表失敗: %w", err)
	}
	
	// 格式化資料
	data := e.formatter.FormatLogEntries(logs)
	
	// 設定串流寫入器（如果啟用）
	if e.streamingMode && len(data) > 1000 {
		return e.writeDataWithStreaming(f, sheetName, data)
	}
	
	// 一般模式寫入
	return e.writeDataNormal(f, sheetName, data)
}

// createStatisticsWorksheet 建立統計資料工作表
func (e *XLSXExporter) createStatisticsWorksheet(f *excelize.File, stats *models.Statistics) error {
	sheetName := "統計資料"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("建立統計工作表失敗: %w", err)
	}
	
	// 格式化統計資料
	data := e.formatter.FormatStatistics(stats)
	
	// 寫入資料
	for rowIdx, row := range data {
		for colIdx, cell := range row {
			cellRef, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				continue
			}
			f.SetCellValue(sheetName, cellRef, cell)
		}
	}
	
	// 設定標題樣式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E6E6FA"}, Pattern: 1},
	})
	
	// 套用標題樣式到包含 "=====" 的行
	for rowIdx, row := range data {
		if len(row) > 0 && len(row[0]) > 5 && row[0][:5] == "=====" {
			cellRef, _ := excelize.CoordinatesToCellName(1, rowIdx+1)
			f.SetCellStyle(sheetName, cellRef, cellRef, headerStyle)
		}
	}
	
	return nil
}

// createBotDetectionWorksheet 建立機器人偵測工作表
func (e *XLSXExporter) createBotDetectionWorksheet(f *excelize.File, logs []*models.LogEntry) error {
	sheetName := "機器人偵測"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("建立機器人偵測工作表失敗: %w", err)
	}
	
	// 格式化機器人偵測資料
	data := e.formatter.FormatBotDetection(logs)
	
	// 寫入資料
	for rowIdx, row := range data {
		for colIdx, cell := range row {
			cellRef, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				continue
			}
			f.SetCellValue(sheetName, cellRef, cell)
		}
	}
	
	// 設定標題列樣式
	if len(data) > 0 {
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
			Border: []excelize.Border{
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})
		
		// 套用到第一行
		for colIdx := range data[0] {
			cellRef, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			f.SetCellStyle(sheetName, cellRef, cellRef, headerStyle)
		}
	}
	
	return nil
}

// writeDataWithStreaming 使用串流模式寫入大量資料
func (e *XLSXExporter) writeDataWithStreaming(f *excelize.File, sheetName string, data [][]string) error {
	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return fmt.Errorf("建立串流寫入器失敗: %w", err)
	}
	defer streamWriter.Flush()
	
	// 設定標題列樣式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
	})
	
	// 寫入資料
	for rowIdx, row := range data {
		cellData := make([]interface{}, len(row))
		for i, cell := range row {
			cellData[i] = cell
		}
		
		cellRef, _ := excelize.CoordinatesToCellName(1, rowIdx+1)
		
		// 標題行使用樣式
		if rowIdx == 0 {
			if err := streamWriter.SetRow(cellRef, cellData, excelize.RowOpts{StyleID: headerStyle}); err != nil {
				return fmt.Errorf("寫入標題行失敗: %w", err)
			}
		} else {
			if err := streamWriter.SetRow(cellRef, cellData); err != nil {
				return fmt.Errorf("寫入資料行 %d 失敗: %w", rowIdx+1, err)
			}
		}
	}
	
	return nil
}

// writeDataNormal 使用一般模式寫入資料
func (e *XLSXExporter) writeDataNormal(f *excelize.File, sheetName string, data [][]string) error {
	// 寫入所有資料
	for rowIdx, row := range data {
		for colIdx, cell := range row {
			cellRef, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				continue
			}
			f.SetCellValue(sheetName, cellRef, cell)
		}
	}
	
	// 設定標題列樣式
	if len(data) > 0 {
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
			Border: []excelize.Border{
				{Type: "bottom", Color: "FFFFFF", Style: 1},
			},
		})
		
		// 套用到第一行
		startCell, _ := excelize.CoordinatesToCellName(1, 1)
		endCell, _ := excelize.CoordinatesToCellName(len(data[0]), 1)
		f.SetCellStyle(sheetName, startCell, endCell, headerStyle)
		
		// 設定自動調整列寬
		for colIdx := range data[0] {
			colName := excelize.ToAlphaString(colIdx)
			f.SetColWidth(sheetName, colName, colName, 15)
		}
	}
	
	return nil
}

// validateInputs 驗證輸入參數
func (e *XLSXExporter) validateInputs(logs []*models.LogEntry, stats *models.Statistics, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("檔案路徑不能為空")
	}
	
	if filepath.Ext(filePath) != ".xlsx" {
		return fmt.Errorf("檔案路徑必須以 .xlsx 結尾")
	}
	
	// 檢查目錄是否可寫
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目錄不存在，嘗試創建
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("無法創建目錄 %s: %w", dir, err)
		}
	}
	
	return nil
}

// SetStreamingMode 設定串流模式
func (e *XLSXExporter) SetStreamingMode(enabled bool) {
	e.streamingMode = enabled
}

// GetStreamingMode 取得當前串流模式狀態
func (e *XLSXExporter) GetStreamingMode() bool {
	return e.streamingMode
}

// exportLogEntriesOnly 建立僅包含日誌條目的工作表（用於基準測試）
func (e *XLSXExporter) exportLogEntriesOnly(logs []*models.LogEntry, filePath string) error {
	f := excelize.NewFile()
	defer f.Close()
	
	// 刪除預設工作表
	f.DeleteSheet(DefaultSheetName)
	
	// 建立日誌條目工作表
	if err := e.createLogEntriesWorksheet(f, logs); err != nil {
		return err
	}
	
	return f.SaveAs(filePath)
}

// exportStatisticsOnly 建立僅包含統計資料的工作表（用於基準測試）
func (e *XLSXExporter) exportStatisticsOnly(stats *models.Statistics, filePath string) error {
	f := excelize.NewFile()
	defer f.Close()
	
	// 刪除預設工作表
	f.DeleteSheet(DefaultSheetName)
	
	// 建立統計資料工作表
	if err := e.createStatisticsWorksheet(f, stats); err != nil {
		return err
	}
	
	return f.SaveAs(filePath)
}

// exportBotDetectionOnly 建立僅包含機器人檢測的工作表（用於基準測試）
func (e *XLSXExporter) exportBotDetectionOnly(logs []*models.LogEntry, filePath string) error {
	f := excelize.NewFile()
	defer f.Close()
	
	// 刪除預設工作表
	f.DeleteSheet(DefaultSheetName)
	
	// 建立機器人偵測工作表
	if err := e.createBotDetectionWorksheet(f, logs); err != nil {
		return err
	}
	
	return f.SaveAs(filePath)
}

// GetEstimatedFileSize 估算檔案大小（位元組）
// 基於記錄數量和平均行大小的粗略估算
func (e *XLSXExporter) GetEstimatedFileSize(recordCount int) int64 {
	// 基於經驗值：每筆記錄約 200-300 位元組
	avgBytesPerRecord := int64(250)
	
	// Excel 檔案開銷（標題、格式、元資料等）
	baseOverhead := int64(50 * 1024) // 50KB 基礎開銷
	
	estimatedSize := baseOverhead + (int64(recordCount) * avgBytesPerRecord)
	
	// 考慮 Excel 壓縮（通常能壓縮 70-80%）
	compressionRatio := 0.3 // 壓縮後約為原始大小的 30%
	
	return int64(float64(estimatedSize) * compressionRatio)
}

// GetSupportedFormats 取得支援的匯出格式
func (e *XLSXExporter) GetSupportedFormats() []string {
	return []string{"xlsx"}
}

// GetMaxRecords 取得最大支援的記錄數
func (e *XLSXExporter) GetMaxRecords() int64 {
	return MaxExcelRows - 1 // 減去標題行
}