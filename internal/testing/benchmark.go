package testing

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// BenchmarkResult 基準測試結果
type BenchmarkResult struct {
	Name           string        // 測試名稱
	Iterations     int           // 迭代次數
	Duration       time.Duration // 總耗時
	AvgDuration    time.Duration // 平均耗時
	OpsPerSecond   float64       // 每秒操作數
	MemAllocBytes  uint64        // 記憶體分配（bytes）
	MemAllocsCount uint64        // 記憶體分配次數
}

// String 返回格式化的基準測試結果
func (r *BenchmarkResult) String() string {
	return fmt.Sprintf(
		"%s: %d 次迭代, 平均 %v, %.2f ops/s, 記憶體 %d bytes (%d 次分配)",
		r.Name,
		r.Iterations,
		r.AvgDuration,
		r.OpsPerSecond,
		r.MemAllocBytes,
		r.MemAllocsCount,
	)
}

// RunBenchmark 執行自訂基準測試
// name: 測試名稱
// fn: 要測試的函式
// iterations: 迭代次數
func RunBenchmark(name string, fn func(), iterations int) *BenchmarkResult {
	// 強制 GC
	runtime.GC()

	// 記錄初始記憶體狀態
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// 執行測試
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	duration := time.Since(start)

	// 記錄最終記憶體狀態
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// 計算結果
	result := &BenchmarkResult{
		Name:           name,
		Iterations:     iterations,
		Duration:       duration,
		AvgDuration:    duration / time.Duration(iterations),
		OpsPerSecond:   float64(iterations) / duration.Seconds(),
		MemAllocBytes:  memAfter.TotalAlloc - memBefore.TotalAlloc,
		MemAllocsCount: memAfter.Mallocs - memBefore.Mallocs,
	}

	return result
}

// BenchmarkThroughput 測試吞吐量
// name: 測試名稱
// fn: 要測試的函式，返回處理的數據量（bytes）
// duration: 測試持續時間
func BenchmarkThroughput(name string, fn func() int64, duration time.Duration) *ThroughputResult {
	start := time.Now()
	var totalBytes int64
	var iterations int

	for time.Since(start) < duration {
		totalBytes += fn()
		iterations++
	}

	elapsed := time.Since(start)

	return &ThroughputResult{
		Name:        name,
		Duration:    elapsed,
		TotalBytes:  totalBytes,
		Iterations:  iterations,
		BytesPerSec: float64(totalBytes) / elapsed.Seconds(),
		MBPerSec:    float64(totalBytes) / (1024 * 1024) / elapsed.Seconds(),
	}
}

// ThroughputResult 吞吐量測試結果
type ThroughputResult struct {
	Name        string        // 測試名稱
	Duration    time.Duration // 測試持續時間
	TotalBytes  int64         // 總處理字節數
	Iterations  int           // 迭代次數
	BytesPerSec float64       // 每秒字節數
	MBPerSec    float64       // 每秒 MB 數
}

// String 返回格式化的吞吐量測試結果
func (r *ThroughputResult) String() string {
	return fmt.Sprintf(
		"%s: %.2f MB/s (%d 次迭代, %d bytes, 耗時 %v)",
		r.Name,
		r.MBPerSec,
		r.Iterations,
		r.TotalBytes,
		r.Duration,
	)
}

// AssertBenchmarkTarget 斷言基準測試符合目標
func AssertBenchmarkTarget(t *testing.T, result *BenchmarkResult, maxAvgDuration time.Duration, msg string) {
	t.Helper()
	if result.AvgDuration > maxAvgDuration {
		t.Errorf("%s: 平均耗時 %v 超過目標 %v", msg, result.AvgDuration, maxAvgDuration)
	}
}

// AssertThroughputTarget 斷言吞吐量符合目標
func AssertThroughputTarget(t *testing.T, result *ThroughputResult, minMBPerSec float64, msg string) {
	t.Helper()
	if result.MBPerSec < minMBPerSec {
		t.Errorf("%s: 吞吐量 %.2f MB/s 低於目標 %.2f MB/s", msg, result.MBPerSec, minMBPerSec)
	}
}

// CompareResults 比較兩個基準測試結果
// 返回性能提升百分比（正數表示更快，負數表示更慢）
func CompareResults(before, after *BenchmarkResult) float64 {
	improvement := float64(before.AvgDuration-after.AvgDuration) / float64(before.AvgDuration) * 100
	return improvement
}

// PrintBenchmarkSummary 打印基準測試摘要
func PrintBenchmarkSummary(results []*BenchmarkResult) {
	fmt.Println("\n=== 基準測試摘要 ===")
	for _, r := range results {
		fmt.Println(r.String())
	}
}

// PrintThroughputSummary 打印吞吐量測試摘要
func PrintThroughputSummary(results []*ThroughputResult) {
	fmt.Println("\n=== 吞吐量測試摘要 ===")
	for _, r := range results {
		fmt.Println(r.String())
	}
}
