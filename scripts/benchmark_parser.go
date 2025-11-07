package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"access-log-analyzer/internal/parser"
	"access-log-analyzer/pkg/logger"
)

// benchmark_parser.go
// 用於驗證解析器效能的獨立測試程式

func main() {
	// 初始化 logger
	logger.Init()

	// 測試檔案列表
	testFiles := []struct {
		path string
		name string
	}{
		{"testdata/valid.log", "小型檔案"},
		{"testdata/100mb.log", "100MB 檔案"},
	}

	fmt.Println("=== Apache Log Parser 效能測試 ===\n")

	for _, tf := range testFiles {
		// 檢查檔案是否存在
		if _, err := os.Stat(tf.path); os.IsNotExist(err) {
			fmt.Printf("⚠️  跳過 %s: 檔案不存在\n\n", tf.name)
			continue
		}

		// 獲取檔案大小
		fileInfo, err := os.Stat(tf.path)
		if err != nil {
			fmt.Printf("❌ 無法讀取檔案資訊 %s: %v\n\n", tf.name, err)
			continue
		}
		fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

		fmt.Printf("測試檔案: %s\n", tf.name)
		fmt.Printf("檔案路徑: %s\n", tf.path)
		fmt.Printf("檔案大小: %.2f MB\n", fileSizeMB)

		// 記錄初始記憶體使用
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)
		initialMemoryMB := float64(m1.Alloc) / (1024 * 1024)

		// 解析檔案
		p := parser.NewParser(parser.FormatCombined, 0) // 0 = 使用 CPU 核心數

		startTime := time.Now()
		result, err := p.ParseFile(tf.path, fileInfo.Size())
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ 解析失敗: %v\n\n", err)
			continue
		}

		// 記錄最終記憶體使用
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)
		finalMemoryMB := float64(m2.Alloc) / (1024 * 1024)
		memoryUsedMB := finalMemoryMB - initialMemoryMB

		// 計算效能指標
		throughputMBps := fileSizeMB / duration.Seconds()
		linesPerSec := float64(result.TotalLines) / duration.Seconds()
		memoryRatio := memoryUsedMB / fileSizeMB

		// 顯示結果
		fmt.Printf("\n結果:\n")
		fmt.Printf("  總行數: %d\n", result.TotalLines)
		fmt.Printf("  成功解析: %d (%.2f%%)\n", result.ParsedLines, float64(result.ParsedLines)/float64(result.TotalLines)*100)
		fmt.Printf("  解析失敗: %d (%.2f%%)\n", result.ErrorLines, float64(result.ErrorLines)/float64(result.TotalLines)*100)
		fmt.Printf("  解析時間: %.2f 秒\n", duration.Seconds())
		fmt.Printf("  吞吐量: %.2f MB/秒\n", throughputMBps)
		fmt.Printf("  每秒行數: %.0f 行/秒\n", linesPerSec)
		fmt.Printf("  記憶體使用: %.2f MB (%.2fx 檔案大小)\n", memoryUsedMB, memoryRatio)

		// 效能判定
		fmt.Printf("\n效能評估:\n")

		// 吞吐量目標: 60-80 MB/秒
		if throughputMBps >= 60 {
			fmt.Printf("  ✅ 吞吐量達標 (>= 60 MB/秒)\n")
		} else {
			fmt.Printf("  ❌ 吞吐量未達標 (< 60 MB/秒)\n")
		}

		// 記憶體目標: <= 1.2x 檔案大小
		if memoryRatio <= 1.2 {
			fmt.Printf("  ✅ 記憶體使用達標 (<= 1.2x 檔案大小)\n")
		} else {
			fmt.Printf("  ❌ 記憶體使用超標 (> 1.2x 檔案大小)\n")
		}

		// 100MB 檔案特定要求: <= 2 秒
		if tf.name == "100MB 檔案" {
			if duration.Seconds() <= 2 {
				fmt.Printf("  ✅ 100MB 檔案解析時間達標 (<= 2 秒)\n")
			} else {
				fmt.Printf("  ❌ 100MB 檔案解析時間超標 (> 2 秒)\n")
			}
		}

		fmt.Printf("\n錯誤樣本 (最多 5 筆):\n")
		maxErrors := 5
		if len(result.ErrorSamples) < maxErrors {
			maxErrors = len(result.ErrorSamples)
		}
		for i := 0; i < maxErrors; i++ {
			sample := result.ErrorSamples[i]
			fmt.Printf("  行 %d: %s\n", sample.LineNumber, sample.Line)
			if sample.Error != "" {
				fmt.Printf("    錯誤: %s\n", sample.Error)
			}
		}

		fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
	}

	// 整體總結
	fmt.Println("測試完成！")
}
