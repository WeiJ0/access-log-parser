package models

// Statistics 表示日誌的統計資料
// 包含所有用戶故事需要的統計指標
type Statistics struct {
	// 基本統計
	TotalRequests int64 `json:"totalRequests"` // 總請求數
	UniqueIPs     int   `json:"uniqueIPs"`     // 唯一 IP 數量
	TotalBytes    int64 `json:"totalBytes"`    // 總傳輸量（位元組）

	// 時間範圍
	StartTime string `json:"startTime"` // 最早請求時間
	EndTime   string `json:"endTime"`   // 最晚請求時間

	// 狀態碼分布（用戶故事 #2）
	StatusCodeDist map[int]int64    `json:"statusCodeDist"` // 狀態碼 -> 次數
	StatusCatDist  map[string]int64 `json:"statusCatDist"`  // 狀態類別 -> 次數（2xx, 3xx, 4xx, 5xx）

	// URL 統計（用戶故事 #2）
	TopURLs []URLStat `json:"topURLs"` // 最熱門 URL（Top 100）

	// 時間分布（用戶故事 #2）
	HourlyDist map[int]int64    `json:"hourlyDist"` // 小時 -> 請求數（0-23）
	DailyDist  map[string]int64 `json:"dailyDist"`  // 日期 -> 請求數（YYYY-MM-DD）

	// IP 統計（用戶故事 #2）
	TopIPs []IPStat `json:"topIPs"` // 最活躍 IP（Top 100）

	// HTTP 方法分布（用戶故事 #2）
	MethodDist map[string]int64 `json:"methodDist"` // 方法 -> 次數（GET, POST, 等）

	// 回應大小統計（用戶故事 #2）
	AvgResponseSize float64 `json:"avgResponseSize"` // 平均回應大小（位元組）
	MinResponseSize int64   `json:"minResponseSize"` // 最小回應大小
	MaxResponseSize int64   `json:"maxResponseSize"` // 最大回應大小

	// 錯誤統計（用戶故事 #3）
	ErrorCount       int64     `json:"errorCount"`       // 錯誤請求數（4xx + 5xx）
	ClientErrorCount int64     `json:"clientErrorCount"` // 客戶端錯誤數（4xx）
	ServerErrorCount int64     `json:"serverErrorCount"` // 伺服器錯誤數（5xx）
	ErrorRate        float64   `json:"errorRate"`        // 錯誤率（百分比）
	TopErrorURLs     []URLStat `json:"topErrorURLs"`     // 最常出錯 URL（Top 50）
	TopErrorIPs      []IPStat  `json:"topErrorIPs"`      // 最常出錯 IP（Top 50）

	// User Agent 統計（用戶故事 #2）
	TopUserAgents []UserAgentStat  `json:"topUserAgents"` // 最常見 User Agent（Top 50）
	BrowserDist   map[string]int64 `json:"browserDist"`   // 瀏覽器分布（簡化版）
	OSdist        map[string]int64 `json:"osDist"`        // 作業系統分布（簡化版）

	// Referer 統計（用戶故事 #2）
	TopReferers []RefererStat `json:"topReferers"` // 最常見來源（Top 50）
}

// URLStat 表示 URL 的統計資料
// 用於 Top URLs 排行
type URLStat struct {
	URL        string  `json:"url"`        // URL 路徑
	Count      int64   `json:"count"`      // 訪問次數
	TotalBytes int64   `json:"totalBytes"` // 總傳輸量（位元組）
	ErrorCount int64   `json:"errorCount"` // 錯誤次數（4xx + 5xx）
	ErrorRate  float64 `json:"errorRate"`  // 錯誤率（百分比）
	AvgBytes   float64 `json:"avgBytes"`   // 平均回應大小
}

