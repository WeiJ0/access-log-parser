package parser

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
	
	"access-log-analyzer/internal/models"
	"access-log-analyzer/pkg/apachelog"
	"access-log-analyzer/pkg/logger"
)

// Parser 提供 Apache log 解析功能
// 使用 worker pool 模式實現並行解析
type Parser struct {
	format      LogFormat
	workerCount int
	maxErrors   int
	log         *logger.Logger
}

// ParseResult 包含解析結果和相關統計資訊
type ParseResult struct {
	Entries       []models.LogEntry  // 成功解析的日誌記錄
	TotalLines    int                // 總行數
	ParsedLines   int                // 成功解析的行數
	ErrorLines    int                // 失敗的行數
	ErrorSamples  []ParseError       // 錯誤樣本（最多 100 筆）
	ParseTime     time.Duration      // 解析耗時
	MemoryUsed    int64              // 記憶體使用量（位元組）
	ThroughputMB  float64            // 吞吐量（MB/秒）
}

// ParseError 記錄解析錯誤的詳細資訊
type ParseError struct {
	LineNumber int    `json:"lineNumber"` // 錯誤行號
	Line       string `json:"line"`       // 原始行內容
	Error      string `json:"error"`      // 錯誤訊息
}

// NewParser 建立新的解析器實例
// workerCount: worker 數量（0 表示使用 CPU 核心數）
func NewParser(format LogFormat, workerCount int) *Parser {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}
	
	return &Parser{
		format:      format,
		workerCount: workerCount,
		maxErrors:   100, // 最多收集 100 個錯誤樣本
		log:         logger.Get().WithModule("parser"),
	}
}

// ParseFile 解析指定的 log 檔案
// 返回解析結果或錯誤
func (p *Parser) ParseFile(filepath string, fileSize int64) (*ParseResult, error) {
	startTime := time.Now()
	
	p.log.Info().
		Str("file", filepath).
		Int64("size", fileSize).
		Int("workers", p.workerCount).
		Msg("開始解析 log 檔案")
	
	// 記錄起始記憶體使用量
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)
	
	// 開啟檔案讀取器
	reader, err := apachelog.NewReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("無法開啟檔案: %w", err)
	}
	defer reader.Close()
	
	// 建立 channels 用於 worker 通訊
	lineChan := make(chan lineData, p.workerCount*2)
	resultChan := make(chan parseResult, p.workerCount*2)
	
	// 啟動 worker pool
	var wg sync.WaitGroup
	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go p.worker(lineChan, resultChan, &wg)
	}
	
	// 啟動結果收集器
	done := make(chan *ParseResult)
	go p.collectResults(resultChan, done)
	
	// 讀取並分發行資料
	for {
		lineNum, line, hasMore := reader.ReadLine()
		if !hasMore {
			break
		}
		
		lineChan <- lineData{
			lineNum: lineNum,
			line:    line,
		}
	}
	
	// 關閉 channel 並等待完成
	close(lineChan)
	wg.Wait()
	close(resultChan)
	
	// 等待結果收集完成
	result := <-done
	
	// 檢查讀取器錯誤
	if err := reader.Error(); err != nil {
		p.log.Warn().Err(err).Msg("讀取檔案時發生錯誤")
	}
	
	// 計算效能指標
	result.ParseTime = time.Since(startTime)
	
	var memEnd runtime.MemStats
	runtime.ReadMemStats(&memEnd)
	result.MemoryUsed = int64(memEnd.Alloc - memStart.Alloc)
	
	if result.ParseTime.Seconds() > 0 {
		sizeMB := float64(fileSize) / (1024 * 1024)
		result.ThroughputMB = sizeMB / result.ParseTime.Seconds()
	}
	
	p.log.Info().
		Int("total", result.TotalLines).
		Int("parsed", result.ParsedLines).
		Int("errors", result.ErrorLines).
		Dur("time", result.ParseTime).
		Float64("throughput_mb_s", result.ThroughputMB).
		Msg("解析完成")
	
	return result, nil
}

// lineData 封裝傳遞給 worker 的行資料
type lineData struct {
	lineNum int
	line    string
}

// parseResult 封裝 worker 的解析結果
type parseResult struct {
	entry  *models.LogEntry
	err    *ParseError
}

// worker 處理行資料並解析為 LogEntry
func (p *Parser) worker(lines <-chan lineData, results chan<- parseResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	pattern := GetPattern(p.format)
	
	for data := range lines {
		entry, err := p.parseLine(data.lineNum, data.line, pattern)
		if err != nil {
			results <- parseResult{
				err: &ParseError{
					LineNumber: data.lineNum,
					Line:       data.line,
					Error:      err.Error(),
				},
			}
		} else {
			results <- parseResult{entry: entry}
		}
	}
}

