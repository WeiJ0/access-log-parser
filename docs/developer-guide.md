# 開發者指南

> Apache Access Log Analyzer 開發文檔

## 目錄

1. [開發環境設置](#開發環境設置)
2. [專案架構](#專案架構)
3. [API 參考](#api-參考)
4. [開發工作流程](#開發工作流程)
5. [測試指南](#測試指南)
6. [建置與部署](#建置與部署)
7. [貢獻指南](#貢獻指南)

## 開發環境設置

### 前置需求

#### 必要工具

| 工具 | 版本要求 | 用途 |
|------|---------|------|
| Go | 1.24+ | 後端開發 |
| Node.js | 18+ | 前端開發 |
| npm | 9+ | 套件管理 |
| Git | 2.0+ | 版本控制 |

#### 安裝 Wails CLI

```bash
# 安裝 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 驗證安裝
wails version

# 檢查系統依賴
wails doctor
```

#### Windows 開發環境

1. **安裝 Go**: https://golang.org/dl/
2. **安裝 Node.js**: https://nodejs.org/
3. **安裝 WebView2**: https://developer.microsoft.com/zh-tw/microsoft-edge/webview2/
4. **安裝 Git**: https://git-scm.com/

**可選工具**:
- **Visual Studio Code**: 推薦的 IDE
- **Go 擴充套件**: 提供語法高亮和自動完成
- **NSIS**: 用於建立安裝程式

#### macOS 開發環境

```bash
# 安裝 Homebrew（如果還沒安裝）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安裝 Go
brew install go

# 安裝 Node.js
brew install node

# 安裝 Git
brew install git
```

**必要的 macOS 工具**:
- **Xcode Command Line Tools**: `xcode-select --install`

#### Linux 開發環境

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang nodejs npm git build-essential libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install golang nodejs npm git gtk3-devel webkit2gtk3-devel

# Arch Linux
sudo pacman -S go nodejs npm git gtk3 webkit2gtk
```

### 取得原始碼

```bash
# 克隆專案
git clone https://github.com/yourusername/access-log-analyzer.git
cd access-log-analyzer

# 初始化 Git submodules（如果有）
git submodule update --init --recursive
```

### 安裝依賴

```bash
# 安裝 Go 依賴
go mod download

# 安裝前端依賴
cd frontend
npm install
cd ..
```

### 驗證環境

```bash
# 執行開發模式
wails dev

# 如果成功，應該看到應用程式視窗開啟
```

## 專案架構

### 目錄結構

```
access-log-analyzer/
├── cmd/                          # 應用程式入口點
│   └── apache-log-analyzer/
│       └── main.go               # 主程式
├── internal/                     # 私有套件（不對外匯出）
│   ├── app/                      # Wails 應用程式邏輯
│   │   ├── app.go                # App 結構和生命週期
│   │   ├── handlers.go           # Wails API 處理器
│   │   ├── handlers_test.go      # 整合測試
│   │   └── state.go              # 應用程式狀態
│   ├── parser/                   # 日誌解析器
│   │   ├── formats.go            # 日誌格式定義
│   │   ├── parser.go             # 解析器實作
│   │   ├── parser_test.go        # 單元測試
│   │   └── parser_bench_test.go  # 基準測試
│   ├── stats/                    # 統計分析引擎
│   │   ├── statistics.go         # 統計計算
│   │   ├── topn.go               # Top-N 堆積
│   │   ├── bot_detector.go       # 機器人偵測
│   │   ├── statistics_test.go    # 單元測試
│   │   ├── topn_test.go          # Top-N 測試
│   │   └── bot_detector_test.go  # 偵測器測試
│   ├── exporter/                 # Excel 匯出器
│   │   ├── xlsx.go               # XLSX 生成
│   │   ├── formatter.go          # 資料格式化
│   │   ├── xlsx_test.go          # 單元測試
│   │   └── xlsx_bench_test.go    # 基準測試
│   ├── models/                   # 資料模型
│   │   ├── log_entry.go          # 日誌條目
│   │   ├── statistics.go         # 統計結構
│   │   ├── errors.go             # 自訂錯誤
│   │   └── metrics.go            # 效能指標
│   ├── monitoring/               # 效能監控
│   │   └── metrics.go            # 指標收集器
│   └── testing/                  # 測試工具
│       ├── helpers.go            # 測試輔助函式
│       └── benchmark.go          # 基準測試框架
├── pkg/                          # 公開套件（可對外匯出）
│   ├── logger/                   # 日誌記錄器
│   │   └── logger.go             # zerolog 配置
│   └── apachelog/                # Apache 日誌讀取器
│       ├── reader.go             # 串流讀取器
│       └── reader_test.go        # 單元測試
├── frontend/                     # React 前端
│   ├── src/
│   │   ├── components/           # UI 組件
│   │   │   ├── LogTable.tsx
│   │   │   ├── Dashboard.tsx
│   │   │   ├── TabPanel.tsx
│   │   │   ├── SearchBar.tsx
│   │   │   ├── FilterPanel.tsx
│   │   │   ├── ErrorSummary.tsx
│   │   │   ├── ProgressIndicator.tsx
│   │   │   ├── ExportProgress.tsx
│   │   │   ├── TopIPsList.tsx
│   │   │   ├── TopPathsList.tsx
│   │   │   ├── StatusCodeDistribution.tsx
│   │   │   └── BotDetection.tsx
│   │   ├── services/             # API 服務層
│   │   │   ├── logService.ts
│   │   │   ├── searchService.ts
│   │   │   ├── searchService.test.ts
│   │   │   ├── filterService.ts
│   │   │   └── filterService.test.ts
│   │   ├── types/                # TypeScript 型別
│   │   │   └── log.ts
│   │   ├── wailsjs/              # Wails 自動生成
│   │   │   └── go/
│   │   ├── App.tsx               # 主應用程式
│   │   ├── theme.ts              # MUI 主題
│   │   └── main.tsx              # 入口點
│   ├── package.json              # npm 配置
│   ├── tsconfig.json             # TypeScript 配置
│   └── vite.config.ts            # Vite 配置
├── docs/                         # 文檔
│   ├── architecture.md           # 架構文檔
│   ├── user-guide.md             # 使用手冊
│   ├── developer-guide.md        # 本文檔
│   └── export-guide.md           # 匯出指南
├── scripts/                      # 建置和工具腳本
│   ├── generate_test_log.go      # 測試資料生成器
│   ├── build.ps1                 # Windows 建置腳本
│   └── benchmark_parser.go       # 效能基準測試
├── configs/                      # 配置檔案
├── testdata/                     # 測試資料
│   ├── valid.log
│   ├── invalid.log
│   └── 100mb.log
├── go.mod                        # Go 模組定義
├── go.sum                        # Go 依賴檢查碼
├── wails.json                    # Wails 配置
├── Makefile                      # Make 建置腳本
└── README.md                     # 專案說明
```

### 設計原則

#### 1. 模組化設計

每個模組職責單一，獨立可測試：
- **parser**: 只負責解析日誌
- **stats**: 只負責統計計算
- **exporter**: 只負責匯出
- **app**: 只負責協調和 Wails 整合

#### 2. 依賴方向

```
App → Parser → Models
App → Stats → Models
App → Exporter → Models
Stats → Parser (僅使用 LogEntry)
```

**原則**: 高層模組依賴低層模組，不產生循環依賴

#### 3. 錯誤處理

```go
// 使用自訂錯誤類型
type ParseError struct {
    Line   int
    Raw    string
    Reason string
}

// 實作 error 介面
func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at line %d: %s", e.Line, e.Reason)
}

// 返回錯誤而非 panic
func ParseFile(path string) (*ParseResult, error) {
    if err := validatePath(path); err != nil {
        return nil, fmt.Errorf("invalid path: %w", err)
    }
    // ...
}
```

#### 4. 並發安全

```go
// 使用 sync.Mutex 保護共享狀態
type State struct {
    mu    sync.RWMutex
    files map[string]*ParseResult
}

func (s *State) AddFile(name string, result *ParseResult) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.files[name] = result
}

func (s *State) GetFile(name string) (*ParseResult, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result, ok := s.files[name]
    return result, ok
}
```

## API 參考

### Wails API 端點

Wails API 端點位於 `internal/app/handlers.go`，透過 Wails 橋接自動暴露給前端。

#### SelectFile

開啟檔案選擇對話框。

```go
func (a *App) SelectFile() (string, error)
```

**返回**:
- `string`: 選擇的檔案路徑
- `error`: 錯誤訊息（使用者取消或其他錯誤）

**前端呼叫**:
```typescript
import { SelectFile } from '../wailsjs/go/app/App';

const filepath = await SelectFile();
```

**錯誤處理**:
- 使用者取消: 返回空字串和 nil 錯誤
- 系統錯誤: 返回空字串和錯誤訊息

#### ParseFile

解析 Apache 日誌檔案。

```go
func (a *App) ParseFile(filepath string) (*ParseResult, error)
```

**參數**:
- `filepath`: 日誌檔案的絕對路徑

**返回**:
- `*ParseResult`: 解析結果（包含日誌條目和統計資料）
- `error`: 解析錯誤

**ParseResult 結構**:
```go
type ParseResult struct {
    Entries    []models.LogEntry
    Statistics *models.Statistics
    Errors     []models.ParseError
    Metrics    models.PerformanceMetrics
}
```

**前端呼叫**:
```typescript
import { ParseFile } from '../wailsjs/go/app/App';

try {
    const result = await ParseFile(filepath);
    console.log(`Parsed ${result.Entries.length} entries`);
} catch (error) {
    console.error('Parse failed:', error);
}
```

**錯誤處理**:
- 檔案不存在: `ErrFileNotFound`
- 檔案太大 (>10GB): `ErrFileTooLarge`
- 格式不支援: `ErrUnsupportedFormat`
- 讀取錯誤: 系統錯誤訊息

**效能考量**:
- 大檔案解析可能需要數秒到數分鐘
- 前端應顯示進度指示器
- 使用 goroutine 避免阻塞 UI

#### ValidateLogFormat

快速驗證日誌格式（只檢查前 100 行）。

```go
func (a *App) ValidateLogFormat(filepath string) error
```

**參數**:
- `filepath`: 日誌檔案路徑

**返回**:
- `error`: 格式錯誤或 nil（格式正確）

**前端呼叫**:
```typescript
import { ValidateLogFormat } from '../wailsjs/go/app/App';

try {
    await ValidateLogFormat(filepath);
    console.log('Format is valid');
} catch (error) {
    alert('Invalid log format');
}
```

**用途**: 在解析大檔案前快速驗證格式，避免浪費時間

#### SelectSaveLocation

開啟儲存對話框。

```go
func (a *App) SelectSaveLocation() (string, error)
```

**返回**:
- `string`: 使用者選擇的儲存路徑
- `error`: 錯誤訊息

**前端呼叫**:
```typescript
import { SelectSaveLocation } from '../wailsjs/go/app/App';

const savePath = await SelectSaveLocation();
if (savePath) {
    // 使用者已選擇路徑
}
```

#### ExportToExcel

匯出分析結果為 XLSX 檔案。

```go
func (a *App) ExportToExcel(result *ParseResult, filepath string) error
```

**參數**:
- `result`: 要匯出的解析結果
- `filepath`: 儲存路徑

**返回**:
- `error`: 匯出錯誤

**前端呼叫**:
```typescript
import { ExportToExcel } from '../wailsjs/go/app/App';

try {
    await ExportToExcel(parseResult, savePath);
    alert('Export successful!');
} catch (error) {
    alert('Export failed: ' + error);
}
```

**錯誤處理**:
- 磁碟空間不足: `ErrDiskFull`
- 無寫入權限: `ErrPermissionDenied`
- 檔案已存在: 自動覆蓋（前端應確認）

**效能考量**:
- 1M 記錄可能需要 30 秒
- 前端應顯示進度對話框
- 使用串流寫入減少記憶體使用

### Go 內部 API

#### Parser API

**位置**: `internal/parser/parser.go`

```go
// NewParser 建立新的解析器實例
func NewParser(format LogFormat) *Parser

// ParseFile 解析檔案並返回結果
func (p *Parser) ParseFile(filepath string) (*ParseResult, error)

// ParseLine 解析單行日誌
func (p *Parser) ParseLine(line string) (*models.LogEntry, error)

// SetWorkers 設定並發 worker 數量（預設 4）
func (p *Parser) SetWorkers(n int)
```

**使用範例**:
```go
package main

import (
    "fmt"
    "access-log-analyzer/internal/parser"
)

func main() {
    p := parser.NewParser(parser.CombinedFormat)
    p.SetWorkers(8) // 使用 8 個並發 workers

    result, err := p.ParseFile("access.log")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Parsed %d entries\n", len(result.Entries))
}
```

#### Statistics API

**位置**: `internal/stats/statistics.go`

```go
// CalculateStatistics 計算日誌統計資訊
func CalculateStatistics(entries []models.LogEntry) *models.Statistics

// CalculateTopIPs 計算 Top N IP
func CalculateTopIPs(entries []models.LogEntry, n int) []models.IPStatistics

// CalculateTopPaths 計算 Top N 路徑
func CalculateTopPaths(entries []models.LogEntry, n int) []models.PathStatistics

// DetectBots 偵測機器人流量
func DetectBots(entries []models.LogEntry) []models.BotStatistics
```

**使用範例**:
```go
package main

import (
    "access-log-analyzer/internal/stats"
    "access-log-analyzer/internal/models"
)

func main() {
    var entries []models.LogEntry
    // ... 填充 entries

    // 計算完整統計
    statistics := stats.CalculateStatistics(entries)

    // 或只計算特定統計
    topIPs := stats.CalculateTopIPs(entries, 10)
    bots := stats.DetectBots(entries)
}
```

#### Exporter API

**位置**: `internal/exporter/xlsx.go`

```go
// ExportToXLSX 匯出為 XLSX 檔案
func ExportToXLSX(result *ParseResult, filepath string) error

// WriteLogsSheet 寫入日誌條目工作表
func WriteLogsSheet(f *excelize.File, entries []models.LogEntry) error

// WriteStatsSheet 寫入統計資料工作表
func WriteStatsSheet(f *excelize.File, stats *models.Statistics) error

// WriteBotsSheet 寫入機器人偵測工作表
func WriteBotsSheet(f *excelize.File, bots []models.BotStatistics) error
```

**使用範例**:
```go
package main

import (
    "access-log-analyzer/internal/exporter"
)

func main() {
    // 假設已有 parseResult
    err := exporter.ExportToXLSX(parseResult, "output.xlsx")
    if err != nil {
        panic(err)
    }
}
```

### 前端 API

#### LogService

**位置**: `frontend/src/services/logService.ts`

```typescript
// 開啟檔案對話框
export async function selectFile(): Promise<string>

// 解析日誌檔案
export async function parseFile(filepath: string): Promise<ParseResult>

// 驗證日誌格式
export async function validateFormat(filepath: string): Promise<boolean>

// 選擇儲存位置
export async function selectSaveLocation(): Promise<string>

// 匯出為 Excel
export async function exportToExcel(
    result: ParseResult,
    filepath: string
): Promise<void>
```

#### SearchService

**位置**: `frontend/src/services/searchService.ts`

```typescript
// 搜尋日誌條目
export function searchLogs(
    entries: LogEntry[],
    query: string
): LogEntry[]

// 搜尋特定欄位
export function searchByField(
    entries: LogEntry[],
    field: keyof LogEntry,
    query: string
): LogEntry[]
```

#### FilterService

**位置**: `frontend/src/services/filterService.ts`

```typescript
// 按狀態碼篩選
export function filterByStatusCode(
    entries: LogEntry[],
    codes: number[]
): LogEntry[]

// 按時間範圍篩選
export function filterByTimeRange(
    entries: LogEntry[],
    start: Date,
    end: Date
): LogEntry[]

// 複合篩選
export function applyFilters(
    entries: LogEntry[],
    filters: FilterOptions
): LogEntry[]
```

## 開發工作流程

### 功能開發流程

1. **建立功能分支**
```bash
git checkout -b feature/your-feature-name
```

2. **TDD 開發流程**
   - 先寫測試（確保失敗）
   - 實作功能
   - 執行測試（確保通過）
   - 重構程式碼

3. **提交程式碼**
```bash
git add .
git commit -m "feat: add your feature description"
```

4. **推送並建立 PR**
```bash
git push origin feature/your-feature-name
```

### 提交訊息格式

遵循 [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**類型 (type)**:
- `feat`: 新功能
- `fix`: 錯誤修復
- `docs`: 文檔變更
- `style`: 程式碼格式（不影響功能）
- `refactor`: 重構
- `test`: 新增測試
- `chore`: 建置流程或輔助工具變更

**範例**:
```
feat(parser): add Nginx log format support

- Implement Nginx combined format parser
- Add unit tests for Nginx format
- Update documentation

Closes #123
```

### 程式碼審查

**審查檢查清單**:
- [ ] 程式碼遵循 Go 風格指南
- [ ] 所有測試通過
- [ ] 測試覆蓋率 ≥80%
- [ ] 沒有 golint/go vet 警告
- [ ] 文檔已更新
- [ ] 效能沒有迴歸
- [ ] 沒有引入新的依賴（或已說明理由）

## 測試指南

### 單元測試

**命名慣例**: `<file>_test.go`

**測試函式命名**: `Test<FunctionName>`

**範例**:
```go
// internal/parser/parser_test.go
package parser_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "access-log-analyzer/internal/parser"
)

func TestParseLineValidCombined(t *testing.T) {
    p := parser.NewParser(parser.CombinedFormat)
    line := `127.0.0.1 - - [15/Jan/2025:14:30:25 +0800] "GET /index.html HTTP/1.1" 200 1024 "-" "Mozilla/5.0"`

    entry, err := p.ParseLine(line)

    assert.NoError(t, err)
    assert.Equal(t, "127.0.0.1", entry.IP)
    assert.Equal(t, "GET", entry.Method)
    assert.Equal(t, 200, entry.Status)
}
```

**執行單元測試**:
```bash
# 執行所有測試
go test ./...

# 執行特定套件測試
go test ./internal/parser

# 顯示詳細輸出
go test -v ./internal/parser

# 執行特定測試
go test -run TestParseLineValidCombined ./internal/parser
```

### 基準測試

**命名慣例**: `Benchmark<FunctionName>`

**範例**:
```go
// internal/parser/parser_bench_test.go
func BenchmarkParseFile100MB(b *testing.B) {
    p := parser.NewParser(parser.CombinedFormat)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := p.ParseFile("testdata/100mb.log")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**執行基準測試**:
```bash
# 執行所有基準測試
go test -bench=. ./...

# 執行特定套件基準測試
go test -bench=. -benchmem ./internal/parser

# 輸出記憶體統計
go test -bench=. -benchmem ./internal/parser
```

### 測試覆蓋率

```bash
# 生成覆蓋率報告
go test -coverprofile=coverage.out ./...

# 查看覆蓋率百分比
go tool cover -func=coverage.out

# 生成 HTML 報告
go tool cover -html=coverage.out -o coverage.html
```

**覆蓋率目標**: 80% 以上

### 整合測試

**位置**: `internal/app/handlers_test.go`

```go
func TestParseFileIntegration(t *testing.T) {
    app := NewApp()
    app.startup(context.Background())
    defer app.shutdown(context.Background())

    result, err := app.ParseFile("testdata/valid.log")

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Greater(t, len(result.Entries), 0)
}
```

### 前端測試

```bash
cd frontend

# 執行所有測試
npm test

# 執行特定測試
npm test -- searchService.test.ts

# 生成覆蓋率報告
npm test -- --coverage
```

## 建置與部署

### 開發建置

```bash
# 熱重載開發模式
wails dev

# 指定瀏覽器
wails dev -browser chrome
```

### 生產建置

```bash
# 基本建置
wails build

# 清理並重新建置
wails build -clean

# 使用 UPX 壓縮
wails build -clean -upx

# 不含 console 視窗（Windows）
wails build -clean -windowsconsole false
```

**輸出位置**: `build/bin/`

### 跨平台建置

```bash
# Windows (amd64)
wails build -platform windows/amd64

# macOS (amd64)
wails build -platform darwin/amd64

# macOS (arm64, M1/M2)
wails build -platform darwin/arm64

# Linux (amd64)
wails build -platform linux/amd64
```

**注意**: 跨平台建置需要對應平台的工具鏈

### 建立安裝程式

#### Windows (NSIS)

```bash
# 建置並建立 NSIS 安裝程式
wails build -nsis

# 自訂安裝程式圖示
wails build -nsis -windowsicon icon.ico
```

**輸出**: `build/bin/<appname>-amd64-installer.exe`

#### macOS (DMG)

```bash
# 建置並建立 DMG
wails build -platform darwin/amd64 -dmg

# 程式碼簽章（需要開發者憑證）
wails build -platform darwin/amd64 -codesign "Developer ID Application: Your Name"
```

#### Linux (AppImage)

```bash
# 建置
wails build -platform linux/amd64

# 使用 appimagetool 建立 AppImage
appimagetool build/bin/apache-log-analyzer apache-log-analyzer.AppImage
```

### 自動化建置腳本

#### Windows (PowerShell)

```powershell
# scripts/build.ps1
param(
    [string]$Version = "1.0.0",
    [switch]$Clean,
    [switch]$UPX
)

Write-Host "Building Apache Log Analyzer v$Version..."

if ($Clean) {
    Remove-Item -Recurse -Force build -ErrorAction SilentlyContinue
}

$buildArgs = @("build", "-clean")
if ($UPX) {
    $buildArgs += "-upx"
}

& wails $buildArgs

Write-Host "Build completed!"
```

**使用**:
```powershell
.\scripts\build.ps1 -Version "1.0.1" -Clean -UPX
```

#### Unix (Bash)

```bash
# scripts/build.sh
#!/bin/bash
set -e

VERSION=${1:-"1.0.0"}
CLEAN=${2:-false}
UPX=${3:-false}

echo "Building Apache Log Analyzer v$VERSION..."

if [ "$CLEAN" = "true" ]; then
    rm -rf build
fi

BUILD_ARGS="build -clean"
if [ "$UPX" = "true" ]; then
    BUILD_ARGS="$BUILD_ARGS -upx"
fi

wails $BUILD_ARGS

echo "Build completed!"
```

**使用**:
```bash
chmod +x scripts/build.sh
./scripts/build.sh 1.0.1 true true
```

### CI/CD 配置

#### GitHub Actions

```yaml
# .github/workflows/build.yml
name: Build

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'

    - name: Install Wails
      run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

    - name: Install dependencies
      run: |
        go mod download
        cd frontend && npm install

    - name: Run tests
      run: go test -v -cover ./...

    - name: Build
      run: wails build -clean

    - name: Upload artifact
      uses: actions/upload-artifact@v3
      with:
        name: ${{ matrix.os }}-build
        path: build/bin/
```

## 貢獻指南

### 報告 Bug

1. 搜尋現有的 Issues，確認問題尚未被報告
2. 建立新 Issue，包含：
   - 詳細的問題描述
   - 重現步驟
   - 預期行為 vs 實際行為
   - 系統資訊（OS, 版本等）
   - 日誌檔案（如適用）
   - 螢幕截圖（如適用）

### 提出新功能

1. 在 Issues 中提出功能請求
2. 描述使用案例和預期效益
3. 等待討論和批准
4. Fork 專案並開始開發

### Pull Request 流程

1. Fork 專案
2. 建立功能分支
3. 實作功能（遵循 TDD）
4. 確保所有測試通過
5. 更新文檔
6. 提交 PR 並描述變更
7. 等待程式碼審查
8. 根據回饋修改
9. PR 合併後刪除分支

### 程式碼風格

**Go 程式碼**:
- 使用 `gofmt` 格式化
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 所有匯出的函式和類型都應有文檔註解
- 錯誤訊息小寫開頭，不含標點符號

**TypeScript 程式碼**:
- 使用 Prettier 格式化
- 遵循 ESLint 規則
- 使用明確的型別標註
- 避免 `any` 類型

**提交訊息**:
- 遵循 Conventional Commits
- 使用現在式（"add" 而非 "added"）
- 首字母小寫
- 不超過 50 字元（標題）

### 發布流程

1. 更新版本號（`wails.json`, `package.json`）
2. 更新 CHANGELOG
3. 建立 Git tag
4. 推送 tag 觸發 CI/CD
5. 建立 GitHub Release
6. 附加建置產出和發行說明

**版本號規則** (Semantic Versioning):
- `MAJOR.MINOR.PATCH`
- `MAJOR`: 不相容的 API 變更
- `MINOR`: 新增功能（向後相容）
- `PATCH`: 錯誤修復（向後相容）

## 參考資源

### 官方文檔

- [Wails 文檔](https://wails.io/docs/introduction)
- [Go 語言規範](https://golang.org/ref/spec)
- [React 文檔](https://react.dev/)
- [TypeScript 手冊](https://www.typescriptlang.org/docs/)

### 依賴套件

- [excelize](https://xuri.me/excelize/) - Excel 處理
- [zerolog](https://github.com/rs/zerolog) - 日誌記錄
- [testify](https://github.com/stretchr/testify) - 測試斷言
- [ag-Grid](https://www.ag-grid.com/react-data-grid/) - 資料表格
- [Material-UI](https://mui.com/) - UI 組件

### 社群

- [GitHub Discussions](https://github.com/yourusername/access-log-analyzer/discussions)
- [Issues](https://github.com/yourusername/access-log-analyzer/issues)

---

**最後更新**: 2025-01-15  
**版本**: 1.0.0
