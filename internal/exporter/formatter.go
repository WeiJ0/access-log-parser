package exporter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"access-log-analyzer/internal/models"
)

// Formatter 負責將 Go 資料結構轉換為 Excel 友善的格式
// 提供一致的資料格式化和表格結構
type Formatter struct {
	timeFormat string // 時間格式化字串
}

// NewFormatter 建立新的格式化器實例
func NewFormatter() *Formatter {
	return &Formatter{
		timeFormat: "2006-01-02 15:04:05", // 標準時間格式
	}
}

// FormatLogEntries 格式化日誌條目為二維字串陣列
// 返回包含標題行和資料行的二維陣列，適用於 Excel 匯出
func (f *Formatter) FormatLogEntries(logs []*models.LogEntry) [][]string {
	// 建立標題行
	headers := []string{
		"IP位址",
		"時間戳",
		"HTTP方法",
		"URL",
		"協定",
		"狀態碼",
		"回應大小",
		"來源頁面",
		"User Agent",
	}

	// 初始化結果陣列
	result := make([][]string, 0, len(logs)+1)
	result = append(result, headers)

	// 處理每筆日誌條目
	for _, log := range logs {
		// 跳過 nil 條目
		if log == nil {
			continue
		}

		row := []string{
			log.IP,
			f.formatTime(log.Timestamp),
			log.Method,
			log.URL,
			log.Protocol,
			strconv.Itoa(log.StatusCode),
			strconv.FormatInt(log.ResponseBytes, 10),
			f.formatReferer(log.Referer),
			log.UserAgent,
		}

		result = append(result, row)
	}

	return result
}

// FormatStatistics 格式化統計資料為二維字串陣列
// 創建包含各種統計指標的結構化表格
func (f *Formatter) FormatStatistics(stats *models.Statistics) [][]string {
	if stats == nil {
		return [][]string{{"統計項目", "數值"}}
	}

	result := make([][]string, 0)

	// 基本統計資料區塊
	result = append(result, []string{"===== 基本統計 ====="})
	result = append(result, []string{"統計項目", "數值"})
	result = append(result, []string{"總請求數", strconv.FormatInt(stats.TotalRequests, 10)})
	result = append(result, []string{"唯一IP數量", strconv.Itoa(stats.UniqueIPs)})
	result = append(result, []string{"總傳輸量 (位元組)", strconv.FormatInt(stats.TotalBytes, 10)})
	result = append(result, []string{"總傳輸量 (MB)", fmt.Sprintf("%.2f", stats.GetTotalBytesMB())})
	result = append(result, []string{"開始時間", stats.StartTime})
	result = append(result, []string{"結束時間", stats.EndTime})
	result = append(result, []string{"平均回應大小 (位元組)", fmt.Sprintf("%.2f", stats.AvgResponseSize)})
	result = append(result, []string{"最小回應大小", strconv.FormatInt(stats.MinResponseSize, 10)})
	result = append(result, []string{"最大回應大小", strconv.FormatInt(stats.MaxResponseSize, 10)})

	// 錯誤統計區塊
	result = append(result, []string{""}) // 空行分隔
	result = append(result, []string{"===== 錯誤統計 ====="})
	result = append(result, []string{"錯誤總數", strconv.FormatInt(stats.ErrorCount, 10)})
	result = append(result, []string{"客戶端錯誤 (4xx)", strconv.FormatInt(stats.ClientErrorCount, 10)})
	result = append(result, []string{"伺服器錯誤 (5xx)", strconv.FormatInt(stats.ServerErrorCount, 10)})
	result = append(result, []string{"錯誤率 (%)", fmt.Sprintf("%.2f", stats.ErrorRate)})

	// 狀態碼分布
	result = append(result, []string{""})
	result = append(result, []string{"===== 狀態碼分布 ====="})
	result = append(result, []string{"狀態碼", "次數"})

	// 排序狀態碼
	var statusCodes []int
	for code := range stats.StatusCodeDist {
		statusCodes = append(statusCodes, code)
	}
	sort.Ints(statusCodes)

	for _, code := range statusCodes {
		count := stats.StatusCodeDist[code]
		result = append(result, []string{
			strconv.Itoa(code),
			strconv.FormatInt(count, 10),
		})
	}

	// 狀態類別分布
	result = append(result, []string{""})
	result = append(result, []string{"===== 狀態類別分布 ====="})
	result = append(result, []string{"類別", "次數"})

	// 排序狀態類別
	var categories []string
	for cat := range stats.StatusCatDist {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		count := stats.StatusCatDist[cat]
		result = append(result, []string{cat, strconv.FormatInt(count, 10)})
	}

	// Top IPs
	result = append(result, []string{""})
	result = append(result, []string{"===== Top 10 IP 位址 ====="})
	result = append(result, []string{"IP位址", "請求次數", "總傳輸量", "錯誤次數", "錯誤率(%)", "唯一URL數"})

	limit := len(stats.TopIPs)
	if limit > 10 {
		limit = 10
	}

	for i := 0; i < limit; i++ {
		ip := stats.TopIPs[i]
		result = append(result, []string{
			ip.IP,
			strconv.FormatInt(ip.Count, 10),
			strconv.FormatInt(ip.TotalBytes, 10),
			strconv.FormatInt(ip.ErrorCount, 10),
			fmt.Sprintf("%.2f", ip.ErrorRate),
			strconv.Itoa(ip.UniqueURLs),
		})
	}

	// Top URLs
	result = append(result, []string{""})
	result = append(result, []string{"===== Top 10 URL 路徑 ====="})
	result = append(result, []string{"URL", "請求次數", "總傳輸量", "錯誤次數", "錯誤率(%)", "平均大小"})

	limit = len(stats.TopURLs)
	if limit > 10 {
		limit = 10
	}

	for i := 0; i < limit; i++ {
		url := stats.TopURLs[i]
		result = append(result, []string{
			url.URL,
			strconv.FormatInt(url.Count, 10),
			strconv.FormatInt(url.TotalBytes, 10),
			strconv.FormatInt(url.ErrorCount, 10),
			fmt.Sprintf("%.2f", url.ErrorRate),
			fmt.Sprintf("%.2f", url.AvgBytes),
		})
	}

	// HTTP 方法分布
	if len(stats.MethodDist) > 0 {
		result = append(result, []string{""})
		result = append(result, []string{"===== HTTP 方法分布 ====="})
		result = append(result, []string{"方法", "次數"})

		// 排序方法
		var methods []string
		for method := range stats.MethodDist {
			methods = append(methods, method)
		}
		sort.Strings(methods)

		for _, method := range methods {
			count := stats.MethodDist[method]
			result = append(result, []string{method, strconv.FormatInt(count, 10)})
		}
	}

	return result
}

