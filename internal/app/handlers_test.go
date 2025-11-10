package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"access-log-analyzer/internal/stats"
)

// TestParseFileWithStatistics 測試 ParseFile API 整合統計計算功能（T079）
func TestParseFileWithStatistics(t *testing.T) {
	// 建立測試用的 log 檔案
	testLog := `127.0.0.1 - - [01/Jan/2024:00:00:00 +0000] "GET /index.html HTTP/1.1" 200 1024 "-" "Mozilla/5.0"
192.168.1.100 - - [01/Jan/2024:00:00:01 +0000] "GET /api/users HTTP/1.1" 200 2048 "-" "Mozilla/5.0"
192.168.1.100 - - [01/Jan/2024:00:00:02 +0000] "POST /api/login HTTP/1.1" 200 512 "-" "Mozilla/5.0"
10.0.0.1 - - [01/Jan/2024:00:00:03 +0000] "GET /index.html HTTP/1.1" 404 256 "-" "Googlebot/2.1"
127.0.0.1 - - [01/Jan/2024:00:00:04 +0000] "GET /about.html HTTP/1.1" 200 768 "-" "Mozilla/5.0"
`

	// 建立暫存測試檔案
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.log")
	err := os.WriteFile(testFile, []byte(testLog), 0644)
	require.NoError(t, err)

	// 建立 App 實例
	app := NewApp()

	// 呼叫 ParseFile API
	req := ParseFileRequest{
		FilePath: testFile,
	}

	resp := app.ParseFile(req)

	// 驗證基本解析結果
	assert.True(t, resp.Success, "ParseFile should succeed")
	assert.Empty(t, resp.ErrorMessage, "Should have no error message")
	require.NotNil(t, resp.LogFile, "LogFile should not be nil")

	// 驗證解析統計
	assert.Equal(t, 5, resp.LogFile.TotalLines, "Should parse 5 lines")
	assert.Equal(t, 5, resp.LogFile.ParsedLines, "Should successfully parse 5 lines")
	assert.Equal(t, 0, resp.LogFile.ErrorLines, "Should have no errors")

	// 驗證統計資訊存在
	assert.NotNil(t, resp.LogFile.Statistics, "Statistics should not be nil")
	assert.GreaterOrEqual(t, resp.LogFile.StatTime, int64(0), "StatTime should be >= 0")

	// 將 Statistics 轉換為實際類型
	statistics, ok := resp.LogFile.Statistics.(stats.Statistics)
	if !ok {
		// 也可能是指標類型
		statisticsPtr, ok := resp.LogFile.Statistics.(*stats.Statistics)
		require.True(t, ok, "Statistics should be stats.Statistics or *stats.Statistics type")
		statistics = *statisticsPtr
	}

	// 驗證統計結果
	assert.Equal(t, 5, statistics.TotalRequests, "Should have 5 total requests")
	assert.Equal(t, 3, statistics.UniqueIPs, "Should have 3 unique IPs")
	assert.Equal(t, int64(4608), statistics.TotalBytes, "Total bytes should be 1024+2048+512+256+768")

	// 驗證 Top IPs（至少前 3 個）
	require.GreaterOrEqual(t, len(statistics.TopIPs), 2, "Should have at least top 2 IPs")
	// 可能是 127.0.0.1 或 192.168.1.100（都有 2 次請求）
	foundIP1 := false
	foundIP2 := false
	for _, ip := range statistics.TopIPs {
		if ip.IP == "127.0.0.1" && ip.RequestCount == 2 {
			foundIP1 = true
		}
		if ip.IP == "192.168.1.100" && ip.RequestCount == 2 {
			foundIP2 = true
		}
	}
	assert.True(t, foundIP1 || foundIP2, "Should have 127.0.0.1 or 192.168.1.100 in top IPs")

	// 驗證 Top Paths
	require.GreaterOrEqual(t, len(statistics.TopPaths), 1, "Should have at least top 1 path")
	assert.Equal(t, "/index.html", statistics.TopPaths[0].Path, "Top path should be /index.html")
	assert.Equal(t, 2, statistics.TopPaths[0].RequestCount, "/index.html should have 2 requests")

	// 驗證狀態碼分布
	assert.Equal(t, 4, statistics.StatusCodeDistribution.Success, "Should have 4 success (2xx) responses")
	assert.Equal(t, 1, statistics.StatusCodeDistribution.ClientError, "Should have 1 client error (4xx)")
	assert.Equal(t, 4, statistics.StatusCodeDistribution.Details[200], "Should have 4 status 200")
	assert.Equal(t, 1, statistics.StatusCodeDistribution.Details[404], "Should have 1 status 404")

	// 驗證機器人偵測
	assert.Equal(t, 1, statistics.BotStats.BotRequests, "Should have 1 bot request")
	assert.InDelta(t, 20.0, statistics.BotStats.BotPercentage, 0.1, "Bot percentage should be ~20%")
}

// TestParseFilePerformance 測試統計計算效能（T071）
func TestParseFilePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// 建立較大的測試檔案（10000 行）
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large.log")

	f, err := os.Create(testFile)
	require.NoError(t, err)
	defer f.Close()

	logLine := `127.0.0.1 - - [01/Jan/2024:00:00:00 +0000] "GET /index.html HTTP/1.1" 200 1024 "-" "Mozilla/5.0"` + "\n"
	for i := 0; i < 10000; i++ {
		_, err := f.WriteString(logLine)
		require.NoError(t, err)
	}
	f.Close()

	// 建立 App 實例
	app := NewApp()

	// 呼叫 ParseFile API
	req := ParseFileRequest{
		FilePath: testFile,
	}

	resp := app.ParseFile(req)

	// 驗證成功
	assert.True(t, resp.Success, "ParseFile should succeed")
	require.NotNil(t, resp.LogFile, "LogFile should not be nil")

	// 驗證效能指標
	assert.Greater(t, resp.LogFile.ParseTime, int64(0), "ParseTime should be > 0")
	assert.Greater(t, resp.LogFile.StatTime, int64(0), "StatTime should be > 0")

	// 統計耗時應該遠小於解析耗時（通常 < 10%）
	t.Logf("ParseTime: %d ms, StatTime: %d ms", resp.LogFile.ParseTime, resp.LogFile.StatTime)

	// 統計耗時應該合理（10K 行應該 < 100ms）
	assert.Less(t, resp.LogFile.StatTime, int64(100), "StatTime should be < 100ms for 10K records")
}