// IPStat 表示 IP 的統計資料
// 用於 Top IPs 排行
type IPStat struct {
	IP         string  `json:"ip"`         // IP 位址
	Count      int64   `json:"count"`      // 請求次數
	TotalBytes int64   `json:"totalBytes"` // 總傳輸量（位元組）
	ErrorCount int64   `json:"errorCount"` // 錯誤次數（4xx + 5xx）
	ErrorRate  float64 `json:"errorRate"`  // 錯誤率（百分比）
	UniqueURLs int     `json:"uniqueUrls"` // 訪問的唯一 URL 數量
}

// UserAgentStat 表示 User Agent 的統計資料
// 用於 User Agent 分析
type UserAgentStat struct {
	UserAgent  string  `json:"userAgent"`  // User Agent 字串
	Count      int64   `json:"count"`      // 出現次數
	Percentage float64 `json:"percentage"` // 佔比（百分比）
}

// RefererStat 表示 Referer 的統計資料
// 用於來源分析
type RefererStat struct {
	Referer    string  `json:"referer"`    // Referer URL
	Count      int64   `json:"count"`      // 出現次數
	Percentage float64 `json:"percentage"` // 佔比（百分比）
}

// NewStatistics 建立新的 Statistics 實例
// 初始化所有映射和切片
func NewStatistics() *Statistics {
	return &Statistics{
		StatusCodeDist: make(map[int]int64),
		StatusCatDist:  make(map[string]int64),
		TopURLs:        make([]URLStat, 0),
		HourlyDist:     make(map[int]int64),
		DailyDist:      make(map[string]int64),
		TopIPs:         make([]IPStat, 0),
		MethodDist:     make(map[string]int64),
		TopErrorURLs:   make([]URLStat, 0),
		TopErrorIPs:    make([]IPStat, 0),
		TopUserAgents:  make([]UserAgentStat, 0),
		BrowserDist:    make(map[string]int64),
		OSdist:         make(map[string]int64),
		TopReferers:    make([]RefererStat, 0),
	}
}

// IsEmpty 檢查統計資料是否為空
// 判斷是否有任何請求資料
func (s *Statistics) IsEmpty() bool {
	return s.TotalRequests == 0
}

// GetTotalBytesMB 取得總傳輸量（MB）
// 返回浮點數，精確到小數點後兩位
func (s *Statistics) GetTotalBytesMB() float64 {
	return float64(s.TotalBytes) / (1024 * 1024)
}

// GetTotalBytesGB 取得總傳輸量（GB）
// 返回浮點數，精確到小數點後兩位
func (s *Statistics) GetTotalBytesGB() float64 {
	return float64(s.TotalBytes) / (1024 * 1024 * 1024)
}

// Validate 驗證 Statistics 的資料一致性
// 返回 ValidationError 如果驗證失敗
func (s *Statistics) Validate() error {
	// 驗證總請求數不為負數
	if s.TotalRequests < 0 {
		return &ValidationError{
			Field:   "TotalRequests",
			Value:   string(rune(s.TotalRequests)),
			Message: "總請求數不能為負數",
		}
	}

	// 驗證唯一 IP 數量不為負數
	if s.UniqueIPs < 0 {
		return &ValidationError{
			Field:   "UniqueIPs",
			Value:   string(rune(s.UniqueIPs)),
			Message: "唯一 IP 數量不能為負數",
		}
	}

	// 驗證總傳輸量不為負數
	if s.TotalBytes < 0 {
		return &ValidationError{
			Field:   "TotalBytes",
			Value:   string(rune(s.TotalBytes)),
			Message: "總傳輸量不能為負數",
		}
	}

	// 驗證錯誤率在有效範圍內（0-100）
	if s.ErrorRate < 0 || s.ErrorRate > 100 {
		return &ValidationError{
			Field:   "ErrorRate",
			Value:   string(rune(int(s.ErrorRate))),
			Message: "錯誤率必須在 0-100 範圍內",
		}
	}

	// 驗證錯誤數量一致性
	if s.ErrorCount != s.ClientErrorCount+s.ServerErrorCount {
		return &ValidationError{
			Field:   "ErrorCount",
			Value:   "",
			Message: "錯誤總數應等於客戶端錯誤數加伺服器錯誤數",
		}
	}

	return nil
}
