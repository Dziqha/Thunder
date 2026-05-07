param(
  [int]$Iterations = 20,
  [string]$Target = "./cmd/thunder",
  [string]$OutJson = "",
  [string]$OutMd = ""
)

$ErrorActionPreference = "Stop"

Write-Host "Thunder benchmark harness"
Write-Host "Iterations: $Iterations"
Write-Host "Target: $Target"

$times = @()
for ($i = 0; $i -lt $Iterations; $i++) {
  $sw = [System.Diagnostics.Stopwatch]::StartNew()
  go build $Target | Out-Null
  $sw.Stop()
  $times += $sw.Elapsed.TotalMilliseconds
}

$avg = ($times | Measure-Object -Average).Average
$min = ($times | Measure-Object -Minimum).Minimum
$max = ($times | Measure-Object -Maximum).Maximum

Write-Host "Build benchmark result (ms):"
Write-Host "  avg=$([math]::Round($avg,2)) min=$([math]::Round($min,2)) max=$([math]::Round($max,2))"

$result = [ordered]@{
  timestamp = (Get-Date).ToString("o")
  target = $Target
  iterations = $Iterations
  average_ms = [math]::Round($avg, 2)
  min_ms = [math]::Round($min, 2)
  max_ms = [math]::Round($max, 2)
  samples_ms = $times
}

if ($OutJson -ne "") {
  $json = $result | ConvertTo-Json -Depth 5
  Set-Content -Path $OutJson -Value $json -Encoding UTF8
  Write-Host "Wrote JSON report: $OutJson"
}

if ($OutMd -ne "") {
  $md = @(
    "# Thunder Benchmark Report",
    "",
    "- Timestamp: $($result.timestamp)",
    "- Target: $($result.target)",
    "- Iterations: $($result.iterations)",
    "- Average (ms): $($result.average_ms)",
    "- Min (ms): $($result.min_ms)",
    "- Max (ms): $($result.max_ms)",
    "",
    "## Samples (ms)",
    "",
    ($times -join ", ")
  ) -join "`n"
  Set-Content -Path $OutMd -Value $md -Encoding UTF8
  Write-Host "Wrote Markdown report: $OutMd"
}
