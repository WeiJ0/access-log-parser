# Apache Access Log Analyzer - 打包腳本
# 用途：生成安裝程式和發布套件

<#
.SYNOPSIS
    Apache Access Log Analyzer 打包腳本

.DESCRIPTION
    生成安裝程式（Windows NSIS）和發布套件（ZIP/DMG/DEB）

.PARAMETER Platform
    目標平台：windows, darwin, linux 或 all（預設：windows）

.PARAMETER Installer
    生成安裝程式（Windows NSIS、macOS DMG、Linux DEB）

.PARAMETER Portable
    生成可攜式 ZIP 版本

.PARAMETER OutputDir
    輸出目錄（預設：release）

.EXAMPLE
    .\package.ps1
    打包 Windows 版本為 ZIP

.EXAMPLE
    .\package.ps1 -Installer
    生成 Windows NSIS 安裝程式

.EXAMPLE
    .\package.ps1 -Platform all -Installer -Portable
    打包所有平台，同時生成安裝程式和可攜式版本
#>

param(
    [ValidateSet('windows', 'darwin', 'linux', 'all')]
    [string]$Platform = 'windows',
    [switch]$Installer,
    [switch]$Portable = $true,
    [string]$OutputDir = "release"
)

# 設定錯誤處理
$ErrorActionPreference = "Stop"

