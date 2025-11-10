# Research: Apache Access Log Analyzer

**Date**: 2025-11-06  
**Feature**: 001-apache-log-analyzer

## Research Tasks

### 1. Go GUI 框架選擇

**研究目標**: 選擇適合跨平台桌面應用程式的 Go GUI 框架，需支援虛擬化列表、分頁介面、文字顯示

**評估的框架**:

#### Option A: Fyne
- **描述**: 純 Go 實作的現代化 GUI 框架，基於 Material Design
- **優點**:
  - 完全用 Go 撰寫，無 CGO 依賴
  - 跨平台支援優秀（Windows, macOS, Linux）
  - 內建豐富的 widgets
  - 活躍的社群和文檔完善
  - 支援表格和列表 widgets
- **缺點**:
  - 大數據集的虛擬化支援有限
  - 對於複雜的資料表格，效能可能受限
- **適用性**: ⭐⭐⭐⭐ (適合，但需要自行實作虛擬化)

#### Option B: Wails
- **描述**: 使用 Web 技術（HTML/CSS/JS）作為前端，Go 作為後端的混合框架
- **優點**:
  - 可以使用成熟的前端框架（React, Vue）實作虛擬化表格
  - 靈活的 UI 設計能力
  - 良好的效能表現
  - 打包體積較小
- **缺點**:
  - 需要學習前端技術
  - 架構較複雜（前後端分離）
  - 除錯可能較困難
- **適用性**: ⭐⭐⭐⭐⭐ (非常適合，虛擬化表格有現成解決方案)

#### Option C: go-gtk (gotk3)
- **描述**: GTK+ 的 Go 綁定
- **優點**:
  - 成熟的 GUI 工具包
  - 原生外觀
  - 支援複雜的表格和樹狀視圖
- **缺點**:
  - 需要 CGO，跨平台編譯複雜
  - Windows 支援相對較弱
  - 依賴 GTK+ 執行環境
- **適用性**: ⭐⭐ (不推薦，跨平台支援問題)

#### Option D: Gio
- **描述**: 即時 GUI 框架，純 Go 實作
- **優點**:
  - 純 Go，無 CGO 依賴
  - 效能優秀，基於即時渲染
  - 適合處理大量資料
- **缺點**:
  - 學習曲線陡峭
  - 社群較小，文檔較少
  - 需要較多自訂 widgets
- **適用性**: ⭐⭐⭐ (可行，但開發成本高)

**決策**: **Wails v2**

**理由**:
1. **虛擬化表格**: 可以使用前端成熟的虛擬化表格庫（如 react-window, ag-grid），已驗證可處理百萬筆資料
2. **分頁介面**: 前端框架（React Tabs, Material-UI Tabs）提供現成的分頁組件
3. **效能**: Go 後端處理解析和統計，前端只負責顯示，職責分離清晰
4. **跨平台**: Wails 打包後可生成原生應用程式，無需使用者安裝額外環境
5. **開發效率**: 前端工具鏈成熟，UI 開發快速

**實作策略**:
- 前端: React + TypeScript + ag-grid (虛擬化表格)
- 後端: Go + Wails v2 綁定
- 通訊: Wails 提供的 Go ↔ JS 橋接

**替代方案被拒理由**:
- Fyne: 虛擬化表格需要大量自訂開發，成本高
- gotk3: 跨平台編譯複雜，Windows 支援弱
- Gio: 學習成本高，開發週期長

### 2. Apache Log 解析策略

**研究目標**: 找到高效解析 Apache access log 的最佳實踐

**決策**: **正規表達式 + 串流讀取 + Goroutine 並發**

**理由**:
1. **正規表達式**: Apache Combined/Common Log Format 有標準格式，正規表達式是最直接的解析方式
2. **串流讀取**: 使用 `bufio.Scanner` 逐行讀取，避免一次性載入整個檔案
3. **並發處理**: 使用 worker pool 模式，多個 goroutines 並發解析行
4. **記憶體管理**: 使用 channel buffer 控制記憶體使用，避免無限增長

