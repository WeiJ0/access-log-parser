# 系統架構文檔

## 概覽

Apache Access Log Analyzer 是一個使用 Go 語言和 Wails 框架開發的跨平台桌面應用程式。採用前後端分離架構，Go 後端負責高效能的日誌解析和資料處理，React 前端提供現代化的使用者介面。

## 架構圖

```
┌─────────────────────────────────────────────────────────────┐
│                      前端層 (React + TypeScript)              │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   LogTable   │  │  Dashboard   │  │  SearchBar   │      │
│  │  (ag-Grid)   │  │ (統計圖表)    │  │ (即時搜尋)    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ FilterPanel  │  │ ErrorSummary │  │ExportProgress│      │
│  │  (篩選器)     │  │ (錯誤摘要)    │  │ (匯出進度)    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ Wails 橋接 (Go ↔ JavaScript)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      應用層 (Wails App)                       │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐   │
│  │              App 結構 (internal/app)                   │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐     │   │
│  │  │  Handlers  │  │   State    │  │  Startup   │     │   │
│  │  │ (API 端點) │  │(狀態管理)   │  │ (生命週期)  │     │   │
│  │  └────────────┘  └────────────┘  └────────────┘     │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ 呼叫服務層
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      服務層 (Business Logic)                  │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Parser     │  │  Statistics  │  │   Exporter   │      │
│  │(日誌解析器)   │  │ (統計分析)    │  │ (Excel匯出)   │      │
│  │              │  │              │  │              │      │
│  │ • 格式識別    │  │ • Top-N 堆積 │  │ • XLSX 生成  │      │
│  │ • 串流讀取    │  │ • 機器人偵測  │  │ • 工作表格式  │      │
│  │ • 並發解析    │  │ • 狀態碼統計  │  │ • 串流寫入    │      │
│  │ • 錯誤處理    │  │ • IP 分析    │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │ 使用模型和工具
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   基礎設施層 (Infrastructure)                 │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Models    │  │   Logger     │  │  Monitoring  │      │
│  │ (資料模型)    │  │ (日誌記錄)    │  │ (效能監控)    │      │
│  │              │  │              │  │              │      │
│  │ • LogEntry   │  │ • zerolog    │  │ • 指標收集    │      │
│  │ • Statistics │  │ • 結構化日誌  │  │ • 吞吐量追蹤  │      │
│  │ • Metrics    │  │ • 多層級輸出  │  │ • 記憶體監控  │      │
│  │ • Errors     │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 核心模組

### 1. Parser Module（解析器模組）

**位置**: `internal/parser/`

**職責**: 高效能 Apache 日誌解析

**核心組件**:
- `formats.go`: 定義 Apache Combined/Common Log Format 正規表達式
- `parser.go`: 主解析邏輯，Worker Pool 並發處理
- `parser_test.go`: 單元測試
- `parser_bench_test.go`: 效能基準測試

**關鍵特性**:
- **串流讀取**: 使用 `bufio.Scanner` 逐行讀取，避免記憶體溢位
- **並發處理**: Worker Pool 模式，充分利用多核心 CPU
- **錯誤容忍**: 跳過損壞行，收集錯誤樣本（最多 100 筆）
- **效能監控**: 追蹤解析速度、記憶體峰值、處理時間

**資料流**:
```
檔案 → bufio.Scanner → Lines Channel → Worker Pool → Results Channel → []LogEntry
```

**效能指標**:
- 解析吞吐量: 134.87 MB/秒
- 並發 Workers: 4 個 goroutines
- 記憶體使用: ~3.28x 檔案大小（需保留所有記錄供 GUI 顯示）

### 2. Statistics Module（統計分析模組）

**位置**: `internal/stats/`

**職責**: 計算日誌統計資訊

**核心組件**:
- `statistics.go`: 主統計計算邏輯
- `topn.go`: Top-N Min Heap 演算法實作
- `bot_detector.go`: 機器人 User-Agent 偵測
- `statistics_test.go`: 單元測試和基準測試

**關鍵演算法**:

#### Top-N Min Heap
```go
// 時間複雜度: O(N log K)，K=10
// 空間複雜度: O(K)
type MinHeap []StatItem
// 維護固定大小的最小堆積，保留前 N 大的項目
```

#### 機器人偵測
```go
// 優先順序匹配：高信心 → 中信心 → 低信心
// 關鍵字: bot, crawler, spider, scraper, etc.
```

**統計類型**:
1. **IPStatistics**: IP 位址統計（請求數、流量、錯誤率）
2. **PathStatistics**: 請求路徑統計（請求數、平均大小、錯誤率）
3. **StatusCodeStatistics**: 狀態碼分布（2xx/3xx/4xx/5xx）
4. **BotStatistics**: 機器人流量分析（類型、信心分數）

**效能指標**:
- 10K 記錄統計: 6ms（目標 <100ms）
- 演算法複雜度: O(N log K)

### 3. Exporter Module（匯出模組）

**位置**: `internal/exporter/`

**職責**: 生成 Excel XLSX 檔案

**核心組件**:
- `xlsx.go`: 主匯出邏輯，使用 excelize
- `formatter.go`: 資料格式化器
- `xlsx_test.go`: 單元測試
- `xlsx_bench_test.go`: 效能基準測試

**工作表結構**:
1. **日誌條目** (Logs): 完整的日誌記錄表格
   - 欄位: IP, 時間戳, 方法, 路徑, 狀態碼, 大小, Referer, User-Agent
2. **統計資料** (Statistics): 彙總統計資訊
   - Top 10 IP
   - Top 10 路徑
   - 狀態碼分布
3. **機器人偵測** (Bots): 識別的機器人流量
   - IP, 類型, 信心分數, 請求數

**效能最佳化**:
- 使用 `StreamWriter` 處理大量資料
- 批次寫入減少 I/O 開銷
- Excel 行數限制處理（1,048,576 行）

### 4. App Module（應用程式模組）

**位置**: `internal/app/`

**職責**: Wails 應用程式邏輯和狀態管理

**核心組件**:
- `app.go`: Wails App 結構和生命週期
- `handlers.go`: Wails API 端點（Go → JavaScript）
- `state.go`: 應用程式狀態管理
- `handlers_test.go`: 整合測試

**Wails API 端點**:
```go
// 檔案操作
SelectFile() (string, error)                    // 開啟檔案對話框
ParseFile(filepath string) (*ParseResult, error) // 解析日誌檔案
ValidateLogFormat(filepath string) error        // 快速格式驗證

