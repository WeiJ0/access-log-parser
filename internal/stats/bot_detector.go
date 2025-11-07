package stats

import (
	"strings"
	"sync"
)

// BotDetector 提供機器人 User-Agent 的偵測功能
// 使用關鍵字匹配來識別常見的機器人、爬蟲和自動化工具
type BotDetector struct {
	patterns map[string]string // 關鍵字 -> 機器人類型的映射
	mu       sync.RWMutex      // 保護統計數據的互斥鎖
	stats    BotStats          // 統計資訊
}

// BotStats 儲存機器人偵測的統計資訊
type BotStats struct {
	Total          int            `json:"total"`          // 總請求數
	BotRequests    int            `json:"botRequests"`    // 機器人請求數
	HumanRequests  int            `json:"humanRequests"`  // 人類請求數
	BotPercentage  float64        `json:"botPercentage"`  // 機器人請求百分比
	BotTypes       map[string]int `json:"botTypes"`       // 各類型機器人的數量
}

// NewBotDetector 建立新的機器人偵測器
func NewBotDetector() *BotDetector {
	detector := &BotDetector{
		patterns: make(map[string]string),
		stats: BotStats{
			BotTypes: make(map[string]int),
		},
	}

	// 初始化機器人模式
	detector.initPatterns()

	return detector
}

// initPatterns 初始化機器人偵測模式
func (d *BotDetector) initPatterns() {
	// 搜尋引擎機器人
	searchEngines := []string{
		"googlebot", "bingbot", "slurp", "duckduckbot",
		"baiduspider", "yandexbot", "sogou", "exabot",
		"facebot", "ia_archiver",
	}
	for _, pattern := range searchEngines {
		d.patterns[pattern] = "搜尋引擎"
	}

	// 爬蟲和抓取工具
	crawlers := []string{
		"bot", "crawler", "spider", "scraper", "scraping",
		"python-requests", "curl", "wget", "httpclient",
		"scrapy", "beautifulsoup", "mechanize", "pycurl",
		"libwww", "okhttp", "go-http-client",
	}
	for _, pattern := range crawlers {
		d.patterns[pattern] = "爬蟲"
	}

	// 社交媒體機器人
	socialMedia := []string{
		"facebookexternalhit", "twitterbot", "linkedinbot",
		"pinterest", "slackbot", "telegrambot", "whatsapp",
		"discordbot",
	}
	for _, pattern := range socialMedia {
		d.patterns[pattern] = "社交媒體"
	}

	// 監控和正常運行時間檢查工具
	monitoring := []string{
		"pingdom", "uptimerobot", "statuscake", "monitor",
		"site24x7", "newrelic", "datadog", "nagios",
	}
	for _, pattern := range monitoring {
		d.patterns[pattern] = "監控工具"
	}

	// SEO 和分析工具
	seo := []string{
		"semrush", "ahrefs", "mj12bot", "majestic",
		"screaming frog", "seokicks", "seoscan",
	}
	for _, pattern := range seo {
		d.patterns[pattern] = "SEO 工具"
	}

	// 安全和掃描工具
	security := []string{
		"nessus", "nikto", "nmap", "masscan", "acunetix",
		"qualys", "securityscanner", "vulnscanner",
	}
	for _, pattern := range security {
		d.patterns[pattern] = "安全掃描"
	}
}

