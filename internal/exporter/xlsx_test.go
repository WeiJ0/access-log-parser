package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"access-log-analyzer/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

// TestXLSXExporter 測試 XLSX 匯出器的基本功能
func TestXLSXExporter(t *testing.T) {
	// 創建測試資料
	logs := createTestLogEntries()
	stats := createTestStatistics()
	
	// 創建臨時檔案
	tempFile := filepath.Join(t.TempDir(), "test_export.xlsx")
	
	// 創建匯出器
	exporter := NewXLSXExporter()
	require.NotNil(t, exporter, "匯出器不應為 nil")
	
	// 執行匯出
	result, err := exporter.Export(logs, stats, tempFile)
	require.NoError(t, err, "匯出應該成功")
	require.NotNil(t, result, "匯出結果不應為 nil")
	
	// 驗證檔案存在
	assert.FileExists(t, tempFile, "匯出檔案應該存在")
	
	// 驗證檔案大小 > 0
	fileInfo, err := os.Stat(tempFile)
	require.NoError(t, err)
	assert.Greater(t, fileInfo.Size(), int64(0), "檔案大小應該大於 0")
	
	// 驗證 Excel 檔案結構
	f, err := excelize.OpenFile(tempFile)
	require.NoError(t, err, "應該能開啟 Excel 檔案")
	defer f.Close()
	
	// 驗證工作表存在
	sheets := f.GetSheetList()
	expectedSheets := []string{"日誌條目", "統計資料", "機器人偵測"}
	
	assert.Len(t, sheets, 3, "應該有 3 個工作表")
	for _, expectedSheet := range expectedSheets {
		assert.Contains(t, sheets, expectedSheet, "應該包含工作表：%s", expectedSheet)
	}
}

// TestLogEntriesWorksheet 測試日誌條目工作表的內容
func TestLogEntriesWorksheet(t *testing.T) {
	logs := createTestLogEntries()
	stats := createTestStatistics()
	tempFile := filepath.Join(t.TempDir(), "test_logs.xlsx")
	
	exporter := NewXLSXExporter()
	_, err := exporter.Export(logs, stats, tempFile)
	require.NoError(t, err)
	
	// 開啟檔案驗證內容
	f, err := excelize.OpenFile(tempFile)
	require.NoError(t, err)
	defer f.Close()
	
	// 驗證日誌條目工作表標題
	expectedHeaders := []string{
		"IP位址", "時間戳", "HTTP方法", "URL", "協定",
		"狀態碼", "回應大小", "來源頁面", "User Agent",
	}
	
	for i, header := range expectedHeaders {
		colName, err := excelize.ColumnNumberToName(i + 1)
		require.NoError(t, err)
		cell, err := f.GetCellValue("日誌條目", colName+"1")
		require.NoError(t, err)
		assert.Equal(t, header, cell, "標題欄位 %d 應該是 %s", i+1, header)
	}
	
	// 驗證資料行數（標題 + 資料行）
	rows, err := f.GetRows("日誌條目")
	require.NoError(t, err)
	assert.Equal(t, len(logs)+1, len(rows), "工作表應該有 %d 行（包含標題）", len(logs)+1)
	
	// 驗證第一筆資料
	if len(logs) > 0 {
		firstLog := logs[0]
		ip, _ := f.GetCellValue("日誌條目", "A2")
		assert.Equal(t, firstLog.IP, ip, "第一筆記錄的 IP 應該正確")
		
		method, _ := f.GetCellValue("日誌條目", "C2")
		assert.Equal(t, firstLog.Method, method, "第一筆記錄的 HTTP 方法應該正確")
	}
}