// FormatBotDetection 格式化機器人偵測結果為二維字串陣列
// 分析 User Agent 並識別機器人請求
func (f *Formatter) FormatBotDetection(logs []*models.LogEntry) [][]string {
	// 建立標題行
	headers := []string{"IP位址", "機器人類型", "信心分數", "請求次數"}
	result := [][]string{headers}

	if logs == nil || len(logs) == 0 {
		return result
	}

	// 統計每個IP的機器人活動
	botIPs := make(map[string]*BotIPStat)

	for _, log := range logs {
		if log == nil {
			continue
		}

		// 使用內建的機器人偵測邏輯
		isBot, botType := f.detectBot(log.UserAgent)

		if isBot {
			if stat, exists := botIPs[log.IP]; exists {
				stat.Count++
				// 如果發現更具體的機器人類型，更新類型
				if f.getBotTypeSpecificity(botType) > f.getBotTypeSpecificity(stat.BotType) {
					stat.BotType = botType
				}
			} else {
				botIPs[log.IP] = &BotIPStat{
					IP:      log.IP,
					BotType: botType,
					Count:   1,
				}
			}
		}
	}

	// 轉換為切片並排序（按請求次數降序）
	var botList []*BotIPStat
	for _, stat := range botIPs {
		botList = append(botList, stat)
	}

	sort.Slice(botList, func(i, j int) bool {
		return botList[i].Count > botList[j].Count
	})

	// 格式化為表格行
	for _, bot := range botList {
		confidence := f.calculateConfidence(bot.BotType, bot.Count)
		row := []string{
			bot.IP,
			bot.BotType,
			confidence,
			strconv.Itoa(bot.Count),
		}
		result = append(result, row)
	}

	return result
}

// BotIPStat 機器人IP統計資料
type BotIPStat struct {
	IP      string
	BotType string
	Count   int
}

// formatTime 格式化時間為字串
func (f *Formatter) formatTime(t time.Time) string {
	if t.IsZero() {
		return "0001-01-01 00:00:00"
	}
	return t.Format(f.timeFormat)
}

// formatReferer 格式化 Referer 欄位
func (f *Formatter) formatReferer(referer string) string {
	if referer == "" || referer == "-" {
		return "-"
	}
	return referer
}

// getBotTypeSpecificity 取得機器人類型的特異性評分
// 特異性越高的類型越具體，優先度越高
func (f *Formatter) getBotTypeSpecificity(botType string) int {
	specificity := map[string]int{
		"搜尋引擎":   5,
		"社交媒體":   4,
		"SEO 工具": 4,
		"監控工具":   3,
		"安全掃描":   3,
		"爬蟲":     1, // 最通用的類別
	}

	if score, exists := specificity[botType]; exists {
		return score
	}
	return 0
}