// parseLine 解析單行 log 資料
func (p *Parser) parseLine(lineNum int, line string, pattern *regexp.Regexp) (*models.LogEntry, error) {
	matches := pattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("無法匹配 log 格式")
	}
	
	// 根據格式類型解析欄位
	var entry models.LogEntry
	entry.LineNumber = lineNum
	entry.RawLine = line
	
	if p.format == FormatCombined {
		// Combined 格式: IP, ident, user, time, method, url, protocol, status, size, referer, ua
		entry.IP = matches[1]
		// matches[2] 是 ident（通常是 -）
		entry.User = matches[3]
		
		// 解析時間戳
		timestamp, err := parseApacheTime(matches[4])
		if err != nil {
			return nil, fmt.Errorf("無法解析時間戳: %w", err)
		}
		entry.Timestamp = timestamp
		
		entry.Method = matches[5]
		entry.URL = matches[6]
		entry.Protocol = matches[7]
		
		// 解析狀態碼
		statusCode, err := strconv.Atoi(matches[8])
		if err != nil {
			return nil, fmt.Errorf("無法解析狀態碼: %w", err)
		}
		entry.StatusCode = statusCode
		
		// 解析回應大小
		if matches[9] != "-" {
			size, err := strconv.ParseInt(matches[9], 10, 64)
			if err == nil {
				entry.ResponseBytes = size
			}
		}
		
		entry.Referer = matches[10]
		entry.UserAgent = matches[11]
		
	} else if p.format == FormatCommon {
		// Common 格式: IP, ident, user, time, method, url, protocol, status, size
		entry.IP = matches[1]
		entry.User = matches[3]
		
		timestamp, err := parseApacheTime(matches[4])
		if err != nil {
			return nil, fmt.Errorf("無法解析時間戳: %w", err)
		}
		entry.Timestamp = timestamp
		
		entry.Method = matches[5]
		entry.URL = matches[6]
		entry.Protocol = matches[7]
		
		statusCode, err := strconv.Atoi(matches[8])
		if err != nil {
			return nil, fmt.Errorf("無法解析狀態碼: %w", err)
		}
		entry.StatusCode = statusCode
		
		if matches[9] != "-" {
			size, err := strconv.ParseInt(matches[9], 10, 64)
			if err == nil {
				entry.ResponseBytes = size
			}
		}
	}
	
	return &entry, nil
}

// parseApacheTime 解析 Apache log 時間格式
// 格式: 06/Nov/2025:14:30:15 +0800
func parseApacheTime(timeStr string) (time.Time, error) {
	// Apache 時間格式
	layout := "02/Jan/2006:15:04:05 -0700"
	return time.Parse(layout, timeStr)
}

// collectResults 收集所有 worker 的解析結果
func (p *Parser) collectResults(results <-chan parseResult, done chan<- *ParseResult) {
	result := &ParseResult{
		Entries:      make([]models.LogEntry, 0),
		ErrorSamples: make([]ParseError, 0),
	}
	
	for res := range results {
		result.TotalLines++
		
		if res.err != nil {
			result.ErrorLines++
			// 只保留前 maxErrors 個錯誤樣本
			if len(result.ErrorSamples) < p.maxErrors {
				result.ErrorSamples = append(result.ErrorSamples, *res.err)
			}
		} else {
			result.ParsedLines++
			result.Entries = append(result.Entries, *res.entry)
		}
	}
	
	done <- result
}

// ValidateFormat 快速驗證檔案格式
// 讀取前 N 行並嘗試解析以驗證格式正確性
func (p *Parser) ValidateFormat(filepath string, sampleLines int) (bool, error) {
	reader, err := apachelog.NewReader(filepath)
	if err != nil {
		return false, err
	}
	defer reader.Close()
	
	pattern := GetPattern(p.format)
	validCount := 0
	totalCount := 0
	
	for i := 0; i < sampleLines; i++ {
		_, line, hasMore := reader.ReadLine()
		if !hasMore {
			break
		}
		
		totalCount++
		if pattern.MatchString(line) {
			validCount++
		}
	}
	
	// 如果超過 80% 的行匹配，則認為格式正確
	if totalCount == 0 {
		return false, fmt.Errorf("檔案為空")
	}
	
	successRate := float64(validCount) / float64(totalCount)
	return successRate >= 0.8, nil
}

// ValidateFirstLine 快速驗證檔案第一行是否為 Apache log 格式
// 用於在選擇檔案後立即檢查格式，提供快速回饋
func (p *Parser) ValidateFirstLine(filepath string) error {
	reader, err := apachelog.NewReader(filepath)
	if err != nil {
		return fmt.Errorf("無法開啟檔案: %w", err)
	}
	defer reader.Close()
	
	// 讀取第一行
	_, line, hasMore := reader.ReadLine()
	if !hasMore {
		return fmt.Errorf("檔案為空")
	}
	
	// 檢查是否符合 Apache log 格式
	pattern := GetPattern(p.format)
	if !pattern.MatchString(line) {
		return fmt.Errorf("第一行不符合 Apache Access Log 格式。請確認檔案是 Apache Combined 或 Common 格式的 access log")
	}
	
	return nil
}
