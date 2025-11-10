# 免安裝版本發布指南

> 如何製作和分發免安裝版本 (Portable Version)

## 快速開始

### 方法 1: 使用自動化打包腳本 (推薦)

```powershell
# 1. 確保應用程式已建置
wails build

# 2. 執行打包腳本 (預設建立 ZIP 免安裝版)
cd scripts
.\create-portable.ps1

# 輸出: release/apache-log-analyzer-1.0.0-windows-x64-portable.zip
# 注意: 如果 package.ps1 出現編碼錯誤,請使用 create-portable.ps1
```

**完成!** ZIP 檔案已在 `release/` 目錄中。

---

### 方法 2: 手動建立 (完全控制)

#### Windows 免安裝版

```powershell
# 1. 建立發布目錄
New-Item -ItemType Directory -Path "release\apache-log-analyzer-portable" -Force

# 2. 複製必要檔案
Copy-Item "build\bin\apache-log-analyzer.exe" -Destination "release\apache-log-analyzer-portable\"
Copy-Item "README.md" -Destination "release\apache-log-analyzer-portable\"
Copy-Item "docs\user-guide.md" -Destination "release\apache-log-analyzer-portable\"

# 3. 建立使用說明
@"
# Apache Access Log Analyzer - 免安裝版

## 快速開始

1. 雙擊 apache-log-analyzer.exe 啟動應用程式
2. 點擊「開啟檔案」載入 Apache 日誌檔案
3. 檢視分析結果並匯出 Excel 報告

## 系統需求

- Windows 10 (版本 1809+) 或 Windows 11
- WebView2 Runtime (Windows 11 內建)
- 4GB RAM 以上建議

## 疑難排解

**問題**: 提示缺少 WebView2
**解決**: 下載並安裝 WebView2 Runtime
https://developer.microsoft.com/microsoft-edge/webview2/

## 技術支援

- 使用手冊: user-guide.md
- 問題回報: [您的 GitHub Issues URL]

"@ | Out-File -FilePath "release\apache-log-analyzer-portable\使用說明.txt" -Encoding UTF8

# 4. 壓縮為 ZIP
Compress-Archive -Path "release\apache-log-analyzer-portable\*" `
                 -DestinationPath "release\apache-log-analyzer-1.0.0-portable.zip" `
                 -Force

Write-Host "✓ 免安裝版本已建立: release\apache-log-analyzer-1.0.0-portable.zip" -ForegroundColor Green
```

---

## 檔案結構

免安裝版本應包含以下檔案:

```
apache-log-analyzer-portable/
├── apache-log-analyzer.exe    # 主程式
├── 使用說明.txt                # 快速入門指南
├── README.md                   # 專案說明
└── user-guide.md               # 完整使用手冊
```

**總大小**: 約 15-20 MB (已壓縮)

---

## 分發方式

### 選項 1: GitHub Releases (推薦)

```bash
# 1. 建立 Git tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 2. 在 GitHub 上建立 Release
# - 上傳 ZIP 檔案
# - 複製 CHANGELOG 作為發行說明
# - 勾選 "Set as the latest release"
```

### 選項 2: 雲端空間

- **Google Drive**: 分享連結設為「知道連結的使用者」
- **OneDrive**: 建立共用連結
- **Dropbox**: 建立公開連結

### 選項 3: 自架伺服器

```nginx
# Nginx 配置範例
location /downloads/ {
    alias /var/www/downloads/;
    autoindex off;
    # 設定下載標頭
    add_header Content-Disposition "attachment";
}
```

---

## 進階: 多平台打包

### Windows + macOS + Linux 一次完成

```powershell
# 建置所有平台
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64

# 打包所有平台
cd scripts
.\package.ps1 -Platform all -Portable
```

**輸出檔案**:
- `apache-log-analyzer-1.0.0-windows-x64-portable.zip` (Windows)
- `apache-log-analyzer-1.0.0-macos-universal.zip` (macOS)
- `apache-log-analyzer-1.0.0-linux-x64.tar.gz` (Linux)

---

## 最佳實務

### ✅ 應該做的事

1. **包含使用說明**: 讓使用者快速上手
2. **提供校驗和**: SHA256 確保檔案完整性
3. **明確標示版本**: 檔名包含版本號
4. **保持檔案小巧**: 避免包含不必要的檔案
5. **測試下載流程**: 確保使用者能順利下載和解壓

### ❌ 避免做的事

1. **不要包含開發檔案**: 如 `.git/`, `node_modules/`
2. **不要包含敏感資訊**: API keys, 密碼等
3. **不要包含測試資料**: 大型 `.log` 檔案
4. **不要包含編譯產物**: 如 `.o`, `.obj` 檔案

---

## 建立校驗和檔案

```powershell
# 計算 SHA256 校驗和
$hash = (Get-FileHash "release\apache-log-analyzer-1.0.0-portable.zip" -Algorithm SHA256).Hash

# 儲存到檔案
@"
# Apache Access Log Analyzer v1.0.0
# SHA256 校驗和

$hash  apache-log-analyzer-1.0.0-portable.zip

## 驗證方式

Windows PowerShell:
```powershell
Get-FileHash apache-log-analyzer-1.0.0-portable.zip -Algorithm SHA256
```

Linux/macOS:
```bash
sha256sum apache-log-analyzer-1.0.0-portable.zip
```
"@ | Out-File -FilePath "release\SHA256SUMS.txt" -Encoding UTF8

Write-Host "✓ 校驗和檔案已建立" -ForegroundColor Green
```

---