// 匯出操作
SelectSaveLocation() (string, error)            // 儲存對話框
ExportToExcel(result *ParseResult, path string) error // 匯出 Excel

// 未來擴展
GetRecentFiles() ([]string, error)              // 最近開啟的檔案
ClearRecentFiles() error                        // 清空列表
```

**狀態管理**:
- 追蹤開啟的檔案
- 管理當前活動分頁
- 快取解析結果

## 前端架構

### 技術棧

- **框架**: React 18 + TypeScript 5
- **UI 庫**: Material-UI (MUI)
- **表格**: ag-Grid Community（虛擬化支援）
- **建置工具**: Vite

### 組件結構

```
frontend/src/
├── components/
│   ├── LogTable.tsx           # 虛擬化日誌表格
│   ├── Dashboard.tsx          # 統計儀表板
│   ├── TabPanel.tsx           # 多檔案分頁管理
│   ├── SearchBar.tsx          # 搜尋框（防抖）
│   ├── FilterPanel.tsx        # 篩選器面板
│   ├── ErrorSummary.tsx       # 錯誤摘要
│   ├── ProgressIndicator.tsx  # 載入進度
│   ├── ExportProgress.tsx     # 匯出進度
│   ├── TopIPsList.tsx         # Top 10 IP 列表
│   ├── TopPathsList.tsx       # Top 10 路徑列表
│   ├── StatusCodeDistribution.tsx # 狀態碼分布
│   └── BotDetection.tsx       # 機器人偵測
├── services/
│   ├── logService.ts          # Wails API 封裝
│   ├── searchService.ts       # 搜尋邏輯
│   └── filterService.ts       # 篩選邏輯
├── types/
│   └── log.ts                 # TypeScript 型別定義
├── App.tsx                    # 主應用程式
└── theme.ts                   # Material-UI 主題配置
```

### 資料流

```
User Action → Component → Service → Wails Bridge → Go Handler
     ↓                                                    ↓
User Event                                          Business Logic
     ↓                                                    ↓
State Update ← Component ← Service ← Wails Bridge ← Result
```

### 虛擬化表格

**為什麼需要虛擬化？**
- 支援顯示百萬級記錄
- 只渲染可見行，減少 DOM 節點
- 流暢的捲動體驗（<100ms 延遲）

**ag-Grid 配置**:
```typescript
<AgGridReact
    rowData={logEntries}
    columnDefs={columnDefs}
    rowBuffer={20}              // 緩衝行數
    rowHeight={35}              // 固定行高（必要）
    virtualizationThreshold={50} // 啟用虛擬化閾值
/>
```

## 資料模型

### LogEntry（日誌條目）

```go
type LogEntry struct {
    IP        string    // IP 位址
    Timestamp time.Time // 請求時間戳
    Method    string    // HTTP 方法
    Path      string    // 請求路徑
    Protocol  string    // HTTP 協定
    Status    int       // 狀態碼
    Size      int64     // 回應大小（bytes）
    Referer   string    // 來源網址
    UserAgent string    // User Agent
}
```

### Statistics（統計資料）

```go
type Statistics struct {
    TopIPs          []IPStatistics
    TopPaths        []PathStatistics
    StatusCodeDist  StatusCodeStatistics
    Bots            []BotStatistics
    TotalRequests   int
    TotalBytes      int64
    TimeRange       TimeRange
}

