package app

import (
	"access-log-analyzer/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// State 應用程式狀態管理
// 追蹤開啟的檔案和當前分頁
type State struct {
	mu           sync.RWMutex
	openFiles    map[string]*models.LogFile // 檔案路徑 -> LogFile
	tabs         []string                   // 已開啟的頁籤順序（檔案路徑列表）
	activeTab    string                     // 目前活動頁籤的檔案路徑
	selectedRows map[string][]int           // 檔案路徑 -> 選中的資料列索引
	recentFiles  []RecentFileRecord         // T151: 最近開啟的檔案列表
}

// RecentFileRecord 最近開啟的檔案記錄（內部使用）
type RecentFileRecord struct {
	Path       string    `json:"path"`       // 檔案路徑
	Name       string    `json:"name"`       // 檔案名稱
	Size       int64     `json:"size"`       // 檔案大小
	OpenedAt   time.Time `json:"openedAt"`   // 開啟時間
	TotalLines int       `json:"totalLines"` // 總行數
}

// NewState 建立新的 State 實例
// 初始化所有狀態映射和列表
func NewState() *State {
	s := &State{
		openFiles:    make(map[string]*models.LogFile),
		tabs:         make([]string, 0),
		selectedRows: make(map[string][]int),
		recentFiles:  make([]RecentFileRecord, 0),
	}
	// 載入最近檔案列表
	s.loadRecentFiles()
	return s
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

// T151: 最近檔案列表管理

// getRecentFilesPath 取得最近檔案列表的儲存路徑
func (s *State) getRecentFilesPath() string {
	// 儲存在使用者主目錄的 .apache-log-analyzer 資料夾
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	
	configDir := filepath.Join(homeDir, ".apache-log-analyzer")
	// 確保目錄存在
	os.MkdirAll(configDir, 0755)
	
	return filepath.Join(configDir, "recent-files.json")
}

// loadRecentFiles 從磁碟載入最近檔案列表
func (s *State) loadRecentFiles() {
	path := s.getRecentFilesPath()
	if path == "" {
		return
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		// 檔案不存在是正常的（首次啟動）
		return
	}
	
	var files []RecentFileRecord
	if err := json.Unmarshal(data, &files); err != nil {
		// 解析失敗，重置列表
		return
	}
	
	s.recentFiles = files
}

// saveRecentFiles 儲存最近檔案列表到磁碟
func (s *State) saveRecentFiles() error {
	path := s.getRecentFilesPath()
	if path == "" {
		return nil
	}
	
	data, err := json.MarshalIndent(s.recentFiles, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// AddRecentFile 新增檔案到最近列表
func (s *State) AddRecentFile(logFile *models.LogFile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 建立最近檔案記錄
	recent := RecentFileRecord{
		Path:       logFile.Path,
		Name:       logFile.Name,
		Size:       logFile.Size,
		OpenedAt:   time.Now(),
		TotalLines: logFile.TotalLines,
	}
	
	// 移除重複項目（如果已存在）
	filtered := make([]RecentFileRecord, 0, len(s.recentFiles))
	for _, f := range s.recentFiles {
		if f.Path != recent.Path {
			filtered = append(filtered, f)
		}
	}
	
	// 新增到列表開頭
	s.recentFiles = append([]RecentFileRecord{recent}, filtered...)
	
	// 限制最多保留 10 個最近檔案
	const maxRecentFiles = 10
	if len(s.recentFiles) > maxRecentFiles {
		s.recentFiles = s.recentFiles[:maxRecentFiles]
	}
	
	// 儲存到磁碟
	s.saveRecentFiles()
}

// GetRecentFiles 取得最近檔案列表（複製）
func (s *State) GetRecentFiles() []RecentFileRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	files := make([]RecentFileRecord, len(s.recentFiles))
	copy(files, s.recentFiles)
	return files
}

// T152: ClearRecentFiles 清空最近檔案列表
func (s *State) ClearRecentFiles() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.recentFiles = make([]RecentFileRecord, 0)
	
	// 儲存到磁碟
	return s.saveRecentFiles()
}
