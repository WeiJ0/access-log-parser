package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseLine_CombinedFormat 測試 Combined 格式單行解析
func TestParseLine_CombinedFormat(t *testing.T) {
	parser := NewParser(FormatCombined, 1)
	pattern := GetPattern(FormatCombined)

	testCases := []struct {
		name         string
		line         string
		expectError  bool
		expectIP     string
		expectMethod string
		expectURL    string
		expectStatus int
	}{
		{
			name:         "標準 Combined 格式",
			line:         `192.168.1.100 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"`,
			expectError:  false,
			expectIP:     "192.168.1.100",
			expectMethod: "GET",
			expectURL:    "/index.html",
			expectStatus: 200,
		},
		{
			name:         "POST 請求",
			line:         `10.0.0.1 - admin [06/Nov/2025:14:30:16 +0800] "POST /api/login HTTP/1.1" 200 512 "-" "curl/7.68.0"`,
			expectError:  false,
			expectIP:     "10.0.0.1",
			expectMethod: "POST",
			expectURL:    "/api/login",
			expectStatus: 200,
		},
		{
			name:         "404 錯誤",
			line:         `172.16.0.50 - - [06/Nov/2025:14:30:17 +0800] "GET /missing.html HTTP/1.1" 404 196 "https://example.com/index.html" "Mozilla/5.0"`,
			expectError:  false,
			expectIP:     "172.16.0.50",
			expectMethod: "GET",
			expectURL:    "/missing.html",
			expectStatus: 404,
		},
		{
			name:        "無效格式",
			line:        "這不是一個有效的 log 行",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := parser.parseLine(1, tc.line, pattern)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectIP, entry.IP)
			assert.Equal(t, tc.expectMethod, entry.Method)
			assert.Equal(t, tc.expectURL, entry.URL)
			assert.Equal(t, tc.expectStatus, entry.StatusCode)
		})
	}
}

// TestParseLine_CommonFormat 測試 Common 格式單行解析
func TestParseLine_CommonFormat(t *testing.T) {
	parser := NewParser(FormatCommon, 1)
	pattern := GetPattern(FormatCommon)

	line := `192.168.1.100 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234`
	entry, err := parser.parseLine(1, line, pattern)

	require.NoError(t, err)
	assert.Equal(t, "192.168.1.100", entry.IP)
	assert.Equal(t, "GET", entry.Method)
	assert.Equal(t, "/index.html", entry.URL)
	assert.Equal(t, 200, entry.StatusCode)
	assert.Equal(t, int64(1234), entry.ResponseBytes)
	assert.Empty(t, entry.Referer) // Common 格式沒有 Referer
}

// TestParseApacheTime 測試時間解析
func TestParseApacheTime(t *testing.T) {
	testCases := []struct {
		name        string
		timeStr     string
		expectError bool
		expectYear  int
		expectMonth time.Month
		expectDay   int
	}{
		{
			name:        "標準格式",
			timeStr:     "06/Nov/2025:14:30:15 +0800",
			expectError: false,
			expectYear:  2025,
			expectMonth: time.November,
			expectDay:   6,
		},
		{
			name:        "不同時區",
			timeStr:     "01/Jan/2024:00:00:00 +0000",
			expectError: false,
			expectYear:  2024,
			expectMonth: time.January,
			expectDay:   1,
		},
		{
			name:        "無效格式",
			timeStr:     "invalid-time-string",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timestamp, err := parseApacheTime(tc.timeStr)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectYear, timestamp.Year())
			assert.Equal(t, tc.expectMonth, timestamp.Month())
			assert.Equal(t, tc.expectDay, timestamp.Day())
		})
	}
}

// TestDetectFormat 測試格式自動偵測
func TestDetectFormat(t *testing.T) {
	testCases := []struct {
		name         string
		line         string
		expectFormat LogFormat
	}{
		{
			name:         "Combined 格式",
			line:         `192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"`,
			expectFormat: FormatCombined,
		},
		{
			name:         "Common 格式",
			line:         `192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234`,
			expectFormat: FormatCommon,
		},
		{
			name:         "無法辨識（預設 Combined）",
			line:         "invalid log line",
			expectFormat: FormatCombined,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			format := DetectFormat(tc.line)
			assert.Equal(t, tc.expectFormat, format)
		})
	}
}

// TestLogEntry_Methods 測試 LogEntry 的輔助方法
func TestLogEntry_Methods(t *testing.T) {
	parser := NewParser(FormatCombined, 1)
	pattern := GetPattern(FormatCombined)

	testCases := []struct {
		name           string
		line           string
		expectError    bool
		expectClient   bool
		expectServer   bool
		expectSuccess  bool
		expectCategory string
	}{
		{
			name:           "2xx 成功",
			line:           `192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /index.html HTTP/1.1" 200 1234 "-" "-"`,
			expectError:    false,
			expectClient:   false,
			expectServer:   false,
			expectSuccess:  true,
			expectCategory: "2xx",
		},
		{
			name:           "3xx 重定向",
			line:           `192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /old.html HTTP/1.1" 301 0 "-" "-"`,
			expectError:    false,
			expectClient:   false,
			expectServer:   false,
			expectSuccess:  true,
			expectCategory: "3xx",
		},
		{
			name:           "4xx 客戶端錯誤",
			line:           `192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /missing.html HTTP/1.1" 404 196 "-" "-"`,
			expectError:    true,
			expectClient:   true,
			expectServer:   false,
			expectSuccess:  false,
			expectCategory: "4xx",
		},
		{
			name:           "5xx 伺服器錯誤",
			line:           `192.168.1.1 - - [06/Nov/2025:14:30:15 +0800] "GET /error.php HTTP/1.1" 500 512 "-" "-"`,
			expectError:    true,
			expectClient:   false,
			expectServer:   true,
			expectSuccess:  false,
			expectCategory: "5xx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := parser.parseLine(1, tc.line, pattern)
			require.NoError(t, err)

			assert.Equal(t, tc.expectError, entry.IsError())
			assert.Equal(t, tc.expectClient, entry.IsClientError())
			assert.Equal(t, tc.expectServer, entry.IsServerError())
			assert.Equal(t, tc.expectSuccess, entry.IsSuccess())
			assert.Equal(t, tc.expectCategory, entry.GetStatusCategory())
		})
	}
}

// TestValidateFirstLine 測試第一行格式驗證
func TestValidateFirstLine(t *testing.T) {
	parser := NewParser(FormatCombined, 1)

	testCases := []struct {
		name        string
		filepath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效的 Apache log 檔案",
			filepath:    "../../testdata/valid.log",
			expectError: false,
		},
		{
			name:        "第一行無效格式檔案",
			filepath:    "../../testdata/invalid_first_line.log",
			expectError: true,
			errorMsg:    "不符合 Apache Access Log 格式",
		},
		{
			name:        "空檔案",
			filepath:    "../../testdata/empty.log",
			expectError: true,
			errorMsg:    "檔案為空",
		},
		{
			name:        "不存在的檔案",
			filepath:    "../../testdata/nonexistent.log",
			expectError: true,
			errorMsg:    "無法開啟檔案",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := parser.ValidateFirstLine(tc.filepath)

			if tc.expectError {
				require.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