// IsBot 判斷給定的 User-Agent 是否為機器人
// 返回值: (是否為機器人, 機器人類型)
func (d *BotDetector) IsBot(userAgent string) (bool, string) {
	// 空字串或無效值不是機器人
	if userAgent == "" || userAgent == "-" {
		d.recordRequest(false, "")
		return false, ""
	}

	// 轉換為小寫進行不區分大小寫的匹配
	lowerUA := strings.ToLower(userAgent)

	// 優先順序匹配：先檢查特定關鍵字，後檢查通用關鍵字
	// 1. 先檢查搜尋引擎的特定名稱
	searchEnginePatterns := []string{
		"googlebot", "bingbot", "slurp", "duckduckbot",
		"baiduspider", "yandexbot", "yandex", "sogou", "exabot",
	}
	for _, pattern := range searchEnginePatterns {
		if strings.Contains(lowerUA, pattern) {
			d.recordRequest(true, "搜尋引擎")
			return true, "搜尋引擎"
		}
	}

	// 2. 檢查社交媒體機器人
	socialMediaPatterns := []string{
		"facebookexternalhit", "twitterbot", "linkedinbot",
		"pinterest", "slackbot", "telegrambot", "whatsapp",
		"discordbot",
	}
	for _, pattern := range socialMediaPatterns {
		if strings.Contains(lowerUA, pattern) {
			d.recordRequest(true, "社交媒體")
			return true, "社交媒體"
		}
	}

	// 3. 檢查監控工具
	monitoringPatterns := []string{
		"pingdom", "uptimerobot", "statuscake", "monitor",
		"site24x7", "newrelic", "datadog", "nagios",
	}
	for _, pattern := range monitoringPatterns {
		if strings.Contains(lowerUA, pattern) {
			d.recordRequest(true, "監控工具")
			return true, "監控工具"
		}
	}

	// 4. 檢查 SEO 工具
	seoPatterns := []string{
		"semrush", "ahrefs", "mj12bot", "majestic",
		"screaming frog", "seokicks", "seoscan",
	}
	for _, pattern := range seoPatterns {
		if strings.Contains(lowerUA, pattern) {
			d.recordRequest(true, "SEO 工具")
			return true, "SEO 工具"
		}
	}

	// 5. 檢查安全掃描工具
	securityPatterns := []string{
		"nessus", "nikto", "nmap", "masscan", "acunetix",
		"qualys", "securityscanner", "vulnscanner",
	}
	for _, pattern := range securityPatterns {
		if strings.Contains(lowerUA, pattern) {
			d.recordRequest(true, "安全掃描")
			return true, "安全掃描"
		}
	}

	// 6. 最後檢查通用爬蟲關鍵字（這些關鍵字較通用，放在最後）
	crawlerPatterns := []string{
		"bot", "crawler", "spider", "scraper", "scraping",
		"python-requests", "curl", "wget", "httpclient",
		"scrapy", "beautifulsoup", "mechanize", "pycurl",
		"libwww", "okhttp", "go-http-client",
	}
	for _, pattern := range crawlerPatterns {
		if strings.Contains(lowerUA, pattern) {
			d.recordRequest(true, "爬蟲")
			return true, "爬蟲"
		}
	}

	// 沒有匹配到任何模式，判定為人類
	d.recordRequest(false, "")
	return false, ""
}

// recordRequest 記錄請求統計
func (d *BotDetector) recordRequest(isBot bool, botType string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.stats.Total++

	if isBot {
		d.stats.BotRequests++
		d.stats.BotTypes[botType]++
	} else {
		d.stats.HumanRequests++
	}

	// 更新百分比
	if d.stats.Total > 0 {
		d.stats.BotPercentage = float64(d.stats.BotRequests) / float64(d.stats.Total) * 100
	}
}

// GetStats 獲取當前的統計資訊
func (d *BotDetector) GetStats() BotStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 複製統計資訊以避免競態條件
	statsCopy := BotStats{
		Total:         d.stats.Total,
		BotRequests:   d.stats.BotRequests,
		HumanRequests: d.stats.HumanRequests,
		BotPercentage: d.stats.BotPercentage,
		BotTypes:      make(map[string]int),
	}

	// 複製機器人類型計數
	for k, v := range d.stats.BotTypes {
		statsCopy.BotTypes[k] = v
	}

	return statsCopy
}

// ResetStats 重置統計資訊
func (d *BotDetector) ResetStats() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.stats = BotStats{
		BotTypes: make(map[string]int),
	}
}

// AddPattern 添加自訂的機器人偵測模式
// pattern: 要匹配的關鍵字（不區分大小寫）
// botType: 機器人類型
func (d *BotDetector) AddPattern(pattern, botType string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.patterns[strings.ToLower(pattern)] = botType
}

// RemovePattern 移除機器人偵測模式
func (d *BotDetector) RemovePattern(pattern string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.patterns, strings.ToLower(pattern))
}