// TestStatisticsWorksheet 測試統計資料工作表的內容
func TestStatisticsWorksheet(t *testing.T) {
	logs := createTestLogEntries()
	stats := createTestStatistics()
	tempFile := filepath.Join(t.TempDir(), "test_stats.xlsx")
	
	exporter := NewXLSXExporter()
	_, err := exporter.Export(logs, stats, tempFile)
	require.NoError(t, err)
	
	// 開啟檔案驗證內容
	f, err := excelize.OpenFile(tempFile)
	require.NoError(t, err)
	defer f.Close()
	
	// 驗證統計資料工作表存在
	rows, err := f.GetRows("統計資料")
	require.NoError(t, err)
	assert.Greater(t, len(rows), 0, "統計資料工作表應該有內容")
	
	// 驗證基本統計資料
	row1Label, _ := f.GetCellValue("統計資料", "A1")
	// 跳過分隔標題行，找到實際的資料行
	if row1Label == "===== 基本統計 =====" {
		// 標題行是第一行，實際標題是第二行
		actualLabelRow, _ := f.GetCellValue("統計資料", "A2")
		assert.Equal(t, "統計項目", actualLabelRow)
		
		// 總請求數在第三行
		totalRequestsLabel, _ := f.GetCellValue("統計資料", "A3")
		assert.Equal(t, "總請求數", totalRequestsLabel)
		
		totalRequestsValue, _ := f.GetCellValue("統計資料", "B3")
		assert.Equal(t, "3", totalRequestsValue) // 基於測試資料
	} else {
		// 舊格式兼容
		assert.Equal(t, "總請求數", row1Label)
		totalRequestsValue, _ := f.GetCellValue("統計資料", "B1")
		assert.Equal(t, "3", totalRequestsValue)
	}
}

// TestBotDetectionWorksheet 測試機器人偵測工作表的內容
func TestBotDetectionWorksheet(t *testing.T) {
	logs := createTestLogEntries()
	stats := createTestStatistics()
	tempFile := filepath.Join(t.TempDir(), "test_bots.xlsx")
	
	exporter := NewXLSXExporter()
	_, err := exporter.Export(logs, stats, tempFile)
	require.NoError(t, err)
	
	// 開啟檔案驗證內容
	f, err := excelize.OpenFile(tempFile)
	require.NoError(t, err)
	defer f.Close()
	
	// 驗證機器人偵測工作表標題
	expectedHeaders := []string{"IP位址", "機器人類型", "信心分數", "請求次數"}
	
	for i, header := range expectedHeaders {
		colName, err := excelize.ColumnNumberToName(i + 1)
		require.NoError(t, err)
		cell, err := f.GetCellValue("機器人偵測", colName+"1")
		require.NoError(t, err)
		assert.Equal(t, header, cell, "標題欄位 %d 應該是 %s", i+1, header)
	}
}

// TestExportEmptyData 測試空資料的匯出
func TestExportEmptyData(t *testing.T) {
	var logs []*models.LogEntry
	stats := models.NewStatistics()
	tempFile := filepath.Join(t.TempDir(), "test_empty.xlsx")
	
	exporter := NewXLSXExporter()
	result, err := exporter.Export(logs, stats, tempFile)
	
	// 空資料應該不會導致錯誤
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalRecords)
	assert.FileExists(t, tempFile)
}

// TestExportLargeData 測試大量資料的匯出
func TestExportLargeData(t *testing.T) {
	// 創建大量測試資料（模擬 10000 筆記錄）
	logs := make([]*models.LogEntry, 10000)
	for i := 0; i < 10000; i++ {
		logs[i] = &models.LogEntry{
			IP:            "192.168.1.1",
			Timestamp:     time.Now().Add(time.Duration(i) * time.Second),
			Method:        "GET",
			URL:           "/test",
			Protocol:      "HTTP/1.1",
			StatusCode:    200,
			ResponseBytes: 1024,
			UserAgent:     "Test Agent",
		}
	}
	
	stats := createTestStatistics()
	tempFile := filepath.Join(t.TempDir(), "test_large.xlsx")
	
	exporter := NewXLSXExporter()
	result, err := exporter.Export(logs, stats, tempFile)
	
	require.NoError(t, err)
	assert.Equal(t, int64(10000), result.TotalRecords)
	assert.FileExists(t, tempFile)
	
	// 驗證檔案大小合理（應該 > 100KB，因為 10K 記錄壓縮後約 300-400KB）
	fileInfo, err := os.Stat(tempFile)
	require.NoError(t, err)
	assert.Greater(t, fileInfo.Size(), int64(100*1024), "大檔案應該 > 100KB")
}

