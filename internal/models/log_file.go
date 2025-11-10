package models

import (
	"time"
)

// LogFile 表示一個已解析的日誌檔案
// 包含所有日誌記錄和相關元資料
type LogFile struct {
	// 檔案資訊
	Path     string    `json:"path"`     // 檔案完整路徑
	Name     string    `json:"name"`     // 檔案名稱
	Size     int64     `json:"size"`     // 檔案大小（位元組）
	LoadedAt time.Time `json:"loadedAt"` // 載入時間

	// 解析統計
	TotalLines  int `json:"totalLines"`  // 總行數
	ParsedLines int `json:"parsedLines"` // 成功解析的行數
	ErrorLines  int `json:"errorLines"`  // 解析失敗的行數

	// 日誌資料
	Entries []LogEntry `json:"entries"` // 所有日誌記錄

	// 統計資訊（User Story 2）
	Statistics interface{} `json:"statistics"` // 統計分析結果

	// 效能指標
	ParseTime  int64 `json:"parseTime"`  // 解析耗時（毫秒）
	StatTime   int64 `json:"statTime"`   // 統計計算耗時（毫秒）
	MemoryUsed int64 `json:"memoryUsed"` // 記憶體使用量（位元組）
}

// NewLogFile 建立新的 LogFile 實例
// 初始化基本欄位，entries 使用預分配容量
func NewLogFile(path, name string, size int64, estimatedLines int) *LogFile {
	return &LogFile{
		Path:       path,
		Name:       name,
		Size:       size,
		LoadedAt:   time.Now(),
		Entries:    make([]LogEntry, 0, estimatedLines),
		TotalLines: 0,
	}
}

// AddEntry 新增一筆日誌記錄
// 自動更新統計資訊
func (f *LogFile) AddEntry(entry LogEntry) {
	f.Entries = append(f.Entries, entry)
	f.ParsedLines++
}

// AddErrorEntry 新增一筆解析失敗的記錄
// 記錄錯誤訊息和原始行內容
func (f *LogFile) AddErrorEntry(lineNumber int, rawLine, errorMsg string) {
	f.Entries = append(f.Entries, LogEntry{
		LineNumber: lineNumber,
		RawLine:    rawLine,
		ParseError: errorMsg,
	})
	f.ErrorLines++
}

// GetSuccessRate 取得解析成功率（百分比）
// 返回 0-100 的浮點數
func (f *LogFile) GetSuccessRate() float64 {
	if f.TotalLines == 0 {
		return 0
	}
	return float64(f.ParsedLines) / float64(f.TotalLines) * 100
}

// GetErrorRate 取得錯誤率（百分比）
// 返回 0-100 的浮點數
func (f *LogFile) GetErrorRate() float64 {
	if f.TotalLines == 0 {
		return 0
	}
	return float64(f.ErrorLines) / float64(f.TotalLines) * 100
}

// GetParseSpeed 取得解析速度（MB/秒）
// 根據檔案大小和解析時間計算
func (f *LogFile) GetParseSpeed() float64 {
	if f.ParseTime == 0 {
		return 0
	}
	sizeMB := float64(f.Size) / (1024 * 1024)
	timeSeconds := float64(f.ParseTime) / 1000.0
	return sizeMB / timeSeconds
}

// GetMemoryEfficiency 取得記憶體效率（記憶體使用/檔案大小）
// 返回倍數（理想值 ≤ 1.2）
func (f *LogFile) GetMemoryEfficiency() float64 {
	if f.Size == 0 {
		return 0
	}
	return float64(f.MemoryUsed) / float64(f.Size)
}

// GetEntryCount 取得記錄總數
// 便捷方法，等同於 len(Entries)
func (f *LogFile) GetEntryCount() int {
	return len(f.Entries)
}

// GetValidEntries 取得所有成功解析的記錄
// 過濾掉包含解析錯誤的記錄
func (f *LogFile) GetValidEntries() []LogEntry {
	valid := make([]LogEntry, 0, f.ParsedLines)
	for _, entry := range f.Entries {
		if entry.ParseError == "" {
			valid = append(valid, entry)
		}
	}
	return valid
}

// GetErrorEntries 取得所有解析失敗的記錄
// 用於錯誤報告和調試
func (f *LogFile) GetErrorEntries() []LogEntry {
	errors := make([]LogEntry, 0, f.ErrorLines)
	for _, entry := range f.Entries {
		if entry.ParseError != "" {
			errors = append(errors, entry)
		}
	}
	return errors
}

// HasErrors 檢查是否存在解析錯誤
// 快速檢查方法
func (f *LogFile) HasErrors() bool {
	return f.ErrorLines > 0
}
