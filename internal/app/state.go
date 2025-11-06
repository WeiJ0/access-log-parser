package app

import (
	"sync"
	"access-log-analyzer/internal/models"
)

// State 管理應用程式的全局狀態
// 追蹤已開啟的檔案、頁籤和目前選擇的資料
type State struct {
	mu           sync.RWMutex
	openFiles    map[string]*models.LogFile     // 檔案路徑 -> LogFile
	tabs         []string                       // 已開啟的頁籤順序（檔案路徑列表）
	activeTab    string                         // 目前活動頁籤的檔案路徑
	selectedRows map[string][]int               // 檔案路徑 -> 選中的資料列索引
}

// NewState 建立新的 State 實例
// 初始化所有狀態映射和列表
func NewState() *State {
	return &State{
		openFiles:    make(map[string]*models.LogFile),
		tabs:         make([]string, 0),
		selectedRows: make(map[string][]int),
	}
}

// AddFile 新增已開啟的檔案
// 如果檔案已存在則更新，否則新增至頁籤列表
func (s *State) AddFile(path string, logFile *models.LogFile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.openFiles[path]; !exists {
		s.tabs = append(s.tabs, path)
	}
	s.openFiles[path] = logFile
	s.activeTab = path
}

// RemoveFile 移除已開啟的檔案
// 同時清理相關的頁籤和選擇狀態
func (s *State) RemoveFile(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delete(s.openFiles, path)
	delete(s.selectedRows, path)
	
	// 從頁籤列表中移除
	for i, tab := range s.tabs {
		if tab == path {
			s.tabs = append(s.tabs[:i], s.tabs[i+1:]...)
			break
		}
	}
	
	// 如果刪除的是活動頁籤，切換到其他頁籤
	if s.activeTab == path {
		if len(s.tabs) > 0 {
			s.activeTab = s.tabs[len(s.tabs)-1]
		} else {
			s.activeTab = ""
		}
	}
}

// GetFile 取得指定路徑的檔案資料
// 執行緒安全的讀取操作
func (s *State) GetFile(path string) (*models.LogFile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	file, exists := s.openFiles[path]
	return file, exists
}

// GetActiveFile 取得目前活動頁籤的檔案資料
// 如果沒有活動頁籤則返回 nil
func (s *State) GetActiveFile() (*models.LogFile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.activeTab == "" {
		return nil, false
	}
	
	file, exists := s.openFiles[s.activeTab]
	return file, exists
}

// SetActiveTab 設定目前活動頁籤
// 只有當檔案已開啟時才會設定成功
func (s *State) SetActiveTab(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.openFiles[path]; !exists {
		return false
	}
	
	s.activeTab = path
	return true
}

// GetActiveTab 取得目前活動頁籤的路徑
// 執行緒安全的讀取操作
func (s *State) GetActiveTab() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.activeTab
}

// GetTabs 取得所有頁籤的路徑列表（複製）
// 返回新的切片避免外部修改
func (s *State) GetTabs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	tabs := make([]string, len(s.tabs))
	copy(tabs, s.tabs)
	return tabs
}

// SetSelectedRows 設定指定檔案的選中資料列
// 用於追蹤使用者在 UI 中的選擇狀態
func (s *State) SetSelectedRows(path string, rows []int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.openFiles[path]; !exists {
		return
	}
	
	s.selectedRows[path] = rows
}

// GetSelectedRows 取得指定檔案的選中資料列（複製）
// 返回新的切片避免外部修改
func (s *State) GetSelectedRows(path string) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	rows, exists := s.selectedRows[path]
	if !exists {
		return []int{}
	}
	
	result := make([]int, len(rows))
	copy(result, rows)
	return result
}

// Cleanup 清理所有狀態資料
// 在應用程式關閉時調用
func (s *State) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.openFiles = make(map[string]*models.LogFile)
	s.tabs = make([]string, 0)
	s.selectedRows = make(map[string][]int)
	s.activeTab = ""
}

// FileCount 返回已開啟的檔案數量
// 執行緒安全的讀取操作
func (s *State) FileCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return len(s.openFiles)
}
