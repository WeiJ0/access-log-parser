package exporter

import (
	"path/filepath"
	"testing"
	"time"

	"access-log-analyzer/internal/models"
)

// BenchmarkXLSXExport 基準測試 XLSX 匯出效能
// 目標：1M 記錄匯出 ≤30 秒
func BenchmarkXLSXExport(b *testing.B) {
	// 測試不同資料量的匯出效能
	benchmarks := []struct {
		name     string
		records  int
		expected time.Duration // 預期最大執行時間
	}{
		{"1K_records", 1000, 1 * time.Second},
		{"10K_records", 10000, 5 * time.Second},
		{"100K_records", 100000, 15 * time.Second},
		{"1M_records", 1000000, 30 * time.Second},
	}
	
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// 創建測試資料
			logs := createBenchmarkLogEntries(bm.records)
			stats := createBenchmarkStatistics(bm.records)
			
			// 重置計時器
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				tempFile := filepath.Join(b.TempDir(), "benchmark.xlsx")
				
				exporter := NewXLSXExporter()
				start := time.Now()
				
				_, err := exporter.Export(logs, stats, tempFile)
				if err != nil {
					b.Fatalf("匯出失敗: %v", err)
				}
				
				elapsed := time.Since(start)
				
				// 驗證效能要求
				if elapsed > bm.expected {
					b.Errorf("匯出 %d 筆記錄耗時 %v，超過預期的 %v", 
						bm.records, elapsed, bm.expected)
				}
				
				// 報告每秒處理記錄數
				recordsPerSec := float64(bm.records) / elapsed.Seconds()
				b.ReportMetric(recordsPerSec, "records/sec")
			}
		})
	}
}

// BenchmarkXLSXExportMemory 測試記憶體使用效率
// 確保記憶體使用在合理範圍內
func BenchmarkXLSXExportMemory(b *testing.B) {
	records := 100000 // 10萬筆記錄
	logs := createBenchmarkLogEntries(records)
	stats := createBenchmarkStatistics(records)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		tempFile := filepath.Join(b.TempDir(), "memory_test.xlsx")
		
		exporter := NewXLSXExporter()
		_, err := exporter.Export(logs, stats, tempFile)
		if err != nil {
			b.Fatalf("匯出失敗: %v", err)
		}
	}
}

// BenchmarkStreamingVsNormal 比較串流寫入與一般寫入的效能
func BenchmarkStreamingVsNormal(b *testing.B) {
	records := 50000
	logs := createBenchmarkLogEntries(records)
	stats := createBenchmarkStatistics(records)
	
	b.Run("Streaming", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tempFile := filepath.Join(b.TempDir(), "streaming.xlsx")
			exporter := NewXLSXExporter()
			exporter.SetStreamingMode(true) // 啟用串流模式
			
			_, err := exporter.Export(logs, stats, tempFile)
			if err != nil {
				b.Fatalf("串流匯出失敗: %v", err)
			}
		}
	})
	
	b.Run("Normal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tempFile := filepath.Join(b.TempDir(), "normal.xlsx")
			exporter := NewXLSXExporter()
			exporter.SetStreamingMode(false) // 一般模式
			
			_, err := exporter.Export(logs, stats, tempFile)
			if err != nil {
				b.Fatalf("一般匯出失敗: %v", err)
			}
		}
	})
}

// BenchmarkWorksheetCreation 測試各工作表創建的效能
func BenchmarkWorksheetCreation(b *testing.B) {
	records := 10000
	logs := createBenchmarkLogEntries(records)
	stats := createBenchmarkStatistics(records)
	
	exporter := NewXLSXExporter()
	
	b.Run("LogEntries", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tempFile := filepath.Join(b.TempDir(), "logs_only.xlsx")
			err := exporter.createLogEntriesWorksheet(logs, tempFile)
			if err != nil {
				b.Fatalf("日誌工作表創建失敗: %v", err)
			}
		}
	})
	
	b.Run("Statistics", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tempFile := filepath.Join(b.TempDir(), "stats_only.xlsx")
			err := exporter.createStatisticsWorksheet(stats, tempFile)
			if err != nil {
				b.Fatalf("統計工作表創建失敗: %v", err)
			}
		}
	})
	
	b.Run("BotDetection", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tempFile := filepath.Join(b.TempDir(), "bots_only.xlsx")
			err := exporter.createBotDetectionWorksheet(logs, tempFile)
			if err != nil {
				b.Fatalf("機器人偵測工作表創建失敗: %v", err)
			}
		}
	})
}

// BenchmarkDataFormatting 測試資料格式化的效能
func BenchmarkDataFormatting(b *testing.B) {
	records := 100000
	logs := createBenchmarkLogEntries(records)
	
	formatter := NewFormatter()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		formatted := formatter.FormatLogEntries(logs)
		if len(formatted) != records+1 { // +1 for header
			b.Fatalf("格式化記錄數量不正確: 期望 %d，實際 %d", records+1, len(formatted))
		}
	}
}

// Helper functions for benchmark data

