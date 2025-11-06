package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkParseLine 基準測試單行解析效能
func BenchmarkParseLine(b *testing.B) {
	parser := NewParser(FormatCombined, 1)
	pattern := GetPattern(FormatCombined)
	line := `192.168.1.100 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.parseLine(1, line, pattern)
	}
}

// BenchmarkParseFile_Small 基準測試小檔案解析（1000 行）
func BenchmarkParseFile_Small(b *testing.B) {
	// 建立測試檔案
	tempFile := createTestLogFile(b, 1000)
	defer os.Remove(tempFile)
	
	fileInfo, _ := os.Stat(tempFile)
	parser := NewParser(FormatCombined, 4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseFile(tempFile, fileInfo.Size())
	}
}

// BenchmarkParseFile_Medium 基準測試中等檔案解析（10,000 行）
func BenchmarkParseFile_Medium(b *testing.B) {
	tempFile := createTestLogFile(b, 10000)
	defer os.Remove(tempFile)
	
	fileInfo, _ := os.Stat(tempFile)
	parser := NewParser(FormatCombined, 4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseFile(tempFile, fileInfo.Size())
	}
}

// BenchmarkParseFile_Large 基準測試大檔案解析（100,000 行）
func BenchmarkParseFile_Large(b *testing.B) {
	tempFile := createTestLogFile(b, 100000)
	defer os.Remove(tempFile)
	
	fileInfo, _ := os.Stat(tempFile)
	parser := NewParser(FormatCombined, 4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseFile(tempFile, fileInfo.Size())
	}
}

// BenchmarkParseFile_Workers 比較不同 worker 數量的效能
func BenchmarkParseFile_Workers(b *testing.B) {
	tempFile := createTestLogFile(b, 50000)
	defer os.Remove(tempFile)
	
	fileInfo, _ := os.Stat(tempFile)
	
	workerCounts := []int{1, 2, 4, 8, 16}
	for _, count := range workerCounts {
		b.Run(fmt.Sprintf("Workers_%d", count), func(b *testing.B) {
			parser := NewParser(FormatCombined, count)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = parser.ParseFile(tempFile, fileInfo.Size())
			}
		})
	}
}

// BenchmarkDetectFormat 基準測試格式偵測效能
func BenchmarkDetectFormat(b *testing.B) {
	line := `192.168.1.100 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DetectFormat(line)
	}
}

// BenchmarkValidateFormat 基準測試格式驗證效能
func BenchmarkValidateFormat(b *testing.B) {
	tempFile := createTestLogFile(b, 1000)
	defer os.Remove(tempFile)
	
	parser := NewParser(FormatCombined, 1)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ValidateFormat(tempFile, 100)
	}
}

// createTestLogFile 建立測試用 log 檔案
// 生成指定行數的標準 Combined 格式 log
func createTestLogFile(tb testing.TB, lines int) string {
	tb.Helper()
	
	tempDir := tb.TempDir()
	tempFile := filepath.Join(tempDir, "test.log")
	
	file, err := os.Create(tempFile)
	if err != nil {
		tb.Fatalf("無法建立測試檔案: %v", err)
	}
	defer file.Close()
	
	// 生成多樣化的 log 行
	templates := []string{
		`192.168.1.%d - - [06/Nov/2025:14:30:%02d +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"` + "\n",
		`10.0.0.%d - - [06/Nov/2025:14:30:%02d +0800] "POST /api/login HTTP/1.1" 200 512 "-" "curl/7.68.0"` + "\n",
		`172.16.0.%d - - [06/Nov/2025:14:30:%02d +0800] "GET /missing.html HTTP/1.1" 404 196 "https://example.com/index.html" "Mozilla/5.0"` + "\n",
		`192.168.100.%d - - [06/Nov/2025:14:30:%02d +0800] "GET /api/data HTTP/1.1" 200 4096 "https://app.example.com" "Mozilla/5.0"` + "\n",
		`203.0.113.%d - admin [06/Nov/2025:14:30:%02d +0800] "PUT /api/users/123 HTTP/1.1" 200 256 "-" "PostmanRuntime/7.28.0"` + "\n",
	}
	
	for i := 0; i < lines; i++ {
		template := templates[i%len(templates)]
		ip := (i % 254) + 1
		sec := (i % 60)
		line := fmt.Sprintf(template, ip, sec)
		
		if _, err := file.WriteString(line); err != nil {
			tb.Fatalf("無法寫入測試資料: %v", err)
		}
	}
	
	return tempFile
}

// TestThroughput 測試解析吞吐量是否達標（≥50 MB/秒）
func TestThroughput(t *testing.T) {
	if testing.Short() {
		t.Skip("跳過吞吐量測試（使用 -short 模式）")
	}
	
	// 建立大約 10MB 的測試檔案（約 50,000 行）
	tempFile := createTestLogFile(t, 50000)
	defer os.Remove(tempFile)
	
	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("無法取得檔案資訊: %v", err)
	}
	
	parser := NewParser(FormatCombined, 0) // 使用所有 CPU 核心
	
	result, err := parser.ParseFile(tempFile, fileInfo.Size())
	if err != nil {
		t.Fatalf("解析失敗: %v", err)
	}
	
	t.Logf("解析統計:")
	t.Logf("  檔案大小: %.2f MB", float64(fileInfo.Size())/(1024*1024))
	t.Logf("  總行數: %d", result.TotalLines)
	t.Logf("  解析耗時: %v", result.ParseTime)
	t.Logf("  吞吐量: %.2f MB/秒", result.ThroughputMB)
	t.Logf("  記憶體使用: %.2f MB", float64(result.MemoryUsed)/(1024*1024))
	
	// 驗證吞吐量至少達到 50 MB/秒
	if result.ThroughputMB < 50.0 {
		t.Errorf("吞吐量未達標: %.2f MB/秒 < 50 MB/秒", result.ThroughputMB)
	}
	
	// 驗證記憶體使用不超過檔案大小的 5 倍（考慮測試環境開銷和結構體封裝）
	memoryRatio := float64(result.MemoryUsed) / float64(fileInfo.Size())
	if memoryRatio > 5.0 {
		t.Errorf("記憶體使用過高: %.2fx 檔案大小", memoryRatio)
	} else {
		t.Logf("  記憶體效率: %.2fx 檔案大小（合格）", memoryRatio)
	}
}
