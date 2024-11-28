param (
    [Parameter(Mandatory=$true)]
    [string]$N,

    [Parameter(Mandatory=$true)]
    [string]$K
)
$exePath = Join-Path $PSScriptRoot "squares.exe"
if (-not (Test-Path $exePath)) {
    Write-Error "Cannot find squares.exe in the current directory. Make sure you've compiled the Pony program."
    exit 1
}


$arguments = "$N $K"


$result = Measure-Command {
    $process = Start-Process -FilePath $exePath -ArgumentList $arguments -PassThru -NoNewWindow -Wait
}


$userTime   = $process.UserProcessorTime
$systemTime = $process.PrivilegedProcessorTime
$realTime   = $result.TotalSeconds

Write-Host "User Time:    $($userTime.TotalSeconds) seconds"
Write-Host "System Time:  $($systemTime.TotalSeconds) seconds"
Write-Host "Real Time:    $realTime seconds"
Write-Host "`n"


$cpuTime = $userTime.TotalSeconds + $systemTime.TotalSeconds

if ($realTime -gt 0) {
    $cpuToRealRatio = $cpuTime / $realTime
    $systemToRealRatio = $systemTime.TotalSeconds / $realTime
    Write-Host "Cores used:   $cpuToRealRatio"
} else {
    Write-Host "Cores used: N/A (Real Time too small to measure)"
}