# 顏色輸出函式
function Write-ColorOutput {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

# 讀取版本資訊
function Get-AppVersion {
    $wailsJson = Get-Content "wails.json" -Raw | ConvertFrom-Json
    return $wailsJson.info.productVersion
}

# 檢查建置檔案
function Test-BuildExists {
    param([string]$PlatformName)
    
    $buildPaths = @{
        "windows" = "build\bin\apache-log-analyzer.exe"
        "darwin"  = "build\bin\apache-log-analyzer.app"
        "linux"   = "build\bin\apache-log-analyzer"
    }
    
    $path = $buildPaths[$PlatformName]
    if (-not (Test-Path $path)) {
        Write-ColorOutput "錯誤：找不到 $PlatformName 建置檔案" "Red"
        Write-ColorOutput "請先執行 .\build.ps1 -Platform $PlatformName" "Yellow"
        return $false
    }
    return $true
}

# 建立輸出目錄
function Initialize-OutputDirectory {
    if (-not (Test-Path $OutputDir)) {
        New-Item -ItemType Directory -Path $OutputDir | Out-Null
        Write-ColorOutput "建立輸出目錄：$OutputDir" "Gray"
    }
}

# 打包 Windows ZIP
function New-WindowsPortable {
    param([string]$Version)
    
    Write-ColorOutput "==> 打包 Windows 可攜式版本..." "Cyan"
    
    $zipName = "apache-log-analyzer-$Version-windows-x64-portable.zip"
    $zipPath = Join-Path $OutputDir $zipName
    $tempDir = "temp_package_windows"
    
    try {
        # 建立暫存目錄
        if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
        New-Item -ItemType Directory -Path $tempDir | Out-Null
        
        # 複製檔案
        Copy-Item "build\bin\apache-log-analyzer.exe" -Destination $tempDir
        Copy-Item "README.md" -Destination $tempDir
        Copy-Item "docs\user-guide.md" -Destination $tempDir
        
        # 建立 ZIP
        Compress-Archive -Path "$tempDir\*" -DestinationPath $zipPath -Force
        
        $size = [math]::Round((Get-Item $zipPath).Length / 1MB, 2)
        Write-ColorOutput "  ✓ 已建立：$zipName ($size MB)" "Green"
    }
    finally {
        if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
    }
}

# 打包 Windows 安裝程式（使用 Wails NSIS）
function New-WindowsInstaller {
    param([string]$Version)
    
    Write-ColorOutput "==> 建立 Windows 安裝程式..." "Cyan"
    
    # 檢查是否有 NSIS
    if (-not (Get-Command "makensis" -ErrorAction SilentlyContinue)) {
        Write-ColorOutput "  ⚠ 未找到 NSIS，跳過安裝程式生成" "Yellow"
        Write-ColorOutput "  提示：可從 https://nsis.sourceforge.io 安裝 NSIS" "Gray"
        return
    }
    
    # 使用 Wails 內建 NSIS 功能
    Write-ColorOutput "  執行：wails build -nsis" "Gray"
    & wails build -nsis
    
    if ($LASTEXITCODE -eq 0) {
        # 查找生成的安裝程式
        $installerPath = Get-ChildItem "build\bin" -Filter "*.exe" -Recurse | 
                         Where-Object { $_.Name -like "*installer*.exe" -or $_.Name -like "*setup*.exe" } |
                         Select-Object -First 1
        
        if ($installerPath) {
            $newName = "apache-log-analyzer-$Version-windows-x64-setup.exe"
            $destPath = Join-Path $OutputDir $newName
            Copy-Item $installerPath.FullName -Destination $destPath -Force
            
            $size = [math]::Round((Get-Item $destPath).Length / 1MB, 2)
            Write-ColorOutput "  ✓ 已建立：$newName ($size MB)" "Green"
        } else {
            Write-ColorOutput "  ⚠ 未找到生成的安裝程式" "Yellow"
        }
    } else {
        Write-ColorOutput "  ✗ 安裝程式建置失敗" "Red"
    }
}

# 打包 macOS DMG
function New-MacOSPackage {
    param([string]$Version)
    
    Write-ColorOutput "==> 打包 macOS 版本..." "Cyan"
    
    if ($Portable) {
        $zipName = "apache-log-analyzer-$Version-macos-universal.zip"
        $zipPath = Join-Path $OutputDir $zipName
        $tempDir = "temp_package_macos"
        
        try {
            if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
            New-Item -ItemType Directory -Path $tempDir | Out-Null
            
            Copy-Item "build\bin\apache-log-analyzer.app" -Destination $tempDir -Recurse
            Copy-Item "README.md" -Destination $tempDir
            Copy-Item "docs\user-guide.md" -Destination $tempDir
            
            Compress-Archive -Path "$tempDir\*" -DestinationPath $zipPath -Force
            
            $size = [math]::Round((Get-Item $zipPath).Length / 1MB, 2)
            Write-ColorOutput "  ✓ 已建立：$zipName ($size MB)" "Green"
        }
        finally {
            if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
        }
    }
    
    if ($Installer) {
        Write-ColorOutput "  ⚠ macOS DMG 生成需要 macOS 系統" "Yellow"
        Write-ColorOutput "  提示：在 macOS 上執行 'hdiutil create' 或使用 Wails 打包" "Gray"
    }
}

# 打包 Linux
function New-LinuxPackage {
    param([string]$Version)
    
    Write-ColorOutput "==> 打包 Linux 版本..." "Cyan"
    
    if ($Portable) {
        $tarName = "apache-log-analyzer-$Version-linux-x64.tar.gz"
        $tarPath = Join-Path $OutputDir $tarName
        $tempDir = "temp_package_linux"
        
        try {
            if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
            New-Item -ItemType Directory -Path $tempDir | Out-Null
            
            Copy-Item "build\bin\apache-log-analyzer" -Destination $tempDir
            Copy-Item "README.md" -Destination $tempDir
            Copy-Item "docs\user-guide.md" -Destination $tempDir
            
            # 使用 tar（需要 WSL 或 Git Bash）
            if (Get-Command "tar" -ErrorAction SilentlyContinue) {
                & tar -czf $tarPath -C $tempDir .
                $size = [math]::Round((Get-Item $tarPath).Length / 1MB, 2)
                Write-ColorOutput "  ✓ 已建立：$tarName ($size MB)" "Green"
            } else {
                # 降級為 ZIP
                $zipName = "apache-log-analyzer-$Version-linux-x64.zip"
                $zipPath = Join-Path $OutputDir $zipName
                Compress-Archive -Path "$tempDir\*" -DestinationPath $zipPath -Force
                Write-ColorOutput "  ✓ 已建立：$zipName (tar 不可用，使用 ZIP)" "Green"
            }
        }
        finally {
            if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
        }
    }
    
    if ($Installer) {
        Write-ColorOutput "  ⚠ Linux DEB/RPM 生成需要 Linux 系統" "Yellow"
        Write-ColorOutput "  提示：在 Linux 上使用 fpm 或 dpkg-deb 工具" "Gray"
    }
}

# 生成校驗和
function New-Checksums {
    Write-ColorOutput "`n==> 生成校驗和..." "Cyan"
    
    $checksumFile = Join-Path $OutputDir "checksums.txt"
    Get-ChildItem $OutputDir -File | Where-Object { $_.Extension -in @('.zip', '.tar', '.gz', '.exe', '.dmg', '.deb') } | ForEach-Object {
        $hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash
        "$hash  $($_.Name)" | Out-File -FilePath $checksumFile -Append
        Write-ColorOutput "  $($_.Name): $hash" "Gray"
    }
    
    Write-ColorOutput "✓ 校驗和已儲存到 checksums.txt" "Green"
}

# 主流程
try {
    $startTime = Get-Date
    $version = Get-AppVersion
    
    Write-ColorOutput "`n╔════════════════════════════════════════════╗" "Magenta"
    Write-ColorOutput "║  Apache Access Log Analyzer 打包工具   ║" "Magenta"
    Write-ColorOutput "╚════════════════════════════════════════════╝`n" "Magenta"
    Write-ColorOutput "版本：$version" "White"
    Write-ColorOutput "輸出目錄：$OutputDir`n" "White"
    
    Initialize-OutputDirectory
    
    # 處理平台
    $platforms = if ($Platform -eq "all") { @("windows", "darwin", "linux") } else { @($Platform) }
    
    foreach ($plat in $platforms) {
        if (-not (Test-BuildExists $plat)) {
            continue
        }
        
        switch ($plat) {
            "windows" {
                if ($Portable) { New-WindowsPortable $version }
                if ($Installer) { New-WindowsInstaller $version }
            }
            "darwin" {
                New-MacOSPackage $version
            }
            "linux" {
                New-LinuxPackage $version
            }
        }
    }
    
    # 生成校驗和
    New-Checksums
    
    # 顯示結果
    Write-ColorOutput "`n==> 打包結果" "Cyan"
    Get-ChildItem $OutputDir -File | ForEach-Object {
        $size = [math]::Round($_.Length / 1MB, 2)
        Write-ColorOutput "  $($_.Name): $size MB" "White"
    }
    
    $elapsed = (Get-Date) - $startTime
    Write-ColorOutput "`n✓ 打包完成！耗時：$($elapsed.Minutes) 分 $($elapsed.Seconds) 秒" "Green"
    
} catch {
    Write-ColorOutput "`n✗ 打包失敗：$($_.Exception.Message)" "Red"
    exit 1
}