// TestExportWithRowLimit 測試 Excel 行數限制處理
func TestExportWithRowLimit(t *testing.T) {
	exporter := NewXLSXExporter()
	
	// 測試接近 Excel 限制的資料量
	maxExcelRows := int64(1048576) // Excel 最大行數
	testRows := maxExcelRows + 1000 // 超過限制
	
	// 模擬大量資料（不實際創建，只測試限制邏輯）
	logs := make([]*models.LogEntry, testRows)
	for i := int64(0); i < testRows; i++ {
		logs[i] = &models.LogEntry{
			IP:         "192.168.1.1",
			Timestamp:  time.Now(),
			Method:     "GET",
			URL:        "/test",
			StatusCode: 200,
		}
	}
	
	stats := createTestStatistics()
	tempFile := filepath.Join(t.TempDir(), "test_limit.xlsx")
	
	// 這應該處理行數限制並返回警告
	result, err := exporter.Export(logs, stats, tempFile)
	
	// 不應該因為超過限制而失敗，而是截斷並給出警告
	require.NoError(t, err)
	assert.True(t, result.TruncatedRows > 0, "應該有截斷的行數")
	// 檢查警告訊息中是否包含關於截斷的訊息
	hasWarning := false
	for _, warning := range result.Warnings {
		t.Logf("警告訊息: %s", warning)
		if strings.Contains(warning, "截斷") || strings.Contains(warning, "限制") {
			hasWarning = true
			break
		}
	}
	assert.True(t, hasWarning, "應該有關於資料截斷的警告，實際警告: %v", result.Warnings)
}

// TestInvalidFilePath 測試無效檔案路徑
func TestInvalidFilePath(t *testing.T) {
	logs := createTestLogEntries()
	stats := createTestStatistics()
	
	// 使用無效路徑（包含非法字元）
	invalidPath := "Z:\\:invalid:path\\test.xlsx" // Windows 不允許冒號在路徑中
	
	exporter := NewXLSXExporter()
	_, err := exporter.Export(logs, stats, invalidPath)
	
	// 應該返回錯誤
	assert.Error(t, err, "無效路徑應該返回錯誤")
}

// Helper functions for test data

// createTestLogEntries 創建測試用的日誌條目
func createTestLogEntries() []*models.LogEntry {
	baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	
	return []*models.LogEntry{
		{
			IP:            "192.168.1.100",
			Timestamp:     baseTime,
			Method:        "GET",
			URL:           "/index.html",
			Protocol:      "HTTP/1.1",
			StatusCode:    200,
			ResponseBytes: 1024,
			Referer:       "https://google.com",
			UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			LineNumber:    1,
		},
		{
			IP:            "192.168.1.101",
			Timestamp:     baseTime.Add(1 * time.Minute),
			Method:        "POST",
			URL:           "/api/login",
			Protocol:      "HTTP/1.1",
			StatusCode:    400,
			ResponseBytes: 512,
			Referer:       "-",
			UserAgent:     "curl/7.68.0",
			LineNumber:    2,
		},
		{
			IP:            "192.168.1.102",
			Timestamp:     baseTime.Add(2 * time.Minute),
			Method:        "GET",
			URL:           "/robots.txt",
			Protocol:      "HTTP/1.1",
			StatusCode:    404,
			ResponseBytes: 256,
			Referer:       "-",
			UserAgent:     "Googlebot/2.1 (+http://www.google.com/bot.html)",
			LineNumber:    3,
		},
	}
}

// createTestStatistics 創建測試用的統計資料
func createTestStatistics() *models.Statistics {
	stats := models.NewStatistics()
	
	stats.TotalRequests = 3
	stats.UniqueIPs = 3
	stats.TotalBytes = 1792
	stats.StartTime = "2024-01-01 10:00:00"
	stats.EndTime = "2024-01-01 10:02:00"
	
	// 狀態碼分布
	stats.StatusCodeDist[200] = 1
	stats.StatusCodeDist[400] = 1
	stats.StatusCodeDist[404] = 1
	
	// 狀態類別分布
	stats.StatusCatDist["2xx"] = 1
	stats.StatusCatDist["4xx"] = 2
	
	// Top IPs
	stats.TopIPs = []models.IPStat{
		{IP: "192.168.1.100", Count: 1, TotalBytes: 1024},
		{IP: "192.168.1.101", Count: 1, TotalBytes: 512},
		{IP: "192.168.1.102", Count: 1, TotalBytes: 256},
	}
	
	// Top URLs
	stats.TopURLs = []models.URLStat{
		{URL: "/index.html", Count: 1, TotalBytes: 1024},
		{URL: "/api/login", Count: 1, TotalBytes: 512},
		{URL: "/robots.txt", Count: 1, TotalBytes: 256},
	}
	
	// 錯誤統計
	stats.ErrorCount = 2
	stats.ClientErrorCount = 2
	stats.ServerErrorCount = 0
	stats.ErrorRate = 66.67
	
	return stats
}