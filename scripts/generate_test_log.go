package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// 測試日誌產生器
// 用於產生各種大小的 Apache access log 測試檔案

var (
	// 命令列參數
	lines       = flag.Int("lines", 10000, "要產生的日誌行數")
	outputPath  = flag.String("output", "testdata/test.log", "輸出檔案路徑")
	errorRate   = flag.Float64("error-rate", 0.05, "錯誤行比率 (0.0-1.0)")
	invalidRate = flag.Float64("invalid-rate", 0.01, "無效行比率 (0.0-1.0)")
)

// 常見的請求路徑
var paths = []string{
	"/",
	"/index.html",
	"/about.html",
	"/contact.html",
	"/products",
	"/products/item1",
	"/products/item2",
	"/api/users",
	"/api/orders",
	"/static/css/style.css",
	"/static/js/app.js",
	"/images/logo.png",
	"/favicon.ico",
	"/robots.txt",
	"/sitemap.xml",
}

// 常見的 User-Agent
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
}

// HTTP 方法
var methods = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}

// 正常狀態碼
var normalStatusCodes = []int{200, 200, 200, 200, 201, 204, 301, 302, 304}

// 錯誤狀態碼
var errorStatusCodes = []int{400, 401, 403, 404, 500, 502, 503}

// 常見的 Referer
var referers = []string{
	"https://www.google.com/",
	"https://www.bing.com/",
	"https://www.facebook.com/",
	"https://twitter.com/",
	"https://www.linkedin.com/",
	"-",
}

func main() {
	flag.Parse()

	// 確保輸出目錄存在
	dir := filepath.Dir(*outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "無法建立目錄 %s: %v\n", dir, err)
		os.Exit(1)
	}

	// 開啟輸出檔案
	file, err := os.Create(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "無法建立檔案 %s: %v\n", *outputPath, err)
		os.Exit(1)
	}
	defer file.Close()

	// 初始化隨機數生成器
	rand.Seed(time.Now().UnixNano())

	// 產生日誌行
	startTime := time.Now().Add(-24 * time.Hour) // 從 24 小時前開始

	for i := 0; i < *lines; i++ {
		// 決定此行是否為無效行
		if rand.Float64() < *invalidRate {
			// 產生無效行
			line := generateInvalidLine()
			fmt.Fprintln(file, line)
			continue
		}

		// 產生正常的日誌行
		timestamp := startTime.Add(time.Duration(i) * time.Second)

		// 決定是否為錯誤請求
		isError := rand.Float64() < *errorRate

		line := generateLogLine(timestamp, isError)
		fmt.Fprintln(file, line)
	}

	fmt.Printf("成功產生 %d 行日誌到 %s\n", *lines, *outputPath)

	// 顯示檔案大小
	if info, err := file.Stat(); err == nil {
		sizeMB := float64(info.Size()) / (1024 * 1024)
		fmt.Printf("檔案大小: %.2f MB\n", sizeMB)
	}
}

// generateLogLine 產生一行標準的 Apache Combined Log Format 日誌
func generateLogLine(timestamp time.Time, isError bool) string {
	// 隨機 IP 位址
	ip := fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256))

	// 隨機選擇路徑和方法
	path := paths[rand.Intn(len(paths))]
	method := methods[rand.Intn(len(methods))]

	// 根據是否錯誤選擇狀態碼
	var statusCode int
	if isError {
		statusCode = errorStatusCodes[rand.Intn(len(errorStatusCodes))]
	} else {
		statusCode = normalStatusCodes[rand.Intn(len(normalStatusCodes))]
	}

	// 隨機回應大小（100 bytes 到 100 KB）
	responseSize := rand.Intn(102400) + 100

	// 隨機 User-Agent 和 Referer
	userAgent := userAgents[rand.Intn(len(userAgents))]
	referer := referers[rand.Intn(len(referers))]

	// 格式化時間戳
	timeStr := timestamp.Format("02/Jan/2006:15:04:05 -0700")

	// Combined Log Format:
	// IP - - [timestamp] "method path protocol" status size "referer" "user-agent"
	return fmt.Sprintf(`%s - - [%s] "%s %s HTTP/1.1" %d %d "%s" "%s"`,
		ip, timeStr, method, path, statusCode, responseSize, referer, userAgent)
}

// generateInvalidLine 產生無效的日誌行（用於測試錯誤處理）
func generateInvalidLine() string {
	invalidLines := []string{
		"This is not a valid log line",
		"",
		"123.456.789.0",
		"[01/Jan/2024:00:00:00 +0000]",
		"GET /path HTTP/1.1",
		"Invalid timestamp format",
		"Missing required fields",
	}

	return invalidLines[rand.Intn(len(invalidLines))]
}