## 範例: 完整發布流程

```powershell
# 完整的發布腳本範例

# 1. 設定版本號
$version = "1.0.0"

# 2. 建置應用程式
Write-Host "==> 建置應用程式..." -ForegroundColor Cyan
wails build -clean

# 3. 建立免安裝版本
Write-Host "`n==> 建立免安裝版本..." -ForegroundColor Cyan
cd scripts
.\package.ps1 -Platform windows -Portable
cd ..

# 4. 生成校驗和
Write-Host "`n==> 生成校驗和..." -ForegroundColor Cyan
$zipFile = "release\apache-log-analyzer-$version-windows-x64-portable.zip"
$hash = (Get-FileHash $zipFile -Algorithm SHA256).Hash
"$hash  $(Split-Path $zipFile -Leaf)" | Out-File "release\SHA256SUMS.txt"

# 5. 建立發行說明
Write-Host "`n==> 建立發行說明..." -ForegroundColor Cyan
@"
# Apache Access Log Analyzer v$version

## 下載

- [Windows 免安裝版 (15MB)](apache-log-analyzer-$version-windows-x64-portable.zip)
- [SHA256 校驗和](SHA256SUMS.txt)

## 新功能

- ⚡ 超快解析速度: 130+ MB/秒
- 📊 完整統計分析: Top IPs, Paths, 狀態碼分布
- 🤖 自動機器人偵測
- 📤 Excel 報告匯出
- 🎨 現代化 Material UI 介面

## 系統需求

- Windows 10 (1809+) / Windows 11
- WebView2 Runtime
- 4GB RAM 建議

## 快速開始

1. 下載並解壓縮 ZIP 檔案
2. 雙擊 apache-log-analyzer.exe
3. 點擊「開啟檔案」載入日誌

## 變更日誌

查看完整變更: [CHANGELOG.md](CHANGELOG.md)
"@ | Out-File "release\RELEASE_NOTES.md" -Encoding UTF8

# 6. 顯示結果
Write-Host "`n==> 發布檔案清單:" -ForegroundColor Cyan
Get-ChildItem release | Format-Table Name, Length, LastWriteTime

Write-Host "`n✓ 發布準備完成!" -ForegroundColor Green
Write-Host "請至 GitHub > Releases > Draft a new release 上傳檔案" -ForegroundColor Yellow
```

---

## 使用者安裝說明 (給使用者看的)

### Windows 免安裝版安裝步驟

1. **下載檔案**
   - 下載 `apache-log-analyzer-1.0.0-windows-x64-portable.zip`

2. **解壓縮**
   - 右鍵點擊 ZIP 檔案
   - 選擇「解壓縮全部...」
   - 選擇解壓縮位置 (例如: `C:\Tools\`)

3. **首次執行**
   - 進入解壓縮的資料夾
   - 雙擊 `apache-log-analyzer.exe`
   - 如果出現 SmartScreen 警告,點擊「更多資訊」→「仍要執行」

4. **WebView2 檢查** (Windows 10 使用者)
   - 如果提示缺少 WebView2
   - 點擊下載連結或前往: https://go.microsoft.com/fwlink/?linkid=2124701
   - 安裝完成後重新啟動應用程式

5. **開始使用**
   - 點擊「開啟檔案」按鈕
   - 選擇 Apache 日誌檔案 (`.log`)
   - 等待解析完成
   - 檢視統計結果或匯出 Excel

---

## 常見問題

### Q: 為什麼要用免安裝版?

**A**: 免安裝版優點:
- ✅ 無需管理員權限
- ✅ 不修改系統設定
- ✅ 可放在 USB 隨身碟執行
- ✅ 解壓縮即可使用
- ✅ 刪除資料夾即完全移除

### Q: 檔案大小為何這麼大?

**A**: 15MB 包含:
- Wails Runtime (~8MB)
- WebView2 綁定 (~2MB)
- 應用程式邏輯 (~3MB)
- 前端資源 (~2MB)

這是正常的桌面應用程式大小,比 Electron 應用程式小很多。

### Q: 如何更新到新版本?

**A**: 
1. 下載新版本 ZIP
2. 解壓縮到新位置
3. (可選) 刪除舊版本資料夾

設定檔和最近檔案列表會保留 (儲存在 `%APPDATA%`)。

### Q: 可以放在公司內網分享嗎?

**A**: 可以!免安裝版適合:
- 內部網路檔案伺服器
- 共用資料夾
- 公司軟體庫
- USB 隨身碟分發

### Q: 支援哪些作業系統?

**A**: 
- ✅ Windows 10 (1809+)
- ✅ Windows 11
- ✅ Windows Server 2019+
- ⚠️ Windows 7/8.1: 不支援 (缺少 WebView2)

---

## 授權與分發

請確保遵守專案的授權條款。如果是開源專案,在分發時應:

1. 保留授權檔案 (`LICENSE`)
2. 註明原作者
3. 提供原始碼連結 (如果是 GPL 授權)

---

## 檢查清單

製作免安裝版本前的檢查清單:

- [ ] 應用程式已建置且功能正常
- [ ] 包含使用說明文件
- [ ] 檔案大小合理 (<50MB)
- [ ] 版本號正確標示
- [ ] 計算並提供 SHA256 校驗和
- [ ] 在乾淨環境測試執行
- [ ] 檢查是否需要 WebView2
- [ ] 準備發行說明
- [ ] 設定下載連結
- [ ] 通知使用者更新

---

**版本**: 1.0.0  
**最後更新**: 2025-11-10  
**維護者**: 專案團隊
