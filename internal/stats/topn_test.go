package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTopNHeap_基本功能 測試 Top-N 堆積的基本插入和排序功能
func TestTopNHeap_基本功能(t *testing.T) {
	// 建立一個 Top 3 的堆積
	heap := NewTopNHeap(3)

	// 插入測試數據
	items := []struct {
		key   string
		value int
	}{
		{"A", 10},
		{"B", 5},
		{"C", 20},
		{"D", 15},
		{"E", 3},
	}

	for _, item := range items {
		heap.Push(item.key, item.value)
	}

	// 獲取結果
	results := heap.GetResults()

	// 驗證結果數量
	assert.Equal(t, 3, len(results), "應該只保留 Top 3 個項目")

	// 驗證結果順序（降序）
	assert.Equal(t, "C", results[0].Key, "第一名應該是 C")
	assert.Equal(t, 20, results[0].Count, "第一名的值應該是 20")

	assert.Equal(t, "D", results[1].Key, "第二名應該是 D")
	assert.Equal(t, 15, results[1].Count, "第二名的值應該是 15")

	assert.Equal(t, "A", results[2].Key, "第三名應該是 A")
	assert.Equal(t, 10, results[2].Count, "第三名的值應該是 10")
}

// TestTopNHeap_空堆積 測試空堆積的情況
func TestTopNHeap_空堆積(t *testing.T) {
	heap := NewTopNHeap(10)
	results := heap.GetResults()

	assert.Empty(t, results, "空堆積應該返回空結果")
}

// TestTopNHeap_單一元素 測試只有一個元素的情況
func TestTopNHeap_單一元素(t *testing.T) {
	heap := NewTopNHeap(5)
	heap.Push("A", 100)

	results := heap.GetResults()

	assert.Equal(t, 1, len(results), "應該有一個元素")
	assert.Equal(t, "A", results[0].Key)
	assert.Equal(t, 100, results[0].Count)
}

// TestTopNHeap_相同值 測試相同值的排序
func TestTopNHeap_相同值(t *testing.T) {
	heap := NewTopNHeap(3)

	heap.Push("A", 10)
	heap.Push("B", 10)
	heap.Push("C", 10)
	heap.Push("D", 5)

	results := heap.GetResults()

	assert.Equal(t, 3, len(results), "應該有 3 個結果")

	// 所有結果的值都應該是 10
	for _, result := range results {
		assert.Equal(t, 10, result.Count, "所有 Top 3 的值都應該是 10")
	}
}

// TestTopNHeap_大量數據 測試大量數據的性能和正確性
func TestTopNHeap_大量數據(t *testing.T) {
	heap := NewTopNHeap(10)

	// 插入 1000 個項目
	for i := 1; i <= 1000; i++ {
		key := string(rune('A' + (i % 26)))
		heap.Push(key, i)
	}

	results := heap.GetResults()

	// 應該只保留 Top 10
	assert.Equal(t, 10, len(results), "應該只有 Top 10 個結果")

	// 驗證結果是降序排列
	for i := 0; i < len(results)-1; i++ {
		assert.GreaterOrEqual(t, results[i].Count, results[i+1].Count,
			"結果應該按降序排列")
	}

	// 驗證最大值
	assert.Equal(t, 1000, results[0].Count, "最大值應該是 1000")
}

// TestTopNHeap_零值 測試零值和負值的處理
func TestTopNHeap_零值和負值(t *testing.T) {
	heap := NewTopNHeap(5)

	heap.Push("A", 10)
	heap.Push("B", 0)
	heap.Push("C", -5)
	heap.Push("D", 20)

	results := heap.GetResults()

	// 應該按值排序
	assert.Equal(t, "D", results[0].Key, "最大值應該是 D")
	assert.Equal(t, "A", results[1].Key, "第二大應該是 A")
	assert.Equal(t, "B", results[2].Key, "第三大應該是 B (0)")
	assert.Equal(t, "C", results[3].Key, "最小值應該是 C (-5)")
}

// TestTopNHeap_更新現有項目 測試更新已存在項目的情況
func TestTopNHeap_更新現有項目(t *testing.T) {
	heap := NewTopNHeap(3)

	// 第一次插入
	heap.Push("A", 10)
	heap.Push("B", 5)
	heap.Push("C", 15)

	// 更新 B 的值
	heap.Push("B", 20)

	results := heap.GetResults()

	// B 應該變成最大值
	assert.Equal(t, "B", results[0].Key, "B 應該是第一名")
	assert.Equal(t, 20, results[0].Count, "B 的值應該是 20（更新後）")
}

// BenchmarkTopNHeap_插入 測試插入性能
func BenchmarkTopNHeap_插入(b *testing.B) {
	heap := NewTopNHeap(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push("key", i)
	}
}

// BenchmarkTopNHeap_大量數據 測試大量數據的整體性能
func BenchmarkTopNHeap_大量數據(b *testing.B) {
	for i := 0; i < b.N; i++ {
		heap := NewTopNHeap(10)

		// 插入 100 萬個項目
		for j := 0; j < 1000000; j++ {
			heap.Push("key", j)
		}

		heap.GetResults()
	}
}