// createBenchmarkLogEntries 創建基準測試用的大量日誌條目
func createBenchmarkLogEntries(count int) []*models.LogEntry {
	logs := make([]*models.LogEntry, count)
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	
	ips := []string{
		"192.168.1.100", "192.168.1.101", "192.168.1.102",
		"10.0.0.100", "10.0.0.101", "203.0.113.1",
	}
	
	methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD"}
	urls := []string{
		"/", "/index.html", "/api/users", "/api/login",
		"/assets/style.css", "/assets/script.js", "/favicon.ico",
	}
	
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"curl/7.68.0",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"facebookexternalhit/1.1",
	}
	
	statusCodes := []int{200, 201, 301, 302, 400, 401, 404, 500}
	
	for i := 0; i < count; i++ {
		logs[i] = &models.LogEntry{
			IP:            ips[i%len(ips)],
			Timestamp:     baseTime.Add(time.Duration(i) * time.Second),
			Method:        methods[i%len(methods)],
			URL:           urls[i%len(urls)],
			Protocol:      "HTTP/1.1",
			StatusCode:    statusCodes[i%len(statusCodes)],
			ResponseBytes: int64(1024 + (i * 100)),
			Referer:       "https://example.com",
			UserAgent:     userAgents[i%len(userAgents)],
			LineNumber:    i + 1,
		}
	}
	
	return logs
}

// createBenchmarkStatistics 創建基準測試用的統計資料
func createBenchmarkStatistics(logCount int) *models.Statistics {
	stats := models.NewStatistics()
	
	stats.TotalRequests = int64(logCount)
	stats.UniqueIPs = 6 // 基於測試資料中的 IP 數量
	stats.TotalBytes = int64(logCount * 1124) // 大約每筆記錄 1KB
	stats.StartTime = "2024-01-01 00:00:00"
	stats.EndTime = time.Date(2024, 1, 1, 0, 0, int(logCount), 0, time.UTC).Format("2006-01-02 15:04:05")
	
	// 狀態碼分布（基於測試資料的模式）
	stats.StatusCodeDist = map[int]int64{
		200: int64(logCount / 8 * 3), // 約 37.5%
		201: int64(logCount / 8),     // 約 12.5%
		301: int64(logCount / 8),     // 約 12.5%
		302: int64(logCount / 8),     // 約 12.5%
		400: int64(logCount / 8),     // 約 12.5%
		401: int64(logCount / 8),     // 約 12.5%
		404: int64(logCount / 8),     // 約 12.5%
		500: int64(logCount / 8),     // 約 12.5%
	}
	
	// 狀態類別分布
	stats.StatusCatDist = map[string]int64{
		"2xx": int64(logCount / 2),   // 50%
		"3xx": int64(logCount / 4),   // 25%
		"4xx": int64(logCount / 8),   // 12.5%
		"5xx": int64(logCount / 8),   // 12.5%
	}
	
	// Top IPs（簡化版）
	stats.TopIPs = []models.IPStat{
		{IP: "192.168.1.100", Count: int64(logCount / 6), TotalBytes: int64(logCount / 6 * 1124)},
		{IP: "192.168.1.101", Count: int64(logCount / 6), TotalBytes: int64(logCount / 6 * 1124)},
		{IP: "192.168.1.102", Count: int64(logCount / 6), TotalBytes: int64(logCount / 6 * 1124)},
		{IP: "10.0.0.100", Count: int64(logCount / 6), TotalBytes: int64(logCount / 6 * 1124)},
		{IP: "10.0.0.101", Count: int64(logCount / 6), TotalBytes: int64(logCount / 6 * 1124)},
		{IP: "203.0.113.1", Count: int64(logCount / 6), TotalBytes: int64(logCount / 6 * 1124)},
	}
	
	// Top URLs（簡化版）
	stats.TopURLs = []models.URLStat{
		{URL: "/", Count: int64(logCount / 7 * 2), TotalBytes: int64(logCount / 7 * 2 * 1124)},
		{URL: "/index.html", Count: int64(logCount / 7), TotalBytes: int64(logCount / 7 * 1124)},
		{URL: "/api/users", Count: int64(logCount / 7), TotalBytes: int64(logCount / 7 * 1124)},
		{URL: "/api/login", Count: int64(logCount / 7), TotalBytes: int64(logCount / 7 * 1124)},
		{URL: "/assets/style.css", Count: int64(logCount / 7), TotalBytes: int64(logCount / 7 * 1124)},
		{URL: "/assets/script.js", Count: int64(logCount / 7), TotalBytes: int64(logCount / 7 * 1124)},
		{URL: "/favicon.ico", Count: int64(logCount / 7), TotalBytes: int64(logCount / 7 * 1124)},
	}
	
	// 錯誤統計
	stats.ErrorCount = int64(logCount / 4)      // 25% 錯誤率
	stats.ClientErrorCount = int64(logCount / 8) // 12.5%
	stats.ServerErrorCount = int64(logCount / 8) // 12.5%
	stats.ErrorRate = 25.0
	
	return stats
}