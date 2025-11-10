package monitoring

import (
	"runtime"
	"sync"
	"time"

	"access-log-analyzer/internal/models"
)

// MetricsCollector 效能指標收集器
// 追蹤解析過程中的效能指標（吞吐量、記憶體使用等）
type MetricsCollector struct {
	mu sync.RWMutex

	// 基本計數
	startTime      time.Time
	endTime        time.Time
	linesProcessed int64
	errorCount     int
	fileSizeBytes  int64

	// 記憶體追蹤
	memStart uint64 // 開始時的記憶體使用（bytes）
	memPeak  uint64 // 峰值記憶體使用（bytes）

	// Worker 資訊
	workerCount int
	workerLoads []int64 // 每個 worker 處理的行數

	// 是否正在收集
	collecting bool
}

// NewMetricsCollector 建立新的指標收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		workerLoads: make([]int64, 0),
		collecting:  false,
	}
}

// Start 開始收集效能指標
// 記錄開始時間和初始記憶體狀態
func (m *MetricsCollector) Start(fileSizeBytes int64, workerCount int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.startTime = time.Now()
	m.fileSizeBytes = fileSizeBytes
	m.workerCount = workerCount
	m.linesProcessed = 0
	m.errorCount = 0
	m.collecting = true

	// 記錄初始記憶體使用
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.memStart = memStats.Alloc
	m.memPeak = memStats.Alloc

	// 初始化 worker 負載追蹤
	m.workerLoads = make([]int64, workerCount)
}

// Stop 停止收集效能指標
// 記錄結束時間
func (m *MetricsCollector) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.endTime = time.Now()
	m.collecting = false

	// 更新最終記憶體峰值
	m.updateMemoryPeak()
}

// IncrementLines 增加已處理行數
// workerID: worker 的索引（用於追蹤負載分布）
func (m *MetricsCollector) IncrementLines(workerID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.linesProcessed++

	// 更新 worker 負載
	if workerID >= 0 && workerID < len(m.workerLoads) {
		m.workerLoads[workerID]++
	}

	// 定期更新記憶體峰值（每 1000 行）
	if m.linesProcessed%1000 == 0 {
		m.updateMemoryPeak()
	}
}

// IncrementErrors 增加錯誤計數
func (m *MetricsCollector) IncrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCount++
}

// updateMemoryPeak 更新記憶體峰值（不加鎖，須由呼叫者加鎖）
func (m *MetricsCollector) updateMemoryPeak() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	currentMem := memStats.Alloc
	if currentMem > m.memPeak {
		m.memPeak = currentMem
	}
}

// GetMetrics 取得當前的效能指標
// 返回 PerformanceMetrics 結構
func (m *MetricsCollector) GetMetrics() *models.PerformanceMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 確保獲取最新的記憶體峰值
	m.mu.RUnlock()
	m.mu.Lock()
	m.updateMemoryPeak()
	m.mu.Unlock()
	m.mu.RLock()

	metrics := &models.PerformanceMetrics{
		StartTime:      m.startTime,
		EndTime:        m.endTime,
		LinesProcessed: m.linesProcessed,
		ErrorCount:     m.errorCount,
		WorkerCount:    m.workerCount,

		// 轉換記憶體單位為 MB
		MemoryUsedMB: float64(m.memPeak-m.memStart) / (1024 * 1024),
		MemoryPeakMB: float64(m.memPeak) / (1024 * 1024),
		FileSizeMB:   float64(m.fileSizeBytes) / (1024 * 1024),
	}

	// 計算平均 worker 負載
	if m.workerCount > 0 {
		var totalLoad int64
		for _, load := range m.workerLoads {
			totalLoad += load
		}
		metrics.AvgWorkerLoad = float64(totalLoad) / float64(m.workerCount)
	}

	// 計算所有衍生指標
	metrics.Calculate()

	return metrics
}

// IsCollecting 檢查是否正在收集指標
func (m *MetricsCollector) IsCollecting() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.collecting
}

// Reset 重置所有指標
// 用於重新開始收集
func (m *MetricsCollector) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.startTime = time.Time{}
	m.endTime = time.Time{}
	m.linesProcessed = 0
	m.errorCount = 0
	m.fileSizeBytes = 0
	m.memStart = 0
	m.memPeak = 0
	m.workerCount = 0
	m.workerLoads = make([]int64, 0)
	m.collecting = false
}

// GetWorkerLoads 取得每個 worker 的負載分布
// 返回每個 worker 處理的行數
func (m *MetricsCollector) GetWorkerLoads() []int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本以避免並發修改
	loads := make([]int64, len(m.workerLoads))
	copy(loads, m.workerLoads)
	return loads
}
