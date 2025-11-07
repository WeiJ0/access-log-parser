package stats

import (
	"container/heap"
	"sort"
)

// TopNItem 代表 Top-N 結果中的一個項目
type TopNItem struct {
	Key   string `json:"key"`   // 項目鍵值（如 IP 位址、路徑等）
	Count int    `json:"count"` // 計數值
}

// TopNHeap 使用最小堆積實現高效的 Top-N 追蹤
// 時間複雜度：O(N log K)，其中 N 是總項目數，K 是 Top-N 的 N
// 空間複雜度：O(K)
type TopNHeap struct {
	n     int              // 保留前 N 個項目
	items *minHeap         // 最小堆積
	index map[string]*item // 快速查找已存在的項目
}

// item 是堆積中的內部節點
type item struct {
	key   string
	count int
	index int // 在堆積中的索引
}

// minHeap 實現 heap.Interface
type minHeap []*item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *minHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*item)
	item.index = n
	*h = append(*h, item)
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免記憶體洩漏
	item.index = -1 // 標記為已移除
	*h = old[0 : n-1]
	return item
}

// NewTopNHeap 建立新的 Top-N 堆積
// n: 保留前 N 個項目
func NewTopNHeap(n int) *TopNHeap {
	h := &minHeap{}
	heap.Init(h)

	return &TopNHeap{
		n:     n,
		items: h,
		index: make(map[string]*item),
	}
}

// Push 添加或更新一個項目
// 如果項目已存在，則更新其計數值
// 如果堆積未滿，直接添加
// 如果堆積已滿且新值大於最小值，則替換最小值
func (t *TopNHeap) Push(key string, count int) {
	// 檢查項目是否已存在
	if existingItem, exists := t.index[key]; exists {
		// 更新現有項目的計數
		existingItem.count = count
		heap.Fix(t.items, existingItem.index)
		return
	}

	// 建立新項目
	newItem := &item{
		key:   key,
		count: count,
	}

	// 如果堆積未滿，直接添加
	if t.items.Len() < t.n {
		heap.Push(t.items, newItem)
		t.index[key] = newItem
		return
	}

	// 堆積已滿，檢查是否需要替換最小值
	minItem := (*t.items)[0]
	if count > minItem.count {
		// 移除最小值
		delete(t.index, minItem.key)
		heap.Pop(t.items)

		// 添加新項目
		heap.Push(t.items, newItem)
		t.index[key] = newItem
	}
}

// GetResults 返回 Top-N 結果，按計數降序排列
func (t *TopNHeap) GetResults() []TopNItem {
	// 複製堆積內容以避免修改原始數據
	results := make([]TopNItem, 0, t.items.Len())
	for _, item := range *t.items {
		results = append(results, TopNItem{
			Key:   item.key,
			Count: item.count,
		})
	}

	// 按計數降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	return results
}

// Size 返回當前堆積中的項目數量
func (t *TopNHeap) Size() int {
	return t.items.Len()
}

// Clear 清空堆積
func (t *TopNHeap) Clear() {
	t.items = &minHeap{}
	heap.Init(t.items)
	t.index = make(map[string]*item)
}
