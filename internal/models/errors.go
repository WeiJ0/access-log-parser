package models

import "fmt"

// ParseError 表示日誌解析過程中的錯誤
// 包含行號、原始行內容和錯誤原因
type ParseError struct {
	LineNumber int    // 行號（從1開始）
	RawLine    string // 原始行內容
	Reason     string // 錯誤原因描述
}

// Error 實現 error 介面
// 返回格式化的錯誤訊息
func (e *ParseError) Error() string {
	return fmt.Sprintf("解析錯誤於第 %d 行: %s (原始: %s)", e.LineNumber, e.Reason, e.RawLine)
}

// ValidationError 表示資料驗證錯誤
// 用於驗證日誌格式、檔案大小等
type ValidationError struct {
	Field   string // 驗證失敗的欄位名稱
	Value   string // 驗證失敗的值
	Message string // 錯誤訊息
}

// Error 實現 error 介面
// 返回格式化的驗證錯誤訊息
func (e *ValidationError) Error() string {
	return fmt.Sprintf("驗證錯誤 [%s]: %s (值: %s)", e.Field, e.Message, e.Value)
}

// FileError 表示檔案相關錯誤
// 包含檔案路徑和具體錯誤原因
type FileError struct {
	Path    string // 檔案路徑
	Message string // 錯誤訊息
}

// Error 實現 error 介面
// 返回格式化的檔案錯誤訊息
func (e *FileError) Error() string {
	return fmt.Sprintf("檔案錯誤 [%s]: %s", e.Path, e.Message)
}
