package stats

import (
	"testing"
	"time"

	"access-log-analyzer/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculator_基本統計 測試基本統計計算功能
func TestCalculator_基本統計(t *testing.T) {
	calc := NewCalculator()

	// 建立測試資料
	entries := []models.LogEntry{
		{
			IP:            "192.168.1.1",
			Timestamp:     time.Now(),
			Method:        "GET",
			URL:           "/index.html",
			StatusCode:    200,
			ResponseBytes: 1024,
			UserAgent:     "Mozilla/5.0",
		},
		{
			IP:            "192.168.1.1",
			Timestamp:     time.Now(),
			Method:        "POST",
			URL:           "/api/users",
			StatusCode:    201,
			ResponseBytes: 512,
			UserAgent:     "Mozilla/5.0",
		},
		{
			IP:            "192.168.1.2",
			Timestamp:     time.Now(),
			Method:        "GET",
			URL:           "/index.html",
			StatusCode:    200,
			ResponseBytes: 1024,
			UserAgent:     "Googlebot/2.1",
		},
	}

	// 計算統計
	stats := calc.Calculate(entries)

	// 驗證基本計數
	assert.Equal(t, 3, stats.TotalRequests, "總請求數應該是 3")
	assert.Equal(t, 2, stats.UniqueIPs, "唯一 IP 數應該是 2")
	assert.Equal(t, 2, stats.UniquePaths, "唯一路徑數應該是 2")
}

// TestCalculator_TopIPs 測試 Top IP 統計
func TestCalculator_TopIPs(t *testing.T) {
	calc := NewCalculator()

	// 建立測試資料，IP 1 有 3 個請求，IP 2 有 2 個請求，IP 3 有 1 個請求
	entries := []models.LogEntry{
		{IP: "192.168.1.1"},
		{IP: "192.168.1.1"},
		{IP: "192.168.1.1"},
		{IP: "192.168.1.2"},
		{IP: "192.168.1.2"},
		{IP: "192.168.1.3"},
	}

	stats := calc.Calculate(entries)

	// 驗證 Top IPs
	require.Len(t, stats.TopIPs, 3, "應該有 3 個 IP")
	assert.Equal(t, "192.168.1.1", stats.TopIPs[0].IP, "第一名應該是 192.168.1.1")
	assert.Equal(t, 3, stats.TopIPs[0].RequestCount, "192.168.1.1 應該有 3 個請求")
	assert.Equal(t, "192.168.1.2", stats.TopIPs[1].IP, "第二名應該是 192.168.1.2")
	assert.Equal(t, 2, stats.TopIPs[1].RequestCount, "192.168.1.2 應該有 2 個請求")
}

// TestCalculator_TopPaths 測試 Top 路徑統計
func TestCalculator_TopPaths(t *testing.T) {
	calc := NewCalculator()

	entries := []models.LogEntry{
		{URL: "/index.html", ResponseBytes: 1000, StatusCode: 200},
		{URL: "/index.html", ResponseBytes: 1200, StatusCode: 200},
		{URL: "/about.html", ResponseBytes: 800, StatusCode: 200},
		{URL: "/index.html", ResponseBytes: 1100, StatusCode: 404},
	}

	stats := calc.Calculate(entries)

	// 驗證 Top Paths
	require.Len(t, stats.TopPaths, 2, "應該有 2 個路徑")

	// 驗證 /index.html 的統計
	indexPath := stats.TopPaths[0]
	assert.Equal(t, "/index.html", indexPath.Path)
	assert.Equal(t, 3, indexPath.RequestCount, "/index.html 應該有 3 個請求")
	assert.Equal(t, int64(1100), indexPath.AverageSize, "平均大小應該是 (1000+1200+1100)/3")
	assert.InDelta(t, 33.33, indexPath.ErrorRate, 0.01, "錯誤率應該是 1/3 = 33.33%")
}

// TestCalculator_StatusCodes 測試狀態碼分布統計
func TestCalculator_StatusCodes(t *testing.T) {
	calc := NewCalculator()

	entries := []models.LogEntry{
		{StatusCode: 200},
		{StatusCode: 200},
		{StatusCode: 201},
		{StatusCode: 301},
		{StatusCode: 404},
		{StatusCode: 500},
	}

	stats := calc.Calculate(entries)

	// 驗證狀態碼分布
	statusCodes := stats.StatusCodeDistribution
	assert.Equal(t, 3, statusCodes.Success, "成功請求 (2xx) 應該有 3 個")
	assert.Equal(t, 1, statusCodes.Redirection, "重定向 (3xx) 應該有 1 個")
	assert.Equal(t, 1, statusCodes.ClientError, "客戶端錯誤 (4xx) 應該有 1 個")
	assert.Equal(t, 1, statusCodes.ServerError, "伺服器錯誤 (5xx) 應該有 1 個")
}

// TestCalculator_BotDetection 測試機器人偵測整合
func TestCalculator_BotDetection(t *testing.T) {
	calc := NewCalculator()

	entries := []models.LogEntry{
		{UserAgent: "Googlebot/2.1"},
		{UserAgent: "Mozilla/5.0 Chrome/120.0.0.0"},
		{UserAgent: "bingbot/2.0"},
		{UserAgent: "Mozilla/5.0 Firefox/121.0"},
	}

	stats := calc.Calculate(entries)

	// 驗證機器人統計
	assert.Equal(t, 4, stats.BotStats.Total, "總請求應該是 4")
	assert.Equal(t, 2, stats.BotStats.BotRequests, "機器人請求應該是 2")
	assert.Equal(t, 2, stats.BotStats.HumanRequests, "人類請求應該是 2")
	assert.InDelta(t, 50.0, stats.BotStats.BotPercentage, 0.01, "機器人百分比應該是 50%")
}

// TestCalculator_空資料 測試空資料集
func TestCalculator_空資料(t *testing.T) {
	calc := NewCalculator()

	stats := calc.Calculate([]models.LogEntry{})

	assert.Equal(t, 0, stats.TotalRequests)
	assert.Equal(t, 0, stats.UniqueIPs)
	assert.Equal(t, 0, stats.UniquePaths)
	assert.Empty(t, stats.TopIPs)
	assert.Empty(t, stats.TopPaths)
}

// TestCalculator_大量資料 測試大量資料的統計計算
func TestCalculator_大量資料(t *testing.T) {
	calc := NewCalculator()

	// 產生 10 萬筆測試資料
	entries := make([]models.LogEntry, 100000)
	for i := 0; i < 100000; i++ {
		entries[i] = models.LogEntry{
			IP:            "192.168.1." + string(rune(i%256)),
			URL:           "/" + string(rune(i%100)),
			StatusCode:    200 + (i % 5),
			ResponseBytes: int64(1000 + (i % 5000)),
			UserAgent:     "Mozilla/5.0",
		}
	}

	startTime := time.Now()
	stats := calc.Calculate(entries)
	duration := time.Since(startTime)

	// 驗證計算完成
	assert.Equal(t, 100000, stats.TotalRequests)

	// 驗證效能（100 萬筆應該在 10 秒內完成，10 萬筆應該在 1 秒內）
	assert.Less(t, duration.Seconds(), 1.0, "10 萬筆資料應該在 1 秒內完成計算")

	t.Logf("100000 筆資料統計耗時: %.2f 秒", duration.Seconds())
}

// TestCalculator_TopN限制 測試 Top-N 數量限制
func TestCalculator_TopN限制(t *testing.T) {
	calc := NewCalculator()

	// 產生 20 個不同的 IP
	entries := make([]models.LogEntry, 20)
	for i := 0; i < 20; i++ {
		entries[i] = models.LogEntry{
			IP:  "192.168.1." + string(rune(i)),
			URL: "/" + string(rune(i)),
		}
	}

	stats := calc.Calculate(entries)

	// 驗證只保留 Top 10
	assert.LessOrEqual(t, len(stats.TopIPs), 10, "TopIPs 應該最多 10 個")
	assert.LessOrEqual(t, len(stats.TopPaths), 10, "TopPaths 應該最多 10 個")
}

// BenchmarkCalculator_小資料集 測試小資料集的計算效能
func BenchmarkCalculator_小資料集(b *testing.B) {
	calc := NewCalculator()

	entries := make([]models.LogEntry, 1000)
	for i := 0; i < 1000; i++ {
		entries[i] = models.LogEntry{
			IP:            "192.168.1.1",
			URL:           "/index.html",
			StatusCode:    200,
			ResponseBytes: 1024,
			UserAgent:     "Mozilla/5.0",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(entries)
	}
}

// BenchmarkCalculator_大資料集 測試大資料集的計算效能
func BenchmarkCalculator_大資料集(b *testing.B) {
	calc := NewCalculator()

	entries := make([]models.LogEntry, 1000000)
	for i := 0; i < 1000000; i++ {
		entries[i] = models.LogEntry{
			IP:            "192.168.1." + string(rune(i%256)),
			URL:           "/" + string(rune(i%100)),
			StatusCode:    200 + (i % 5),
			ResponseBytes: int64(1000 + (i % 5000)),
			UserAgent:     "Mozilla/5.0",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(entries)
	}
}