**實作細節**:
```go
// 正規表達式範例（Combined Log Format）
pattern := `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d{3}) (\S+) "([^"]*)" "([^"]*)"`

// 串流讀取
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // 傳送到 worker channel
}

// Worker pool 並發解析
for i := 0; i < numWorkers; i++ {
    go func() {
        for line := range linesChan {
            entry := parseLogLine(line)
            resultsChan <- entry
        }
    }()
}
```

**效能優化**:
- 預編譯正規表達式 (`regexp.MustCompile`)
- 使用 `strings.Builder` 減少字串分配
- 批次處理結果，減少 channel 開銷

**替代方案**:
- 手動字串解析: 更快但程式碼複雜度高，維護困難
- 第三方解析庫: 大多不支援串流或效能不佳

### 3. 統計分析演算法

**研究目標**: 高效計算 Top-N IP、頁面、狀態碼分布等統計資訊

**決策**: **Map + Min Heap 組合**

**理由**:
1. **Map 計數**: 使用 `map[string]int` 快速統計每個項目的出現次數 (O(1) 插入/查詢)
2. **Min Heap 維護 Top-N**: 使用固定大小的最小堆維護 Top 10，空間複雜度 O(10)，時間複雜度 O(N log 10)
3. **狀態碼分類**: 使用 map 累加各類別 (2xx, 3xx, 4xx, 5xx) 的計數

**實作策略**:
```go
// 計數階段
ipCount := make(map[string]int)
for _, entry := range entries {
    ipCount[entry.IP]++
}

// Top-N 提取（使用 container/heap）
type IPStat struct {
    IP    string
    Count int
}

h := &MinHeap{}
heap.Init(h)
for ip, count := range ipCount {
    if h.Len() < 10 {
        heap.Push(h, IPStat{ip, count})
    } else if count > (*h)[0].Count {
        heap.Pop(h)
        heap.Push(h, IPStat{ip, count})
    }
}
```

**機器人偵測**:
- 使用預編譯的關鍵字列表（bot, crawler, spider, googlebot, bingbot 等）
- User Agent 字串轉小寫後進行子字串匹配
- 使用 `strings.Contains` 或 Aho-Corasick 演算法（如果關鍵字很多）

**替代方案**:
- 完整排序: O(N log N)，對於大數據集效能較差
- 使用資料庫: 增加依賴和複雜度，不適合桌面應用

### 4. XLSX 匯出技術

**研究目標**: 選擇 Go 語言的 Excel 匯出庫

**決策**: **excelize v2**

**理由**:
1. **純 Go 實作**: 無需 CGO，跨平台相容性好
2. **效能**: 支援串流寫入，適合大數據集
3. **功能完整**: 支援樣式、公式、多工作表
4. **活躍維護**: GitHub 15k+ stars，持續更新
5. **文檔完善**: 中英文文檔齊全

**實作策略**:
```go
import "github.com/xuri/excelize/v2"

f := excelize.NewFile()
// 寫入資料
for i, entry := range entries {
    f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), entry.IP)
    // ...
}
// 設定標題樣式
// 儲存檔案
f.SaveAs("output.xlsx")
```

**效能優化**:
- 使用 `StreamWriter` 模式處理超大資料集
- 分批寫入，避免記憶體峰值

**替代方案**:
- xlsx 套件: 功能較少，效能一般
- tealeg/xlsx: 較舊，維護較少

### 5. 錯誤處理與日誌記錄

**研究目標**: 實作可觀測性要求的日誌和錯誤處理策略

**決策**: **Structured Logging (zerolog) + 自訂錯誤類型**

**理由**:
1. **zerolog**: 高效能的結構化日誌庫，零分配設計
2. **結構化**: JSON 格式輸出，方便分析和除錯
3. **等級控制**: 支援 Debug, Info, Warn, Error 等級
4. **效能**: Benchmark 顯示比 logrus 快 10 倍

**錯誤處理策略**:
```go
// 自訂錯誤類型
type ParseError struct {
    Line   int
    Raw    string
    Reason string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at line %d: %s", e.Line, e.Reason)
}

// 結構化日誌
log.Error().
    Int("line", lineNum).
    Str("raw", line).
    Str("reason", "invalid format").
    Msg("failed to parse log line")
```

**錯誤摘要策略**:
- 收集所有解析錯誤到切片
- 在解析完成後生成摘要報告
- 提供前 10 個錯誤樣本供使用者檢視

### 6. 虛擬化表格實作

**研究目標**: 前端虛擬化表格的最佳實踐（配合 Wails）

**決策**: **ag-Grid Community (React)**

**理由**:
1. **虛擬化**: 內建 Row Virtualization，只渲染可見行
2. **效能**: 經過實戰驗證，可處理百萬級資料
3. **功能豐富**: 排序、篩選、搜尋都內建
4. **免費版本**: Community 版本功能已足夠

**實作策略**:
```tsx
import { AgGridReact } from 'ag-grid-react';

<AgGridReact
    rowData={logEntries}
    columnDefs={columnDefs}
    rowBuffer={20}
    rowHeight={35}
    virtualizationThreshold={50}
    onGridReady={onGridReady}
/>
```

**替代方案**:
- react-window: 需要較多自訂開發
- react-virtualized: 較舊，效能不如 ag-grid

## 技術堆疊總結

**後端 (Go)**:
- Wails v2 (GUI 框架)
- excelize v2 (Excel 匯出)
- zerolog (結構化日誌)
- testify (測試斷言)
- Go 標準庫 (regexp, bufio, container/heap)

**前端 (React + TypeScript)**:
- React 18
- TypeScript 5
- ag-Grid Community
- Material-UI (按鈕、對話框等 UI 組件)
- Recharts (如未來需要圖表可擴展)

**開發工具**:
- Go 1.21+
- Node.js 18+ (前端建置)
- Wails CLI

**效能預期**:
- 解析速度: 60-80 MB/秒（使用 4 個 worker goroutines）
- 記憶體使用: 約 1.2x 檔案大小
- GUI 響應: <100ms（得益於虛擬化）
- 匯出速度: 15,000 筆/秒

## 風險與緩解

**風險 1**: Wails 學習曲線
- **緩解**: 提供詳細的 quickstart 文檔，參考官方範例

**風險 2**: 前後端通訊效能
- **緩解**: 批次傳輸資料，避免逐筆呼叫

**風險 3**: 跨平台打包複雜度
- **緩解**: 使用 GitHub Actions 自動化建置流程

**風險 4**: 記憶體限制（處理多個 1GB 檔案）
- **緩解**: 實作記憶體壓力監控，必要時使用磁碟快取
