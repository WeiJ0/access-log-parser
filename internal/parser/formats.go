package parser

import (
	"regexp"
)

// LogFormat 定義 Apache log 格式
type LogFormat int

const (
	// FormatCombined 對應 Apache Combined Log Format
	// 格式: %h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-agent}i\"
	FormatCombined LogFormat = iota
	
	// FormatCommon 對應 Apache Common Log Format
	// 格式: %h %l %u %t \"%r\" %>s %b
	FormatCommon
)

// 正規表達式模式定義
var (
	// combinedLogPattern 匹配 Combined Log Format
	// 範例: 192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"
	combinedLogPattern = regexp.MustCompile(
		`^(\S+) ` + // IP 位址
		`(\S+) ` + // 識別符號（通常是 -）
		`(\S+) ` + // 使用者 ID（通常是 -）
		`\[([^\]]+)\] ` + // 時間戳 [日期時間]
		`"([A-Z]+) ([^\s]+) ([^"]+)" ` + // 請求方法 URL 協定
		`(\d{3}) ` + // 狀態碼
		`(\S+) ` + // 回應大小（可能是 -）
		`"([^"]*)" ` + // Referer
		`"([^"]*)"`, // User-Agent
	)
	
	// commonLogPattern 匹配 Common Log Format
	// 範例: 192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234
	commonLogPattern = regexp.MustCompile(
		`^(\S+) ` + // IP 位址
		`(\S+) ` + // 識別符號（通常是 -）
		`(\S+) ` + // 使用者 ID（通常是 -）
		`\[([^\]]+)\] ` + // 時間戳 [日期時間]
		`"([A-Z]+) ([^\s]+) ([^"]+)" ` + // 請求方法 URL 協定
		`(\d{3}) ` + // 狀態碼
		`(\S+)`, // 回應大小（可能是 -）
	)
)

// GetPattern 根據格式類型返回對應的正規表達式
func GetPattern(format LogFormat) *regexp.Regexp {
	switch format {
	case FormatCombined:
		return combinedLogPattern
	case FormatCommon:
		return commonLogPattern
	default:
		return combinedLogPattern
	}
}

// DetectFormat 自動偵測 log 格式
// 透過分析 log 行的欄位數量來判斷
func DetectFormat(line string) LogFormat {
	// 嘗試匹配 Combined 格式（12 個群組）
	if combinedLogPattern.MatchString(line) {
		return FormatCombined
	}
	
	// 嘗試匹配 Common 格式（10 個群組）
	if commonLogPattern.MatchString(line) {
		return FormatCommon
	}
	
	// 預設使用 Combined 格式
	return FormatCombined
}
