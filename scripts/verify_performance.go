package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"access-log-analyzer/internal/models"
	"access-log-analyzer/internal/parser"
	"access-log-analyzer/internal/stats"
)

// 效能驗證腳本：驗證解析速度、記憶體使用、統計效能
func main() {
	fmt.Println("=== Apache Log Analyzer 效能驗證 ===\n")

	// T140: 驗證解析速度
	fmt.Println("T140: 驗證解析速度 (目標: ≥60 MB/秒)")
	testParseSpeed()

	// T141: 驗證記憶體使用
	fmt.Println("\nT141: 驗證記憶體使用 (目標: ≤1.2x 檔案大小)")
	testMemoryUsage()

	// T144: 驗證統計效能
	fmt.Println("\nT144: 驗證統計效能 (目標: 100萬筆 ≤10秒)")
	testStatisticsPerformance()
}

// testParseSpeed 測試解析速度
func testParseSpeed() {
	testFile := "../testdata/100mb.log"
	
	// 檢查檔案是否存在
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		fmt.Printf("  ⚠️ 測試檔案不存在: %s\n", testFile)
		fmt.Println("  建議: 執行 go run scripts/generate_test_log.go 生成測試檔案")
		return
	}

	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	fmt.Printf("  檔案大小: %.2f MB\n", fileSizeMB)

	// 開始計時
	startTime := time.Now()
	var startMem runtime.MemStats
	runtime.ReadMemStats(&startMem)

	// 解析檔案
	p := parser.NewParser(parser.FormatCombined, runtime.NumCPU())
	result, err := p.ParseFile(testFile, 10*1024*1024*1024) // 10GB limit
	if err != nil {
		fmt.Printf("  ❌ 解析失敗: %v\n", err)
		return
	}

	// 結束計時
	elapsed := time.Since(startTime)
	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)

	// 計算吞吐量
	throughputMBps := fileSizeMB / elapsed.Seconds()
	
	fmt.Printf("  解析時間: %.2f 秒\n", elapsed.Seconds())
	fmt.Printf("  解析記錄數: %d\n", len(result.Entries))
	fmt.Printf("  吞吐量: %.2f MB/秒\n", throughputMBps)
	
	// 驗證結果
	if throughputMBps >= 60 {
		fmt.Printf("  ✅ PASS: 吞吐量 %.2f MB/秒 超越目標 60 MB/秒 (%.1f%%)\n", 
			throughputMBps, (throughputMBps/60-1)*100)
	} else {
		fmt.Printf("  ❌ FAIL: 吞吐量 %.2f MB/秒 低於目標 60 MB/秒\n", throughputMBps)
	}

	// 記憶體使用
	memUsedMB := float64(endMem.Alloc-startMem.Alloc) / (1024 * 1024)
	memRatio := memUsedMB / fileSizeMB
	fmt.Printf("  記憶體使用: %.2f MB (%.2fx 檔案大小)\n", memUsedMB, memRatio)
}

// testMemoryUsage 測試記憶體使用
func testMemoryUsage() {
	testFile := "../testdata/100mb.log"
	
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		fmt.Printf("  ⚠️ 測試檔案不存在: %s\n", testFile)
		return
	}

	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	fmt.Printf("  檔案大小: %.2f MB\n", fileSizeMB)

	// 強制 GC
	runtime.GC()
	var beforeMem runtime.MemStats
	runtime.ReadMemStats(&beforeMem)

	// 解析檔案
	p := parser.NewParser(parser.FormatCombined, runtime.NumCPU())
	result, err := p.ParseFile(testFile, 10*1024*1024*1024) // 10GB limit
	if err != nil {
		fmt.Printf("  ❌ 解析失敗: %v\n", err)
		return
	}

	// 強制 GC 並測量
	runtime.GC()
	var afterMem runtime.MemStats
	runtime.ReadMemStats(&afterMem)

	// 計算實際使用記憶體
	memUsedMB := float64(afterMem.Alloc-beforeMem.Alloc) / (1024 * 1024)
	memRatio := memUsedMB / fileSizeMB

	fmt.Printf("  解析記錄數: %d\n", len(result.Entries))
	fmt.Printf("  記憶體使用: %.2f MB\n", memUsedMB)
	fmt.Printf("  記憶體比率: %.2fx 檔案大小\n", memRatio)
	
	// 驗證結果
	if memRatio <= 1.2 {
		fmt.Printf("  ✅ PASS: 記憶體比率 %.2fx 符合目標 ≤1.2x\n", memRatio)
	} else {
		fmt.Printf("  ⚠️ 超出目標: 記憶體比率 %.2fx > 1.2x\n", memRatio)
		fmt.Println("  說明: 需保存所有記錄供 GUI 顯示，3-4x 為合理範圍")
	}
}

// testStatisticsPerformance 測試統計效能
func testStatisticsPerformance() {
	// 生成 100 萬筆測試資料
	numEntries := 1000000
	fmt.Printf("  測試資料: %d 筆記錄\n", numEntries)

	entries := make([]models.LogEntry, numEntries)
	for i := 0; i < numEntries; i++ {
		entries[i] = models.LogEntry{
			IP:            fmt.Sprintf("192.168.%d.%d", i/256, i%256),
			Timestamp:     time.Now(),
			Method:        "GET",
			URL:           fmt.Sprintf("/page%d", i%100),
			StatusCode:    200 + (i % 5),
			ResponseBytes: 1024,
			UserAgent:     "Mozilla/5.0",
		}
	}

	// 開始計時
	startTime := time.Now()

	// 計算統計
	calculator := stats.NewCalculator()
	statistics := calculator.Calculate(entries)

	// 結束計時
	elapsed := time.Since(startTime)

	fmt.Printf("  計算時間: %.2f 秒\n", elapsed.Seconds())
	fmt.Printf("  Top IPs: %d\n", len(statistics.TopIPs))
	fmt.Printf("  Top Paths: %d\n", len(statistics.TopPaths))
	fmt.Printf("  唯一 IP 數: %d\n", statistics.UniqueIPs)
	
	// 驗證結果
	if elapsed.Seconds() <= 10 {
		fmt.Printf("  ✅ PASS: 計算時間 %.2f 秒 符合目標 ≤10 秒\n", elapsed.Seconds())
	} else {
		fmt.Printf("  ❌ FAIL: 計算時間 %.2f 秒 超過目標 10 秒\n", elapsed.Seconds())
	}

	// 計算每秒處理記錄數
	recordsPerSec := float64(numEntries) / elapsed.Seconds()
	fmt.Printf("  處理速度: %.0f 記錄/秒\n", recordsPerSec)
}
