# GUI Startup Time Test Script
# Target: PC-004 - GUI startup time <= 2 seconds

$exePath = "..\build\bin\apache-log-analyzer.exe"

if (-not (Test-Path $exePath)) {
    Write-Host "[ERROR] Executable not found: $exePath" -ForegroundColor Red
    Write-Host "Please run 'wails build' first" -ForegroundColor Yellow
    exit 1
}

Write-Host "`n=== T142: GUI Startup Time Test (Target: <=2 sec) ===" -ForegroundColor Cyan
Write-Host ""

# Test 5 times and calculate average
$attempts = 5
$startupTimes = @()

for ($i = 1; $i -le $attempts; $i++) {
    Write-Host "Test $i/$attempts : " -NoNewline
    
    # Record start time
    $startTime = Get-Date
    
    # Launch application
    $process = Start-Process -FilePath $exePath -PassThru -WindowStyle Normal
    
    # Wait for window to appear
    $timeout = 10
    $elapsed = 0
    $windowReady = $false
    
    while ($elapsed -lt $timeout) {
        Start-Sleep -Milliseconds 100
        $elapsed += 0.1
        
        # Check if window is created
        if ($process.MainWindowHandle -ne 0) {
            $windowReady = $true
            break
        }
    }
    
    if ($windowReady) {
        $endTime = Get-Date
        $startupTime = ($endTime - $startTime).TotalSeconds
        $startupTimes += $startupTime
        
        Write-Host ("{0:N2} sec" -f $startupTime) -ForegroundColor Green
        
        # Close application immediately
        Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
        Start-Sleep -Milliseconds 500
    } else {
        Write-Host "Timeout (>10 sec)" -ForegroundColor Red
        Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
    }
}

Write-Host ""

if ($startupTimes.Count -gt 0) {
    $avgStartupTime = ($startupTimes | Measure-Object -Average).Average
    $minStartupTime = ($startupTimes | Measure-Object -Minimum).Minimum
    $maxStartupTime = ($startupTimes | Measure-Object -Maximum).Maximum
    
    Write-Host "Test Summary:" -ForegroundColor Cyan
    Write-Host "  Attempts: $($startupTimes.Count)"
    Write-Host ("  Average: {0:N2} sec" -f $avgStartupTime)
    Write-Host ("  Min: {0:N2} sec" -f $minStartupTime)
    Write-Host ("  Max: {0:N2} sec" -f $maxStartupTime)
    Write-Host ""
    
    # Check if passed
    if ($avgStartupTime -le 2.0) {
        Write-Host ("[PASS] Average startup time {0:N2} sec meets target <=2 sec" -f $avgStartupTime) -ForegroundColor Green
        exit 0
    } else {
        $diff = (($avgStartupTime / 2.0) - 1) * 100
        Write-Host ("[WARNING] Average startup time {0:N2} sec exceeds target <=2 sec (+{1:N1}%)" -f $avgStartupTime, $diff) -ForegroundColor Yellow
        exit 1
    }
} else {
    Write-Host "[FAIL] All tests failed" -ForegroundColor Red
    exit 1
}
