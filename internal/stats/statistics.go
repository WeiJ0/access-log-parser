package stats

import (
	"access-log-analyzer/internal/models"
	"access-log-analyzer/pkg/logger"
)

// Calculator 提供日誌統計計算功能
// 整合 Top-N 演算法和機器人偵測功能
type Calculator struct {
	topN        int          // Top-N 的 N 值
	botDetector *BotDetector // 機器人偵測器
	log         *logger.Logger
}

// Statistics 包含完整的統計資訊
type Statistics struct {
	TotalRequests          int                  `json:"totalRequests"`          // 總請求數
	UniqueIPs              int                  `json:"uniqueIPs"`              // 唯一 IP 數量
	UniquePaths            int                  `json:"uniquePaths"`            // 唯一路徑數量
	TotalBytes             int64                `json:"totalBytes"`             // 總傳輸量（位元組）
	AverageResponseSize    int64                `json:"averageResponseSize"`    // 平均回應大小
	TopIPs                 []IPStatistics       `json:"topIPs"`                 // Top IP 統計
	TopPaths               []PathStatistics     `json:"topPaths"`               // Top 路徑統計
	StatusCodeDistribution StatusCodeStatistics `json:"statusCodeDistribution"` // 狀態碼分布
	BotStats               BotStats             `json:"botStats"`               // 機器人統計
}

// IPStatistics IP 統計資訊
type IPStatistics struct {
	IP           string `json:"ip"`           // IP 位址
	RequestCount int    `json:"requestCount"` // 請求次數
	TotalBytes   int64  `json:"totalBytes"`   // 總傳輸量（位元組）
}

// PathStatistics 路徑統計資訊
type PathStatistics struct {
	Path         string  `json:"path"`         // 路徑
	RequestCount int     `json:"requestCount"` // 請求次數
	AverageSize  int64   `json:"averageSize"`  // 平均大小
	ErrorRate    float64 `json:"errorRate"`    // 錯誤率（百分比）
}

// StatusCodeStatistics 狀態碼統計資訊
type StatusCodeStatistics struct {
	Success     int         `json:"success"`     // 2xx 成功
	Redirection int         `json:"redirection"` // 3xx 重定向
	ClientError int         `json:"clientError"` // 4xx 客戶端錯誤
	ServerError int         `json:"serverError"` // 5xx 伺服器錯誤
	Details     map[int]int `json:"details"`     // 詳細狀態碼分布
}

// NewCalculator 建立新的統計計算器
func NewCalculator() *Calculator {
	return &Calculator{
		topN:        10, // 預設保留 Top 10
		botDetector: NewBotDetector(),
		log:         logger.Get().WithModule("stats"),
	}
}

// SetTopN 設定 Top-N 的 N 值
func (c *Calculator) SetTopN(n int) {
	c.topN = n
}

