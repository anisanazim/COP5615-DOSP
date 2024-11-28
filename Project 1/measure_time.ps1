Param (
    [string]$ProgramPath,
    [string[]]$Arguments
)

# Start timing
$start_time = Get-Date

# Run the program and measure CPU and system time
$proc = Start-Process -FilePath $ProgramPath -ArgumentList $Arguments -PassThru -Wait

# End timing
$end_time = Get-Date

# Calculate real time
$real_time = $end_time - $start_time

# Retrieve process CPU times
$user_time = $proc.UserProcessorTime.TotalSeconds
$system_time = $proc.PrivilegedProcessorTime.TotalSeconds

# Output the results
Write-Host "Real time: $($real_time.TotalSeconds) seconds"
Write-Host "User Time: $user_time seconds"
Write-Host "System Time: $system_time seconds"

# Calculate cores used
$cores_used = ($user_time + $system_time) / $real_time.TotalSeconds
Write-Host "Cores used: $cores_used"