// detectBot 內建機器人偵測功能（簡化版）
func (f *Formatter) detectBot(userAgent string) (bool, string) {
	if userAgent == "" || userAgent == "-" {
		return false, ""
	}

	lowerUA := strings.ToLower(userAgent)

	// 搜尋引擎機器人
	searchEnginePatterns := []string{
		"googlebot", "bingbot", "slurp", "duckduckbot",
		"baiduspider", "yandexbot", "sogou", "exabot",
	}
	for _, pattern := range searchEnginePatterns {
		if strings.Contains(lowerUA, pattern) {
			return true, "搜尋引擎"
		}
	}

	// 社交媒體機器人
	socialMediaPatterns := []string{
		"facebookexternalhit", "twitterbot", "linkedinbot",
		"pinterest", "slackbot", "telegrambot", "whatsapp",
		"discordbot",
	}
	for _, pattern := range socialMediaPatterns {
		if strings.Contains(lowerUA, pattern) {
			return true, "社交媒體"
		}
	}

	// 監控工具
	monitoringPatterns := []string{
		"pingdom", "uptimerobot", "statuscake", "monitor",
		"site24x7", "newrelic", "datadog", "nagios",
	}
	for _, pattern := range monitoringPatterns {
		if strings.Contains(lowerUA, pattern) {
			return true, "監控工具"
		}
	}

	// 通用爬蟲關鍵字
	crawlerPatterns := []string{
		"bot", "crawler", "spider", "scraper", "scraping",
		"python-requests", "curl", "wget", "httpclient",
		"scrapy", "beautifulsoup", "mechanize", "pycurl",
	}
	for _, pattern := range crawlerPatterns {
		if strings.Contains(lowerUA, pattern) {
			return true, "爬蟲"
		}
	}

	return false, ""
}
func (f *Formatter) calculateConfidence(botType string, requestCount int) string {
	// 基於機器人類型的基礎信心分數
	baseScore := f.getBotTypeSpecificity(botType)

	// 基於請求次數調整信心分數
	// 機器人通常會有較高的請求頻率
	var frequencyScore int
	if requestCount >= 100 {
		frequencyScore = 3
	} else if requestCount >= 50 {
		frequencyScore = 2
	} else if requestCount >= 10 {
		frequencyScore = 1
	} else {
		frequencyScore = 0
	}

	totalScore := baseScore + frequencyScore

	// 轉換為信心等級
	if totalScore >= 7 {
		return "極高"
	} else if totalScore >= 5 {
		return "高"
	} else if totalScore >= 3 {
		return "中"
	} else {
		return "低"
	}
}

// SetTimeFormat 設定時間格式化字串
func (f *Formatter) SetTimeFormat(format string) {
	f.timeFormat = format
}

// GetTimeFormat 取得當前時間格式化字串
func (f *Formatter) GetTimeFormat() string {
	return f.timeFormat
}

// FormatFileSize 格式化檔案大小為人類可讀的格式
func (f *Formatter) FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// EscapeCSVField 轉義 CSV 欄位中的特殊字元
// 處理包含逗號、引號或換行的欄位
func (f *Formatter) EscapeCSVField(field string) string {
	// 如果欄位包含逗號、引號或換行，需要用引號包圍
	if strings.Contains(field, ",") || strings.Contains(field, "\"") || strings.Contains(field, "\n") || strings.Contains(field, "\r") {
		// 轉義內部的引號（雙引號變成四引號）
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		return fmt.Sprintf("\"%s\"", escaped)
	}
	return field
}

// ValidateData 驗證資料完整性
// 檢查是否有必要的欄位缺失或格式錯誤
func (f *Formatter) ValidateData(logs []*models.LogEntry, stats *models.Statistics) []string {
	var warnings []string

	// 檢查日誌條目
	if logs == nil {
		warnings = append(warnings, "日誌條目為空")
	} else {
		for i, log := range logs {
			if log == nil {
				warnings = append(warnings, fmt.Sprintf("第 %d 筆日誌條目為 nil", i+1))
				continue
			}

			// 檢查必要欄位
			if log.IP == "" {
				warnings = append(warnings, fmt.Sprintf("第 %d 筆記錄缺少IP位址", i+1))
			}
			if log.Timestamp.IsZero() {
				warnings = append(warnings, fmt.Sprintf("第 %d 筆記錄缺少時間戳", i+1))
			}
			if log.StatusCode == 0 {
				warnings = append(warnings, fmt.Sprintf("第 %d 筆記錄缺少狀態碼", i+1))
			}
		}
	}

	// 檢查統計資料
	if stats == nil {
		warnings = append(warnings, "統計資料為空")
	} else {
		if stats.TotalRequests == 0 {
			warnings = append(warnings, "統計資料中總請求數為0")
		}
		if len(logs) > 0 && int64(len(logs)) != stats.TotalRequests {
			warnings = append(warnings, fmt.Sprintf("日誌條目數量 (%d) 與統計中的總請求數 (%d) 不符", len(logs), stats.TotalRequests))
		}
	}

	return warnings
}