// Calculate 計算日誌條目的統計資訊
// 使用單次遍歷和 Top-N 堆積實現高效計算
func (c *Calculator) Calculate(entries []models.LogEntry) Statistics {
	c.log.Info().Int("count", len(entries)).Msg("開始計算統計資訊")

	// 初始化統計結構
	stats := Statistics{
		TotalRequests: len(entries),
		StatusCodeDistribution: StatusCodeStatistics{
			Details: make(map[int]int),
		},
	}

	// 特殊情況：空資料集
	if len(entries) == 0 {
		return stats
	}

	// 初始化 Top-N 堆積
	ipHeap := NewTopNHeap(c.topN)
	pathHeap := NewTopNHeap(c.topN)

	// 用於追蹤唯一 IP 和路徑
	uniqueIPs := make(map[string]bool)
	uniquePaths := make(map[string]bool)

	// 用於計算路徑詳細統計
	pathStats := make(map[string]*pathStatAccumulator)

	// 用於計算 IP 詳細統計
	ipStats := make(map[string]*ipStatAccumulator)

	// 重置機器人偵測器統計
	c.botDetector.ResetStats()

	// 單次遍歷所有記錄
	for _, entry := range entries {
		// 統計唯一 IP 和路徑
		uniqueIPs[entry.IP] = true
		uniquePaths[entry.URL] = true

		// 累計總傳輸量
		stats.TotalBytes += entry.ResponseBytes

		// 統計 IP
		if _, exists := ipStats[entry.IP]; !exists {
			ipStats[entry.IP] = &ipStatAccumulator{}
		}
		ipStats[entry.IP].requestCount++
		ipStats[entry.IP].totalBytes += entry.ResponseBytes

		// 統計路徑
		if _, exists := pathStats[entry.URL]; !exists {
			pathStats[entry.URL] = &pathStatAccumulator{}
		}
		acc := pathStats[entry.URL]
		acc.requestCount++
		acc.totalBytes += entry.ResponseBytes
		if entry.StatusCode >= 400 {
			acc.errorCount++
		}

		// 統計狀態碼
		c.updateStatusCodeStats(&stats.StatusCodeDistribution, entry.StatusCode)

		// 機器人偵測
		c.botDetector.IsBot(entry.UserAgent)
	}

	// 設定唯一計數
	stats.UniqueIPs = len(uniqueIPs)
	stats.UniquePaths = len(uniquePaths)

	// 計算平均回應大小
	if stats.TotalRequests > 0 {
		stats.AverageResponseSize = stats.TotalBytes / int64(stats.TotalRequests)
	}

	// 建立 Top IP 統計
	for ip, acc := range ipStats {
		ipHeap.Push(ip, acc.requestCount)
	}

	topIPItems := ipHeap.GetResults()
	stats.TopIPs = make([]IPStatistics, len(topIPItems))
	for i, item := range topIPItems {
		acc := ipStats[item.Key]
		stats.TopIPs[i] = IPStatistics{
			IP:           item.Key,
			RequestCount: item.Count,
			TotalBytes:   acc.totalBytes,
		}
	}

	// 建立 Top 路徑統計
	for path, acc := range pathStats {
		pathHeap.Push(path, acc.requestCount)
	}

	topPathItems := pathHeap.GetResults()
	stats.TopPaths = make([]PathStatistics, len(topPathItems))
	for i, item := range topPathItems {
		acc := pathStats[item.Key]
		averageSize := int64(0)
		if acc.requestCount > 0 {
			averageSize = acc.totalBytes / int64(acc.requestCount)
		}
		errorRate := 0.0
		if acc.requestCount > 0 {
			errorRate = float64(acc.errorCount) / float64(acc.requestCount) * 100
		}

		stats.TopPaths[i] = PathStatistics{
			Path:         item.Key,
			RequestCount: item.Count,
			AverageSize:  averageSize,
			ErrorRate:    errorRate,
		}
	}

	// 獲取機器人統計
	stats.BotStats = c.botDetector.GetStats()

	c.log.Info().
		Int("totalRequests", stats.TotalRequests).
		Int("uniqueIPs", stats.UniqueIPs).
		Int("uniquePaths", stats.UniquePaths).
		Int("botRequests", stats.BotStats.BotRequests).
		Msg("統計計算完成")

	return stats
}

// updateStatusCodeStats 更新狀態碼統計
func (c *Calculator) updateStatusCodeStats(stats *StatusCodeStatistics, statusCode int) {
	// 增加詳細計數
	stats.Details[statusCode]++

	// 根據狀態碼範圍分類
	switch statusCode / 100 {
	case 2:
		stats.Success++
	case 3:
		stats.Redirection++
	case 4:
		stats.ClientError++
	case 5:
		stats.ServerError++
	}
}

// pathStatAccumulator 累積路徑統計資訊
type pathStatAccumulator struct {
	requestCount int
	totalBytes   int64
	errorCount   int
}

// ipStatAccumulator 累積 IP 統計資訊
type ipStatAccumulator struct {
	requestCount int
	totalBytes   int64
}
