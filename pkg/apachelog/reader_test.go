package apachelog

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewReader 測試建立讀取器
func TestNewReader(t *testing.T) {
	// 建立測試檔案
	tempFile := createTestFile(t, []string{
		"line 1",
		"line 2",
		"line 3",
	})
	defer os.Remove(tempFile)
	
	// 測試成功建立
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	require.NotNil(t, reader)
	defer reader.Close()
	
	// 測試檔案不存在
	_, err = NewReader("/path/to/nonexistent/file.log")
	assert.Error(t, err)
}

// TestReadLine 測試逐行讀取
func TestReadLine(t *testing.T) {
	lines := []string{
		"first line",
		"second line",
		"third line",
	}
	
	tempFile := createTestFile(t, lines)
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	// 讀取所有行
	for i, expectedLine := range lines {
		lineNum, line, hasMore := reader.ReadLine()
		
		assert.True(t, hasMore, "應該有更多資料")
		assert.Equal(t, i+1, lineNum, "行號應該從 1 開始")
		assert.Equal(t, expectedLine, line)
	}
	
	// 讀取結束後應該返回 false
	_, _, hasMore := reader.ReadLine()
	assert.False(t, hasMore, "不應該有更多資料")
	
	// 驗證沒有錯誤（EOF 不是錯誤）
	assert.NoError(t, reader.Error())
}

// TestReadLine_EmptyFile 測試空檔案
func TestReadLine_EmptyFile(t *testing.T) {
	tempFile := createTestFile(t, []string{})
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	lineNum, line, hasMore := reader.ReadLine()
	
	assert.False(t, hasMore)
	assert.Equal(t, 0, lineNum)
	assert.Empty(t, line)
	assert.NoError(t, reader.Error())
}

// TestReadLine_SingleLine 測試單行檔案
func TestReadLine_SingleLine(t *testing.T) {
	tempFile := createTestFile(t, []string{"only one line"})
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	// 讀取第一行
	lineNum, line, hasMore := reader.ReadLine()
	assert.True(t, hasMore)
	assert.Equal(t, 1, lineNum)
	assert.Equal(t, "only one line", line)
	
	// 第二次讀取應該結束
	_, _, hasMore = reader.ReadLine()
	assert.False(t, hasMore)
}

// TestLineNumber 測試行號追蹤
func TestLineNumber(t *testing.T) {
	tempFile := createTestFile(t, []string{"line 1", "line 2", "line 3"})
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	assert.Equal(t, 0, reader.LineNumber(), "初始行號應該是 0")
	
	reader.ReadLine()
	assert.Equal(t, 1, reader.LineNumber())
	
	reader.ReadLine()
	assert.Equal(t, 2, reader.LineNumber())
	
	reader.ReadLine()
	assert.Equal(t, 3, reader.LineNumber())
}

// TestReset 測試重置讀取器
func TestReset(t *testing.T) {
	lines := []string{"line 1", "line 2", "line 3"}
	tempFile := createTestFile(t, lines)
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	// 讀取所有行
	for range lines {
		reader.ReadLine()
	}
	assert.Equal(t, 3, reader.LineNumber())
	
	// 重置
	err = reader.Reset()
	require.NoError(t, err)
	assert.Equal(t, 0, reader.LineNumber())
	
	// 重新讀取第一行
	lineNum, line, hasMore := reader.ReadLine()
	assert.True(t, hasMore)
	assert.Equal(t, 1, lineNum)
	assert.Equal(t, "line 1", line)
}

// TestClose 測試關閉讀取器
func TestClose(t *testing.T) {
	tempFile := createTestFile(t, []string{"test line"})
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	
	// 第一次關閉
	err = reader.Close()
	assert.NoError(t, err)
}

// TestReadLine_LongLines 測試超長行處理
func TestReadLine_LongLines(t *testing.T) {
	// 建立包含超長行的檔案（超過 16KB 緩衝區）
	longLine := string(make([]byte, 100*1024)) // 100KB 行
	for i := range longLine {
		longLine = longLine[:i] + "x" + longLine[i+1:]
	}
	
	tempFile := createTestFile(t, []string{longLine})
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	lineNum, line, hasMore := reader.ReadLine()
	
	// 應該能成功讀取（最大支援 1MB）
	assert.True(t, hasMore)
	assert.Equal(t, 1, lineNum)
	assert.Equal(t, len(longLine), len(line))
}

// TestReadLine_ManyLines 測試大量行讀取
func TestReadLine_ManyLines(t *testing.T) {
	if testing.Short() {
		t.Skip("跳過大量行測試（使用 -short 模式）")
	}
	
	// 建立包含 10,000 行的檔案
	lineCount := 10000
	lines := make([]string, lineCount)
	for i := 0; i < lineCount; i++ {
		lines[i] = "test line"
	}
	
	tempFile := createTestFile(t, lines)
	defer os.Remove(tempFile)
	
	reader, err := NewReader(tempFile)
	require.NoError(t, err)
	defer reader.Close()
	
	count := 0
	for {
		_, _, hasMore := reader.ReadLine()
		if !hasMore {
			break
		}
		count++
	}
	
	assert.Equal(t, lineCount, count)
	assert.NoError(t, reader.Error())
}

// BenchmarkReadLine 基準測試讀取效能
func BenchmarkReadLine(b *testing.B) {
	lines := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		lines[i] = `192.168.1.100 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"`
	}
	
	tempFile := createTestFile(b, lines)
	defer os.Remove(tempFile)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader, _ := NewReader(tempFile)
		
		for {
			_, _, hasMore := reader.ReadLine()
			if !hasMore {
				break
			}
		}
		
		reader.Close()
	}
}

// createTestFile 建立測試用檔案
func createTestFile(t testing.TB, lines []string) string {
	t.Helper()
	
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.log")
	
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("無法建立測試檔案: %v", err)
	}
	defer file.Close()
	
	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			t.Fatalf("無法寫入測試資料: %v", err)
		}
	}
	
	return tempFile
}
