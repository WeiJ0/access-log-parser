package apachelog

import (
	"bufio"
	"io"
	"os"
)

// Reader 提供高效的 log 檔案讀取功能
// 使用 bufio.Scanner 進行逐行讀取，減少記憶體使用
type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
	lineNum int
	err     error
}

// NewReader 建立新的 log 檔案讀取器
// 使用 16KB 緩衝區大小以平衡效能和記憶體使用
func NewReader(filepath string) (*Reader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	// 設定最大單行大小為 1MB（處理超長 log 行）
	buf := make([]byte, 0, 16*1024) // 16KB 初始緩衝區
	scanner.Buffer(buf, 1024*1024)  // 最大 1MB

	return &Reader{
		file:    file,
		scanner: scanner,
		lineNum: 0,
	}, nil
}

// ReadLine 讀取下一行 log 資料
// 返回行號、行內容和是否還有更多資料
func (r *Reader) ReadLine() (lineNum int, line string, hasMore bool) {
	if r.err != nil {
		return 0, "", false
	}

	if !r.scanner.Scan() {
		r.err = r.scanner.Err()
		if r.err == nil {
			r.err = io.EOF
		}
		return 0, "", false
	}

	r.lineNum++
	return r.lineNum, r.scanner.Text(), true
}

// Close 關閉檔案並釋放資源
// 必須在讀取完成後調用以避免資源洩漏
func (r *Reader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// Error 返回讀取過程中發生的錯誤
// 如果是 EOF 則返回 nil（正常結束）
func (r *Reader) Error() error {
	if r.err == io.EOF {
		return nil
	}
	return r.err
}

// LineNumber 返回當前已讀取的行號
// 可用於錯誤報告和進度追蹤
func (r *Reader) LineNumber() int {
	return r.lineNum
}

// Reset 重置讀取器到檔案開頭
// 用於需要多次讀取同一檔案的場景
func (r *Reader) Reset() error {
	if r.file == nil {
		return nil
	}

	_, err := r.file.Seek(0, 0)
	if err != nil {
		return err
	}

	r.scanner = bufio.NewScanner(r.file)
	buf := make([]byte, 0, 16*1024)
	r.scanner.Buffer(buf, 1024*1024)
	r.lineNum = 0
	r.err = nil

	return nil
}
