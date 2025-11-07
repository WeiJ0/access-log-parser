package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBotDetector_常見機器人 測試常見機器人的偵測
func TestBotDetector_常見機器人(t *testing.T) {
	detector := NewBotDetector()

	testCases := []struct {
		userAgent string
		isBot     bool
		botType   string
	}{
		// 搜尋引擎機器人
		{"Googlebot/2.1 (+http://www.google.com/bot.html)", true, "搜尋引擎"},
		{"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", true, "搜尋引擎"},
		{"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)", true, "搜尋引擎"},
		{"DuckDuckBot/1.0; (+http://duckduckgo.com/duckduckbot.html)", true, "搜尋引擎"},
		{"Baiduspider+(+http://www.baidu.com/search/spider.htm)", true, "搜尋引擎"},
		{"Yandex/1.01.001 (compatible; Win16; I)", true, "搜尋引擎"},

		// 爬蟲
		{"Scrapy/2.5.0 (+https://scrapy.org)", true, "爬蟲"},
		{"python-requests/2.28.0", true, "爬蟲"},
		{"Apache-HttpClient/4.5.13", true, "爬蟲"},
		{"curl/7.68.0", true, "爬蟲"},
		{"wget/1.20.3", true, "爬蟲"},

		// 社交媒體機器人
		{"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)", true, "社交媒體"},
		{"Twitterbot/1.0", true, "社交媒體"},
		{"LinkedInBot/1.0 (compatible; Mozilla/5.0; Apache-HttpClient +http://www.linkedin.com)", true, "社交媒體"},
		{"Slackbot-LinkExpanding 1.0 (+https://api.slack.com/robots)", true, "社交媒體"},

		// 監控工具
		{"Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)", true, "監控工具"},
		{"UptimeRobot/2.0; http://www.uptimerobot.com/", true, "監控工具"},
		{"StatusCake/1.0 (+http://www.statuscake.com/)", true, "監控工具"},

		// SEO 工具
		{"SemrushBot/7~bl", true, "SEO 工具"},
		{"AhrefsBot/7.0", true, "SEO 工具"},
		{"MJ12bot/v1.4.8", true, "SEO 工具"},

		// 正常瀏覽器（應該不被標記為機器人）
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", false, ""},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", false, ""},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0", false, ""},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15", false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.userAgent[:min(30, len(tc.userAgent))], func(t *testing.T) {
			isBot, botType := detector.IsBot(tc.userAgent)

			assert.Equal(t, tc.isBot, isBot,
				"User-Agent: %s 的機器人判定錯誤", tc.userAgent)

			if tc.isBot {
				assert.Equal(t, tc.botType, botType,
					"User-Agent: %s 的機器人類型錯誤", tc.userAgent)
			}
		})
	}
}

// TestBotDetector_空字串 測試空字串和無效輸入
func TestBotDetector_空字串(t *testing.T) {
	detector := NewBotDetector()

	isBot, botType := detector.IsBot("")
	assert.False(t, isBot, "空字串不應被判定為機器人")
	assert.Empty(t, botType, "空字串應該返回空的機器人類型")

	isBot, botType = detector.IsBot("-")
	assert.False(t, isBot, "- 符號不應被判定為機器人")
	assert.Empty(t, botType)
}

// TestBotDetector_大小寫不敏感 測試大小寫不敏感的匹配
func TestBotDetector_大小寫不敏感(t *testing.T) {
	detector := NewBotDetector()

	testCases := []string{
		"googlebot",
		"GOOGLEBOT",
		"GoogleBot",
		"gOoGlEbOt",
	}

	for _, ua := range testCases {
		isBot, botType := detector.IsBot(ua)
		assert.True(t, isBot, "應該不區分大小寫: %s", ua)
		assert.Equal(t, "搜尋引擎", botType)
	}
}

// TestBotDetector_統計功能 測試統計功能
func TestBotDetector_統計功能(t *testing.T) {
	detector := NewBotDetector()

	userAgents := []string{
		"Googlebot/2.1",
		"Mozilla/5.0 Chrome/120.0.0.0",
		"bingbot/2.0",
		"Mozilla/5.0 Firefox/121.0",
		"Googlebot/2.1",
		"python-requests/2.28.0",
	}

	for _, ua := range userAgents {
		detector.IsBot(ua)
	}

	stats := detector.GetStats()

	// 驗證統計數據
	assert.Equal(t, 6, stats.Total, "總請求數應該是 6")
	assert.Equal(t, 4, stats.BotRequests, "機器人請求應該是 4")
	assert.Equal(t, 2, stats.HumanRequests, "人類請求應該是 2")

	// 驗證百分比
	expectedBotPercentage := float64(4) / float64(6) * 100
	assert.InDelta(t, expectedBotPercentage, stats.BotPercentage, 0.01,
		"機器人百分比計算錯誤")
}

// TestBotDetector_機器人類型統計 測試機器人類型統計
func TestBotDetector_機器人類型統計(t *testing.T) {
	detector := NewBotDetector()

	userAgents := []struct {
		ua       string
		expected string
	}{
		{"Googlebot/2.1", "搜尋引擎"},
		{"bingbot/2.0", "搜尋引擎"},
		{"python-requests/2.28.0", "爬蟲"},
		{"Scrapy/2.5.0", "爬蟲"},
		{"facebookexternalhit/1.1", "社交媒體"},
	}

	for _, tc := range userAgents {
		detector.IsBot(tc.ua)
	}

	stats := detector.GetStats()

	// 驗證機器人類型計數
	assert.Equal(t, 2, stats.BotTypes["搜尋引擎"], "搜尋引擎機器人應該有 2 個")
	assert.Equal(t, 2, stats.BotTypes["爬蟲"], "爬蟲應該有 2 個")
	assert.Equal(t, 1, stats.BotTypes["社交媒體"], "社交媒體機器人應該有 1 個")
}

// TestBotDetector_重置統計 測試重置統計功能
func TestBotDetector_重置統計(t *testing.T) {
	detector := NewBotDetector()

	// 添加一些數據
	detector.IsBot("Googlebot/2.1")
	detector.IsBot("Mozilla/5.0 Chrome/120.0.0.0")

	stats := detector.GetStats()
	assert.Equal(t, 2, stats.Total, "重置前應該有 2 筆數據")

	// 重置
	detector.ResetStats()

	stats = detector.GetStats()
	assert.Equal(t, 0, stats.Total, "重置後應該是 0")
	assert.Equal(t, 0, stats.BotRequests)
	assert.Equal(t, 0, stats.HumanRequests)
	assert.Empty(t, stats.BotTypes)
}

// BenchmarkBotDetector_偵測 測試機器人偵測的性能
func BenchmarkBotDetector_偵測(b *testing.B) {
	detector := NewBotDetector()
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.IsBot(userAgent)
	}
}

// BenchmarkBotDetector_混合請求 測試混合請求的性能
func BenchmarkBotDetector_混合請求(b *testing.B) {
	detector := NewBotDetector()
	userAgents := []string{
		"Googlebot/2.1",
		"Mozilla/5.0 Chrome/120.0.0.0",
		"bingbot/2.0",
		"python-requests/2.28.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ua := userAgents[i%len(userAgents)]
		detector.IsBot(ua)
	}
}

// 輔助函數
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
