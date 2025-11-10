package exporter

import (
	"testing"
	"time"

	"access-log-analyzer/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatter 測試資料格式化器的基本功能
func TestFormatter(t *testing.T) {
	formatter := NewFormatter()
	require.NotNil(t, formatter, "格式化器不應為 nil")
}

// TestFormatLogEntries 測試日誌條目的格式化
func TestFormatLogEntries(t *testing.T) {
	logs := []*models.LogEntry{
		{
			IP:            "192.168.1.100",
			Timestamp:     time.Date(2024, 1, 1, 10, 30, 45, 0, time.UTC),
			Method:        "GET",
			URL:           "/index.html",
			Protocol:      "HTTP/1.1",
			StatusCode:    200,
			ResponseBytes: 1024,
			Referer:       "https://google.com",
			UserAgent:     "Mozilla/5.0 (Windows NT 10.0)",
			LineNumber:    1,
		},
		{
			IP:            "192.168.1.101",
			Timestamp:     time.Date(2024, 1, 1, 10, 31, 0, 0, time.UTC),
			Method:        "POST",
			URL:           "/api/login",
			Protocol:      "HTTP/1.1",
			StatusCode:    400,
			ResponseBytes: 512,
			Referer:       "-",
			UserAgent:     "curl/7.68.0",
			LineNumber:    2,
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatLogEntries(logs)

	// 檢查返回的資料結構
	require.Len(t, result, 3, "應該有 3 行（1 標題 + 2 資料）")

	// 檢查標題行
	expectedHeaders := []string{
		"IP位址", "時間戳", "HTTP方法", "URL", "協定",
		"狀態碼", "回應大小", "來源頁面", "User Agent",
	}
	assert.Equal(t, expectedHeaders, result[0], "標題行應該正確")

	// 檢查第一筆資料
	expectedFirstRow := []string{
		"192.168.1.100",
		"2024-01-01 10:30:45",
		"GET",
		"/index.html",
		"HTTP/1.1",
		"200",
		"1024",
		"https://google.com",
		"Mozilla/5.0 (Windows NT 10.0)",
	}
	assert.Equal(t, expectedFirstRow, result[1], "第一行資料應該正確")

	// 檢查第二筆資料
	expectedSecondRow := []string{
		"192.168.1.101",
		"2024-01-01 10:31:00",
		"POST",
		"/api/login",
		"HTTP/1.1",
		"400",
		"512",
		"-",
		"curl/7.68.0",
	}
	assert.Equal(t, expectedSecondRow, result[2], "第二行資料應該正確")
}

// TestFormatLogEntriesEmpty 測試空日誌條目的格式化
func TestFormatLogEntriesEmpty(t *testing.T) {
	var logs []*models.LogEntry

	formatter := NewFormatter()
	result := formatter.FormatLogEntries(logs)

	// 即使沒有資料，也應該有標題行
	require.Len(t, result, 1, "應該有 1 行（僅標題）")

	expectedHeaders := []string{
		"IP位址", "時間戳", "HTTP方法", "URL", "協定",
		"狀態碼", "回應大小", "來源頁面", "User Agent",
	}
	assert.Equal(t, expectedHeaders, result[0], "標題行應該正確")
}

// TestFormatStatistics 測試統計資料的格式化
func TestFormatStatistics(t *testing.T) {
	stats := &models.Statistics{
		TotalRequests:    1000,
		UniqueIPs:        50,
		TotalBytes:       1048576,
		StartTime:        "2024-01-01 10:00:00",
		EndTime:          "2024-01-01 11:00:00",
		ErrorCount:       100,
		ClientErrorCount: 80,
		ServerErrorCount: 20,
		ErrorRate:        10.0,
		StatusCodeDist: map[int]int64{
			200: 800,
			404: 150,
			500: 50,
		},
		StatusCatDist: map[string]int64{
			"2xx": 800,
			"4xx": 150,
			"5xx": 50,
		},
		TopIPs: []models.IPStat{
			{IP: "192.168.1.100", Count: 100, TotalBytes: 102400},
			{IP: "192.168.1.101", Count: 80, TotalBytes: 81920},
		},
		TopURLs: []models.URLStat{
			{URL: "/index.html", Count: 200, TotalBytes: 204800},
			{URL: "/api/users", Count: 150, TotalBytes: 153600},
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatStatistics(stats)

	// 檢查返回的資料結構不為空
	require.Greater(t, len(result), 0, "格式化的統計資料應該不為空")

	// 檢查基本統計資料是否包含
	found := false
	for _, row := range result {
		if len(row) >= 2 && row[0] == "總請求數" && row[1] == "1000" {
			found = true
			break
		}
	}
	assert.True(t, found, "應該包含總請求數統計")

	// 檢查是否包含 Top IPs 資料
	foundTopIP := false
	for _, row := range result {
		if len(row) >= 2 && row[0] == "192.168.1.100" {
			foundTopIP = true
			break
		}
	}
	assert.True(t, foundTopIP, "應該包含 Top IP 資料")
}

// TestFormatBotDetection 測試機器人偵測資料的格式化
func TestFormatBotDetection(t *testing.T) {
	logs := []*models.LogEntry{
		{
			IP:         "192.168.1.100",
			UserAgent:  "Googlebot/2.1 (+http://www.google.com/bot.html)",
			LineNumber: 1,
		},
		{
			IP:         "192.168.1.101",
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			LineNumber: 2,
		},
		{
			IP:         "192.168.1.102",
			UserAgent:  "facebookexternalhit/1.1",
			LineNumber: 3,
		},
		{
			IP:         "192.168.1.100", // 同一個 Googlebot
			UserAgent:  "Googlebot/2.1 (+http://www.google.com/bot.html)",
			LineNumber: 4,
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatBotDetection(logs)

	// 檢查返回的資料結構
	require.Greater(t, len(result), 1, "應該有標題行和資料行")

	// 檢查標題行
	expectedHeaders := []string{"IP位址", "機器人類型", "信心分數", "請求次數"}
	assert.Equal(t, expectedHeaders, result[0], "標題行應該正確")

	// 檢查是否偵測到 Googlebot
	foundGooglebot := false
	foundFacebook := false

	for i := 1; i < len(result); i++ {
		row := result[i]
		require.Len(t, row, 4, "每行資料應該有 4 個欄位")

		if row[0] == "192.168.1.100" && row[1] == "搜尋引擎" {
			foundGooglebot = true
			assert.Equal(t, "高", row[2], "Googlebot 信心分數應該是高")
			assert.Equal(t, "2", row[3], "Googlebot 請求次數應該是 2")
		}

		if row[0] == "192.168.1.102" && row[1] == "社交媒體" {
			foundFacebook = true
			// 社交媒體機器人，請求次數只有 1，信心分數應該是「中」（特異性 4 + 頻率 0 = 4）
			assert.Equal(t, "中", row[2], "Facebook 信心分數應該是中（單次請求）")
		}
	}

	assert.True(t, foundGooglebot, "應該偵測到 Googlebot")
	assert.True(t, foundFacebook, "應該偵測到 Facebook 機器人")
}

// TestFormatBotDetectionNoBots 測試沒有機器人的情況
func TestFormatBotDetectionNoBots(t *testing.T) {
	logs := []*models.LogEntry{
		{
			IP:         "192.168.1.100",
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			LineNumber: 1,
		},
		{
			IP:         "192.168.1.101",
			UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			LineNumber: 2,
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatBotDetection(logs)

	// 即使沒有機器人，也應該有標題行
	require.Len(t, result, 1, "應該只有標題行")

	expectedHeaders := []string{"IP位址", "機器人類型", "信心分數", "請求次數"}
	assert.Equal(t, expectedHeaders, result[0], "標題行應該正確")
}

// TestTimeFormatting 測試時間格式化
func TestTimeFormatting(t *testing.T) {
	// 測試不同時區的時間
	utcTime := time.Date(2024, 1, 1, 15, 30, 45, 0, time.UTC)

	logs := []*models.LogEntry{
		{
			IP:         "192.168.1.100",
			Timestamp:  utcTime,
			Method:     "GET",
			URL:        "/test",
			StatusCode: 200,
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatLogEntries(logs)

	// 檢查時間格式
	require.Len(t, result, 2, "應該有標題和一行資料")
	timeStr := result[1][1] // 時間戳欄位
	assert.Equal(t, "2024-01-01 15:30:45", timeStr, "時間格式應該正確")
}

// TestSpecialCharacterHandling 測試特殊字元處理
func TestSpecialCharacterHandling(t *testing.T) {
	logs := []*models.LogEntry{
		{
			IP:         "192.168.1.100",
			Method:     "GET",
			URL:        "/path with spaces",
			UserAgent:  "User Agent with \"quotes\" and ,commas,",
			Referer:    "https://example.com/path?param=value&other=測試",
			StatusCode: 200,
			Timestamp:  time.Now(),
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatLogEntries(logs)

	require.Len(t, result, 2, "應該有標題和一行資料")

	// 檢查特殊字元是否正確處理
	dataRow := result[1]
	assert.Equal(t, "/path with spaces", dataRow[3], "URL 中的空格應該保留")
	assert.Contains(t, dataRow[8], "quotes", "User Agent 中的引號應該正確處理")
	assert.Contains(t, dataRow[7], "測試", "Referer 中的中文字元應該正確處理")
}

// TestLargeDataFormatting 測試大量資料格式化效能
func TestLargeDataFormatting(t *testing.T) {
	// 創建大量測試資料
	count := 10000
	logs := make([]*models.LogEntry, count)

	for i := 0; i < count; i++ {
		logs[i] = &models.LogEntry{
			IP:         "192.168.1.100",
			Timestamp:  time.Now().Add(time.Duration(i) * time.Second),
			Method:     "GET",
			URL:        "/test",
			StatusCode: 200,
		}
	}

	formatter := NewFormatter()
	result := formatter.FormatLogEntries(logs)

	// 檢查結果數量
	assert.Len(t, result, count+1, "應該有正確數量的資料行（包含標題）")

	// 檢查標題行
	expectedHeaders := []string{
		"IP位址", "時間戳", "HTTP方法", "URL", "協定",
		"狀態碼", "回應大小", "來源頁面", "User Agent",
	}
	assert.Equal(t, expectedHeaders, result[0], "標題行應該正確")
}

// TestNilDataHandling 測試 nil 資料處理
func TestNilDataHandling(t *testing.T) {
	formatter := NewFormatter()

	// 測試 nil 日誌條目列表
	result := formatter.FormatLogEntries(nil)
	require.Len(t, result, 1, "nil 列表應該只返回標題行")

	// 測試包含 nil 條目的列表
	logs := []*models.LogEntry{
		{
			IP:         "192.168.1.100",
			Method:     "GET",
			StatusCode: 200,
			Timestamp:  time.Now(),
		},
		nil, // nil 條目
		{
			IP:         "192.168.1.101",
			Method:     "POST",
			StatusCode: 404,
			Timestamp:  time.Now(),
		},
	}

	result = formatter.FormatLogEntries(logs)
	// 應該跳過 nil 條目，只處理有效的條目
	assert.Len(t, result, 3, "應該跳過 nil 條目，只處理有效條目")
}

// TestZeroValueHandling 測試零值處理
func TestZeroValueHandling(t *testing.T) {
	logs := []*models.LogEntry{
		{
			IP:            "192.168.1.100",
			Timestamp:     time.Time{}, // 零值時間
			Method:        "",          // 空字串
			URL:           "",
			StatusCode:    0, // 零值狀態碼
			ResponseBytes: 0,
		},
	}

	formatter := NewFormatter()
	result := formatter.FormatLogEntries(logs)

	require.Len(t, result, 2, "應該有標題和一行資料")

	dataRow := result[1]
	assert.Equal(t, "192.168.1.100", dataRow[0], "IP 應該正確")
	assert.Equal(t, "0001-01-01 00:00:00", dataRow[1], "零值時間應該有默認格式")
	assert.Equal(t, "", dataRow[2], "空 Method 應該保持空字串")
	assert.Equal(t, "0", dataRow[5], "零值狀態碼應該格式化為 '0'")
}
