package testing

import (
	"runtime"
	"testing"
	"time"

	"access-log-analyzer/internal/models"
)

// AssertEqual 斷言兩個值相等
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: 期望 %v, 實際 %v", msg, expected, actual)
	}
}

// AssertNotEqual 斷言兩個值不相等
func AssertNotEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected == actual {
		t.Errorf("%s: 不應該相等, 實際值 %v", msg, actual)
	}
}

// AssertTrue 斷言條件為真
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: 條件應為 true", msg)
	}
}

// AssertFalse 斷言條件為假
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("%s: 條件應為 false", msg)
	}
}

// AssertNil 斷言值為 nil
func AssertNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: 應為 nil, 實際 %v", msg, value)
	}
}

// AssertNotNil 斷言值不為 nil
func AssertNotNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: 不應為 nil", msg)
	}
}

// AssertError 斷言有錯誤發生
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: 應該有錯誤", msg)
	}
}

// AssertNoError 斷言沒有錯誤
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: 不應該有錯誤, 錯誤: %v", msg, err)
	}
}

// AssertGreaterThan 斷言 a > b
func AssertGreaterThan(t *testing.T, a, b float64, msg string) {
	t.Helper()
	if a <= b {
		t.Errorf("%s: %v 應該大於 %v", msg, a, b)
	}
}

// AssertLessThan 斷言 a < b
func AssertLessThan(t *testing.T, a, b float64, msg string) {
	t.Helper()
	if a >= b {
		t.Errorf("%s: %v 應該小於 %v", msg, a, b)
	}
}

// AssertInRange 斷言值在範圍內 [min, max]
func AssertInRange(t *testing.T, value, min, max float64, msg string) {
	t.Helper()
	if value < min || value > max {
		t.Errorf("%s: %v 應該在 [%v, %v] 範圍內", msg, value, min, max)
	}
}

// CreateTestLogEntry 建立測試用的 LogEntry
func CreateTestLogEntry(ip string, statusCode int, timestamp time.Time) *models.LogEntry {
	return &models.LogEntry{
		IP:            ip,
		Timestamp:     timestamp,
		Method:        "GET",
		URL:           "/test",
		Protocol:      "HTTP/1.1",
		StatusCode:    statusCode,
		ResponseBytes: 1024,
		Referer:       "https://example.com",
		UserAgent:     "TestAgent/1.0",
		LineNumber:    1,
		RawLine:       "test log line",
	}
}

// CreateTestStatistics 建立測試用的 Statistics
func CreateTestStatistics() *models.Statistics {
	stats := models.NewStatistics()
	stats.TotalRequests = 100
	stats.UniqueIPs = 50
	stats.TotalBytes = 102400
	stats.ErrorCount = 5
	stats.ClientErrorCount = 3
	stats.ServerErrorCount = 2
	stats.ErrorRate = 5.0
	return stats
}

// MeasurePerformance 測量函式執行時間
// 返回執行時間（毫秒）
func MeasurePerformance(fn func()) int64 {
	start := time.Now()
	fn()
	return time.Since(start).Milliseconds()
}

// MeasureMemory 測量函式記憶體使用
// 返回記憶體增長量（bytes）
func MeasureMemory(fn func()) uint64 {
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	fn()

	runtime.ReadMemStats(&memAfter)
	return memAfter.Alloc - memBefore.Alloc
}

// AssertPerformance 斷言效能符合要求
// maxDurationMs: 最大允許執行時間（毫秒）
func AssertPerformance(t *testing.T, fn func(), maxDurationMs int64, msg string) {
	t.Helper()
	duration := MeasurePerformance(fn)
	if duration > maxDurationMs {
		t.Errorf("%s: 執行時間 %dms 超過限制 %dms", msg, duration, maxDurationMs)
	}
}

// AssertMemoryUsage 斷言記憶體使用符合要求
// maxMemoryBytes: 最大允許記憶體使用（bytes）
func AssertMemoryUsage(t *testing.T, fn func(), maxMemoryBytes uint64, msg string) {
	t.Helper()
	memUsed := MeasureMemory(fn)
	if memUsed > maxMemoryBytes {
		t.Errorf("%s: 記憶體使用 %d bytes 超過限制 %d bytes", msg, memUsed, maxMemoryBytes)
	}
}