type IPStatistics struct {
    IP         string
    Count      int
    TotalBytes int64
    ErrorRate  float64
}

type PathStatistics struct {
    Path        string
    Count       int
    AvgSize     float64
    ErrorRate   float64
}

type StatusCodeStatistics struct {
    Status2xx int
    Status3xx int
    Status4xx int
    Status5xx int
}

type BotStatistics struct {
    IP             string
    Type           string  // Crawler, Bot, Scraper
    ConfidenceScore float64
    RequestCount   int
}
```

## 效能優化策略

### 1. 解析器優化

- **預編譯正規表達式**: 避免重複編譯開銷
- **Worker Pool**: 並發處理多行，充分利用多核心
- **Channel Buffer**: 減少 goroutine 阻塞
- **記憶體池**: 重用 buffer，減少 GC 壓力

### 2. 統計計算優化

- **Min Heap**: O(N log K) 而非 O(N log N) 排序
- **單次掃描**: 一次遍歷計算所有統計
- **預分配記憶體**: 使用 `make(map, estimatedSize)`

### 3. 前端優化

- **虛擬化表格**: 只渲染可見行
- **防抖 (Debounce)**: 搜尋輸入延遲 300ms
- **批次更新**: 使用 `useMemo` 快取計算結果
- **Web Workers**: 未來可用於大型資料集處理

## 可觀測性

### 結構化日誌

使用 `zerolog` 記錄所有關鍵操作：

```go
log.Info().
    Str("file", filepath).
    Int("lines", count).
    Float64("duration_sec", duration.Seconds()).
    Msg("parsing completed")
```

### 效能指標

`internal/monitoring/metrics.go` 收集：
- 解析吞吐量（MB/秒）
- 記憶體使用峰值
- 處理時間
- 錯誤率

### 錯誤追蹤

- 所有錯誤包含上下文資訊
- 使用自訂錯誤類型（`ParseError`, `ValidationError`）
- 錯誤樣本收集（最多 100 筆）

## 安全性考量

### 輸入驗證

```go
// 檔案路徑驗證
func validateFilePath(path string) error {
    // 1. 路徑遍歷防護
    if strings.Contains(path, "..") {
        return ErrInvalidPath
    }
    // 2. 檔案大小限制
    if size > 10*1024*1024*1024 { // 10GB
        return ErrFileTooLarge
    }
    // 3. 檔案類型檢查
    if !isLogFile(path) {
        return ErrInvalidFileType
    }
    return nil
}
```

### 錯誤訊息處理

- 不洩漏系統路徑
- 不顯示內部錯誤細節
- 提供使用者友善的錯誤訊息

### Panic 恢復

```go
func (a *App) startup(ctx context.Context) {
    defer func() {
        if r := recover(); r != nil {
            log.Error().
                Interface("panic", r).
                Str("stack", string(debug.Stack())).
                Msg("application panic recovered")
        }
    }()
    // 初始化邏輯...
}
```

## 測試策略

### 單元測試

- 每個模組都有獨立的測試檔案
- 使用 `testify` 斷言庫
- 覆蓋率目標: 80%+

### 整合測試

- 測試 Wails API 端點
- 驗證前後端資料流
- 使用真實的測試資料

### 基準測試

```go
func BenchmarkParseFile(b *testing.B) {
    for i := 0; i < b.N; i++ {
        parser.ParseFile("testdata/100mb.log")
    }
}
```

### 效能測試

- 100MB 檔案解析時間
- 記憶體使用峰值
- 統計計算耗時
- GUI 回應速度

## 建置與部署

### 開發建置

```bash
wails dev  # 熱重載開發模式
```

### 生產建置

```bash
wails build                    # 基本建置
wails build -clean -upx        # 清理並壓縮
wails build -nsis              # 生成 Windows 安裝程式
```

### 跨平台建置

```bash
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

## 未來擴展

### 短期計畫

- [ ] 深色模式支援
- [ ] 鍵盤快捷鍵（Ctrl+O, Ctrl+S 等）
- [ ] 最近開啟的檔案列表
- [ ] 圖表視覺化（使用 Recharts）

### 長期計畫

- [ ] 支援 Nginx 日誌格式
- [ ] 自訂日誌格式解析
- [ ] 即時日誌監控（tail -f 模式）
- [ ] 資料庫匯出（SQLite, PostgreSQL）
- [ ] API 端點統計（RESTful API 分析）
- [ ] 地理位置分析（IP → 國家/城市）

## 參考資源

- [Wails 官方文檔](https://wails.io/)
- [Go 語言規範](https://golang.org/ref/spec)
- [Apache Log Format 文檔](https://httpd.apache.org/docs/current/logs.html)
- [ag-Grid 文檔](https://www.ag-grid.com/react-data-grid/)
- [excelize 文檔](https://xuri.me/excelize/)
