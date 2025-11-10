package models

import (
	"time"
)

// LogEntry 表示一筆 Apache access log 記錄
// 對應標準 Combined Log Format 格式
type LogEntry struct {
	// 基本欄位
	IP            string    `json:"ip"`            // 客戶端 IP 位址
	Timestamp     time.Time `json:"timestamp"`     // 請求時間
	Method        string    `json:"method"`        // HTTP 方法（GET, POST, 等）
	URL           string    `json:"url"`           // 請求的 URL 路徑
	Protocol      string    `json:"protocol"`      // HTTP 協定版本
	StatusCode    int       `json:"statusCode"`    // HTTP 狀態碼
	ResponseBytes int64     `json:"responseBytes"` // 回應大小（位元組）
	Referer       string    `json:"referer"`       // 來源頁面
	UserAgent     string    `json:"userAgent"`     // 使用者代理字串

	// 擴展欄位
	User        string `json:"user,omitempty"`        // 認證使用者名稱（如果有）
	RequestTime int64  `json:"requestTime,omitempty"` // 請求處理時間（微秒）

	// 內部欄位
	LineNumber int    `json:"lineNumber"`           // 原始檔案中的行號
	RawLine    string `json:"rawLine"`              // 原始日誌行（用於錯誤分析）
	ParseError string `json:"parseError,omitempty"` // 解析錯誤訊息（如果有）
}

// IsError 檢查此記錄是否為錯誤狀態
// HTTP 狀態碼 >= 400 視為錯誤
func (e *LogEntry) IsError() bool {
	return e.StatusCode >= 400
}

// IsClientError 檢查是否為客戶端錯誤（4xx）
// HTTP 狀態碼 400-499
func (e *LogEntry) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError 檢查是否為伺服器錯誤（5xx）
// HTTP 狀態碼 500-599
func (e *LogEntry) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// IsSuccess 檢查請求是否成功（2xx 或 3xx）
// HTTP 狀態碼 200-399
func (e *LogEntry) IsSuccess() bool {
	return e.StatusCode >= 200 && e.StatusCode < 400
}

// GetStatusCategory 取得狀態碼類別字串
// 返回 "2xx", "3xx", "4xx", "5xx" 或 "其他"
func (e *LogEntry) GetStatusCategory() string {
	switch {
	case e.StatusCode >= 200 && e.StatusCode < 300:
		return "2xx"
	case e.StatusCode >= 300 && e.StatusCode < 400:
		return "3xx"
	case e.StatusCode >= 400 && e.StatusCode < 500:
		return "4xx"
	case e.StatusCode >= 500 && e.StatusCode < 600:
		return "5xx"
	default:
		return "其他"
	}
}

// ResponseSizeKB 取得回應大小（KB）
// 返回浮點數，精確到小數點後兩位
func (e *LogEntry) ResponseSizeKB() float64 {
	return float64(e.ResponseBytes) / 1024.0
}

// ResponseSizeMB 取得回應大小（MB）
// 返回浮點數，精確到小數點後兩位
func (e *LogEntry) ResponseSizeMB() float64 {
	return float64(e.ResponseBytes) / (1024.0 * 1024.0)
}

// Validate 驗證 LogEntry 的必填欄位
// 返回 ValidationError 如果驗證失敗
func (e *LogEntry) Validate() error {
	// 驗證 IP 位址不為空
	if e.IP == "" {
		return &ValidationError{
			Field:   "IP",
			Value:   "",
			Message: "IP 位址不能為空",
		}
	}

	// 驗證時間戳不為零值
	if e.Timestamp.IsZero() {
		return &ValidationError{
			Field:   "Timestamp",
			Value:   "",
			Message: "時間戳不能為空",
		}
	}

	// 驗證 HTTP 方法不為空
	if e.Method == "" {
		return &ValidationError{
			Field:   "Method",
			Value:   "",
			Message: "HTTP 方法不能為空",
		}
	}

	// 驗證 URL 不為空
	if e.URL == "" {
		return &ValidationError{
			Field:   "URL",
			Value:   "",
			Message: "URL 不能為空",
		}
	}

	// 驗證狀態碼在有效範圍內（100-599）
	if e.StatusCode < 100 || e.StatusCode > 599 {
		return &ValidationError{
			Field:   "StatusCode",
			Value:   string(rune(e.StatusCode)),
			Message: "HTTP 狀態碼必須在 100-599 範圍內",
		}
	}

	// 驗證回應大小不為負數
	if e.ResponseBytes < 0 {
		return &ValidationError{
			Field:   "ResponseBytes",
			Value:   string(rune(e.ResponseBytes)),
			Message: "回應大小不能為負數",
		}
	}

	return nil
}
