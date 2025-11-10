// T139: 1GB 完整工作流程測試
// 測試完整流程：載入 → 解析 → 統計分析 → 匯出 Excel
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"access-log-analyzer/internal/exporter"
	"access-log-analyzer/internal/models"
	"access-log-analyzer/internal/parser"
	"access-log-analyzer/internal/stats"
)

func main() {
	fmt.Println("\n=== T139: 1GB 完整工作流程測試 ===")
	fmt.Println()

	testFile := "../testdata/1gb.log"
	outputFile := "../testdata/1gb_report.xlsx"

	// 檢查測試檔案
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		fmt.Printf("錯誤：找不到測試檔案 %s\n", testFile)
		fmt.Println("請先執行：go run scripts/generate_test_log.go -lines 7000000 -output testdata/1gb.log")
		os.Exit(1)
	}

	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	fmt.Printf("測試檔案: %s\n", testFile)
	fmt.Printf("檔案大小: %.2f MB\n", fileSizeMB)
	fmt.Println()

	// 記錄總開始時間
	totalStart := time.Now()

	// 步驟 1: 載入並解析檔案
	fmt.Println("步驟 1/3: 載入並解析日誌檔案...")
	parseStart := time.Now()

	p := parser.NewParser(parser.FormatCombined, runtime.NumCPU())
	result, err := p.ParseFile(testFile, 10*1024*1024*1024) // 10GB limit
	if err != nil {
		fmt.Printf("❌ 解析失敗: %v\n", err)
		os.Exit(1)
	}

	parseElapsed := time.Since(parseStart)
	parseThroughput := fileSizeMB / parseElapsed.Seconds()

	fmt.Printf("  ✓ 解析完成\n")
	fmt.Printf("  記錄數: %d 筆\n", len(result.Entries))
	fmt.Printf("  解析時間: %.2f 秒\n", parseElapsed.Seconds())
	fmt.Printf("  吞吐量: %.2f MB/秒\n", parseThroughput)
	fmt.Println()

	// 步驟 2: 統計分析
	fmt.Println("步驟 2/3: 執行統計分析...")
	statsStart := time.Now()

	calculator := stats.NewCalculator()
	statisticsPtr := calculator.Calculate(result.Entries)

	statsElapsed := time.Since(statsStart)
	statsSpeed := float64(len(result.Entries)) / statsElapsed.Seconds()

	fmt.Printf("  ✓ 統計完成\n")
	fmt.Printf("  唯一 IP 數: %d\n", statisticsPtr.UniqueIPs)
	fmt.Printf("  Top IPs: %d 個\n", len(statisticsPtr.TopIPs))
	fmt.Printf("  Top Paths: %d 個\n", len(statisticsPtr.TopPaths))
	fmt.Printf("  計算時間: %.2f 秒\n", statsElapsed.Seconds())
	fmt.Printf("  處理速度: %.0f 筆/秒\n", statsSpeed)
	fmt.Println()

	// 步驟 3: 匯出 Excel
	fmt.Println("步驟 3/3: 匯出 Excel 報告...")
	exportStart := time.Now()

	// 刪除舊檔案（如果存在）
	os.Remove(outputFile)

	// 轉換為指標類型
	entryPointers := make([]*models.LogEntry, len(result.Entries))
	for i := range result.Entries {
		entryPointers[i] = &result.Entries[i]
	}

	exp := exporter.NewXLSXExporter()
	_, err = exp.ExportWithStatsStatistics(entryPointers, &statisticsPtr, outputFile)
	if err != nil {
		fmt.Printf("❌ 匯出失敗: %v\n", err)
		os.Exit(1)
	}

	exportElapsed := time.Since(exportStart)

	// 檢查輸出檔案
	outputInfo, err := os.Stat(outputFile)
	if err != nil {
		fmt.Printf("❌ 找不到輸出檔案: %v\n", err)
		os.Exit(1)
	}

	outputSizeMB := float64(outputInfo.Size()) / (1024 * 1024)

	fmt.Printf("  ✓ 匯出完成\n")
	fmt.Printf("  輸出檔案: %s\n", outputFile)
	fmt.Printf("  檔案大小: %.2f MB\n", outputSizeMB)
	fmt.Printf("  匯出時間: %.2f 秒\n", exportElapsed.Seconds())
	fmt.Println()

	// 總結
	totalElapsed := time.Since(totalStart)

	fmt.Println("=== 測試總結 ===")
	fmt.Println()
	fmt.Printf("總處理時間: %.2f 秒 (%.2f 分鐘)\n", totalElapsed.Seconds(), totalElapsed.Minutes())
	fmt.Println()
	fmt.Println("時間分佈:")
	parsePercent := (parseElapsed.Seconds() / totalElapsed.Seconds()) * 100
	statsPercent := (statsElapsed.Seconds() / totalElapsed.Seconds()) * 100
	exportPercent := (exportElapsed.Seconds() / totalElapsed.Seconds()) * 100

	fmt.Printf("  解析: %.2f 秒 (%.1f%%)\n", parseElapsed.Seconds(), parsePercent)
	fmt.Printf("  統計: %.2f 秒 (%.1f%%)\n", statsElapsed.Seconds(), statsPercent)
	fmt.Printf("  匯出: %.2f 秒 (%.1f%%)\n", exportElapsed.Seconds(), exportPercent)
	fmt.Println()

	// 記憶體使用情況
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	allocMB := float64(m.Alloc) / (1024 * 1024)
	sysMemMB := float64(m.Sys) / (1024 * 1024)
	memRatio := allocMB / fileSizeMB

	fmt.Println("記憶體使用:")
	fmt.Printf("  已分配: %.2f MB\n", allocMB)
	fmt.Printf("  系統記憶體: %.2f MB\n", sysMemMB)
	fmt.Printf("  記憶體比例: %.2fx 檔案大小\n", memRatio)
	fmt.Println()

	// 效能評估
	fmt.Println("效能評估:")
	fmt.Printf("  解析吞吐量: %.2f MB/秒 ", parseThroughput)
	if parseThroughput >= 60 {
		fmt.Printf("✅ (目標: ≥60 MB/秒)\n")
	} else {
		fmt.Printf("⚠️ (目標: ≥60 MB/秒)\n")
	}

	fmt.Printf("  統計速度: %.0f 筆/秒 ", statsSpeed)
	if statsSpeed >= 100000 {
		fmt.Printf("✅ (參考: 685K 筆/秒)\n")
	} else {
		fmt.Printf("⚠️ (參考: 685K 筆/秒)\n")
	}

	fmt.Printf("  總處理時間: %.2f 分鐘 ", totalElapsed.Minutes())
	if totalElapsed.Minutes() <= 5 {
		fmt.Printf("✅ (預期: ≤5 分鐘)\n")
	} else {
		fmt.Printf("⚠️ (預期: ≤5 分鐘)\n")
	}
	fmt.Println()

	// 清理測試檔案
	fmt.Println("清理測試檔案...")
	os.Remove(outputFile)

	fmt.Println("✅ T139 測試完成！")
}
