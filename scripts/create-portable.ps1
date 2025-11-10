# Apache Access Log Analyzer - Portable Version Builder
# Simple script to create portable ZIP package

param(
    [string]$Version = "1.0.0",
    [string]$OutputDir = "..\release"
)

$ErrorActionPreference = "Stop"

Write-Host "`n=== Apache Access Log Analyzer - Portable Package Builder ===" -ForegroundColor Cyan
Write-Host "Version: $Version`n" -ForegroundColor White

# Check if build exists
$exePath = "..\build\bin\apache-log-analyzer.exe"
if (-not (Test-Path $exePath)) {
    Write-Host "[ERROR] Build not found: $exePath" -ForegroundColor Red
    Write-Host "Please run 'wails build' first" -ForegroundColor Yellow
    exit 1
}

# Create output directory
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# Create temp directory
$tempDir = "temp_portable"
$portableDir = Join-Path $tempDir "apache-log-analyzer-portable"

if (Test-Path $tempDir) {
    Remove-Item $tempDir -Recurse -Force
}
New-Item -ItemType Directory -Path $portableDir | Out-Null

Write-Host "==> Copying files..." -ForegroundColor Cyan

# Copy executable
Copy-Item $exePath -Destination $portableDir
Write-Host "  + apache-log-analyzer.exe" -ForegroundColor Gray

# Copy README if exists
if (Test-Path "..\README.md") {
    Copy-Item "..\README.md" -Destination $portableDir
    Write-Host "  + README.md" -ForegroundColor Gray
}

# Create user guide
$userGuide = @"
# Apache Access Log Analyzer - Portable Edition v$Version

## Quick Start

1. Double-click apache-log-analyzer.exe to launch
2. Click "Open File" to load Apache log file
3. View analysis results and export Excel report

## System Requirements

- Windows 10 (version 1809+) or Windows 11
- WebView2 Runtime (built-in on Windows 11)
- 4GB+ RAM recommended

## WebView2 Installation

If prompted for WebView2:
https://go.microsoft.com/fwlink/?linkid=2124701

## Features

- Fast parsing: 130+ MB/sec
- Complete statistics: Top IPs, Paths, Status codes
- Bot detection
- Excel report export
- Modern Material UI interface

## Support

Issues: https://github.com/yourusername/access-log-analyzer/issues

---
Version: $Version | Build Date: $(Get-Date -Format 'yyyy-MM-dd')
"@

$userGuide | Out-File -FilePath (Join-Path $portableDir "README.txt") -Encoding UTF8
Write-Host "  + README.txt" -ForegroundColor Gray

# Create ZIP
$zipName = "apache-log-analyzer-$Version-windows-x64-portable.zip"
$zipPath = Join-Path $OutputDir $zipName

Write-Host "`n==> Creating ZIP archive..." -ForegroundColor Cyan
Compress-Archive -Path "$portableDir\*" -DestinationPath $zipPath -Force

# Calculate size and hash
$zipFile = Get-Item $zipPath
$sizeMB = [math]::Round($zipFile.Length / 1MB, 2)
$hash = (Get-FileHash $zipPath -Algorithm SHA256).Hash

# Clean up temp directory
Remove-Item $tempDir -Recurse -Force

# Display results
Write-Host "`n==> Package Created!" -ForegroundColor Green
Write-Host "  File: $zipName" -ForegroundColor White
Write-Host "  Size: $sizeMB MB" -ForegroundColor White
Write-Host "  Path: $($zipFile.FullName)" -ForegroundColor White
Write-Host "`n  SHA256: $hash" -ForegroundColor Gray

# Save checksum
$checksumFile = Join-Path $OutputDir "SHA256SUMS.txt"
"$hash  $zipName" | Out-File -FilePath $checksumFile -Encoding UTF8
Write-Host "`n  Checksum saved to: SHA256SUMS.txt" -ForegroundColor Gray

Write-Host "`n==> Distribution Ready!" -ForegroundColor Green
Write-Host "Upload to GitHub Releases or share via cloud storage" -ForegroundColor Yellow
Write-Host ""
