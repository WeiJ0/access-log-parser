package models

import "time"

// PerformanceMetrics 記錄解析和處理過程的效能指標
// 用於監控和優化系統效能
type PerformanceMetrics struct {
	// 解析效能
	StartTime      time.Time `json:"startTime"`      // 開始時間
	EndTime        time.Time `json:"endTime"`        // 結束時間
	Duration       int64     `json:"duration"`       // 總耗時(毫秒)
	LinesProcessed int64     `json:"linesProcessed"` // 處理的行數
	LinesPerSecond float64   `json:"linesPerSecond"` // 每秒處理行數

	// 記憶體使用
	MemoryUsedMB float64 `json:"memoryUsedMB"` // 使用的記憶體(MB)
	MemoryPeakMB float64 `json:"memoryPeakMB"` // 峰值記憶體(MB)
	MemoryRatio  float64 `json:"memoryRatio"`  // 記憶體使用率 (相對於檔案大小)

	// 檔案資訊
	FileSizeMB     float64 `json:"fileSizeMB"`     // 檔案大小(MB)
	ThroughputMBPS float64 `json:"throughputMBPS"` // 吞吐量(MB/秒)

	// 錯誤統計
	ErrorCount int     `json:"errorCount"` // 錯誤數量
	ErrorRate  float64 `json:"errorRate"`  // 錯誤率(%)

	// Goroutine 資訊
	WorkerCount   int     `json:"workerCount"`   // Worker goroutine 數量
	AvgWorkerLoad float64 `json:"avgWorkerLoad"` // 平均 Worker 負載
}

// Calculate 計算所有衍生指標
// 應在設定基本值後呼叫
func (m *PerformanceMetrics) Calculate() {
	// 計算總耗時(毫秒)
	m.Duration = m.EndTime.Sub(m.StartTime).Milliseconds()

	// 計算每秒處理行數
	if m.Duration > 0 {
		seconds := float64(m.Duration) / 1000.0
		m.LinesPerSecond = float64(m.LinesProcessed) / seconds

		// 計算吞吐量(MB/秒)
		m.ThroughputMBPS = m.FileSizeMB / seconds
	}

	// 計算記憶體使用率
	if m.FileSizeMB > 0 {
		m.MemoryRatio = m.MemoryPeakMB / m.FileSizeMB
	}

	// 計算錯誤率
	if m.LinesProcessed > 0 {
		m.ErrorRate = (float64(m.ErrorCount) / float64(m.LinesProcessed)) * 100
	}
}

// IsAcceptable 檢查效能指標是否符合目標要求
// 根據憲法定義的效能標準: 60-80 MB/秒, 記憶體 ≤1.2x 檔案大小
func (m *PerformanceMetrics) IsAcceptable() bool {
	// 檢查吞吐量: 應該 ≥60 MB/秒
	if m.ThroughputMBPS < 60 {
		return false
	}

	// 檢查記憶體使用率: 應該 ≤1.2x
	if m.MemoryRatio > 1.2 {
		return false
	}

	return true
}
