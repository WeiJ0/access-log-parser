# Apache Log Analyzer - 建置腳本
# PowerShell 版本，適用於 Windows 環境

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Apache Log Analyzer - 建置腳本" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "可用命令:" -ForegroundColor Yellow
    Write-Host "  .\build.ps1 build              - 建置應用程式"
    Write-Host "  .\build.ps1 build-compressed   - 建置應用程式（啟用壓縮）"
    Write-Host "  .\build.ps1 build-installer    - 建置 Windows 安裝程式"
    Write-Host "  .\build.ps1 dev                - 啟動開發模式"
    Write-Host "  .\build.ps1 test               - 執行單元測試"
    Write-Host "  .\build.ps1 test-coverage      - 執行測試並產生覆蓋率報告"
    Write-Host "  .\build.ps1 bench              - 執行基準測試"
    Write-Host "  .\build.ps1 bench-parser       - 執行解析器基準測試"
    Write-Host "  .\build.ps1 install            - 安裝所有依賴"
    Write-Host "  .\build.ps1 lint               - 執行程式碼檢查"
    Write-Host "  .\build.ps1 fmt                - 格式化程式碼"
    Write-Host "  .\build.ps1 clean              - 清理建置產物"
    Write-Host "  .\build.ps1 generate-testdata  - 產生測試資料"
    Write-Host "  .\build.ps1 check              - 執行所有檢查"
    Write-Host "  .\build.ps1 help               - 顯示此幫助訊息"
}

function Build-App {
    Write-Host "建置應用程式..." -ForegroundColor Green
    wails build
}

function Build-Compressed {
    Write-Host "建置應用程式（啟用 UPX 壓縮）..." -ForegroundColor Green
    wails build -clean -upx
}

function Build-Installer {
    Write-Host "建置 Windows 安裝程式..." -ForegroundColor Green
    wails build -nsis
}

function Start-Dev {
    Write-Host "啟動開發模式..." -ForegroundColor Green
    wails dev
}

function Run-Tests {
    Write-Host "執行單元測試..." -ForegroundColor Green
    go test -v ./...
}

function Run-TestCoverage {
    Write-Host "執行測試並產生覆蓋率報告..." -ForegroundColor Green
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    Write-Host "覆蓋率報告已產生: coverage.html" -ForegroundColor Cyan
}

function Run-Bench {
    Write-Host "執行基準測試..." -ForegroundColor Green
    go test -bench=. -benchmem ./...
}

function Run-BenchParser {
    Write-Host "執行解析器基準測試..." -ForegroundColor Green
    go test -bench=. -benchmem ./internal/parser
}

function Install-Dependencies {
    Write-Host "安裝 Go 依賴..." -ForegroundColor Green
    go mod download
    Write-Host "安裝前端依賴..." -ForegroundColor Green
    Push-Location frontend
    npm install
    Pop-Location
}

function Run-Lint {
    Write-Host "執行程式碼檢查..." -ForegroundColor Green
    go vet ./...
    gofmt -l .
}

function Format-Code {
    Write-Host "格式化程式碼..." -ForegroundColor Green
    gofmt -w .
}

function Clean-Build {
    Write-Host "清理建置產物..." -ForegroundColor Green
    if (Test-Path "build/bin") {
        Remove-Item -Recurse -Force "build/bin"
    }
    if (Test-Path "frontend/dist") {
        Remove-Item -Recurse -Force "frontend/dist"
    }
    if (Test-Path "coverage.out") {
        Remove-Item -Force "coverage.out"
    }
    if (Test-Path "coverage.html") {
        Remove-Item -Force "coverage.html"
    }
    Write-Host "清理完成" -ForegroundColor Cyan
}

function Generate-TestData {
    Write-Host "產生測試資料..." -ForegroundColor Green
    go run scripts/generate_test_log.go -lines 100 -output testdata/valid.log -error-rate 0.05 -invalid-rate 0
    go run scripts/generate_test_log.go -lines 100 -output testdata/invalid.log -error-rate 0.05 -invalid-rate 0.2
    go run scripts/generate_test_log.go -lines 1000000 -output testdata/100mb.log -error-rate 0.05 -invalid-rate 0.01
    Write-Host "測試資料產生完成" -ForegroundColor Cyan
}

function Run-Check {
    Write-Host "執行所有檢查..." -ForegroundColor Green
    Format-Code
    Run-Lint
    Run-Tests
    Write-Host "所有檢查通過！" -ForegroundColor Cyan
}

# 主要命令分派
switch ($Command.ToLower()) {
    "build" { Build-App }
    "build-compressed" { Build-Compressed }
    "build-installer" { Build-Installer }
    "dev" { Start-Dev }
    "test" { Run-Tests }
    "test-coverage" { Run-TestCoverage }
    "bench" { Run-Bench }
    "bench-parser" { Run-BenchParser }
    "install" { Install-Dependencies }
    "lint" { Run-Lint }
    "fmt" { Format-Code }
    "clean" { Clean-Build }
    "generate-testdata" { Generate-TestData }
    "check" { Run-Check }
    "help" { Show-Help }
    default {
        Write-Host "未知命令: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
    }
}
