param(
  [ValidateSet("thunder", "air")]
  [string]$Tool = "thunder",
  [string]$Fixture = "benchmarks/fixture-small-api",
  [int]$Iterations = 5,
  [string]$HealthUrl = "http://127.0.0.1:18080/health",
  [int]$TimeoutSeconds = 20,
  [string]$OutJson = "",
  [string]$OutMd = ""
)

$ErrorActionPreference = "Stop"

function Wait-Health {
  param([string]$Url, [int]$TimeoutSeconds)

  $sw = [System.Diagnostics.Stopwatch]::StartNew()
  while ($sw.Elapsed.TotalSeconds -lt $TimeoutSeconds) {
    try {
      $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 2
      if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 400) {
        $sw.Stop()
        return $sw.Elapsed.TotalMilliseconds
      }
    } catch {
    }
    Start-Sleep -Milliseconds 100
  }
  throw "Health check timeout after $TimeoutSeconds seconds"
}

function Stop-ProcessTree {
  param([int]$ParentId)

  $children = Get-CimInstance Win32_Process -Filter "ParentProcessId = $ParentId" -ErrorAction SilentlyContinue
  foreach ($child in $children) {
    Stop-ProcessTree -ParentId $child.ProcessId
  }

  Stop-Process -Id $ParentId -Force -ErrorAction SilentlyContinue
}

function New-CommandSpec {
  param(
    [string]$Tool,
    [string]$RepoRoot
  )

  switch ($Tool) {
    "thunder" {
      $localBin = Join-Path $RepoRoot ".tmp-bench/thunder-bench.exe"
      $localBinDir = Split-Path -Parent $localBin
      New-Item -ItemType Directory -Path $localBinDir -Force | Out-Null
      Push-Location $RepoRoot
      try {
        go build -o $localBin ./cmd/thunder | Out-Null
      } finally {
        Pop-Location
      }
      return @{
        FilePath = $localBin
        Arguments = @("run", "./cmd/api")
      }
    }
    "air" {
      return @{
        FilePath = "air"
        Arguments = @("-c", ".air.toml")
      }
    }
    default { throw "Unsupported tool: $Tool" }
  }
}

$repoRoot = (Get-Location).Path
$fixturePath = Join-Path $repoRoot $Fixture
if (-not (Test-Path -LiteralPath $fixturePath)) {
  throw "Fixture path not found: $fixturePath"
}

$commandSpec = New-CommandSpec -Tool $Tool -RepoRoot $repoRoot
$samples = @()

Write-Host "Benchmark compare runner"
Write-Host "Tool: $Tool"
Write-Host "Fixture: $Fixture"
Write-Host "Iterations: $Iterations"

for ($i = 1; $i -le $Iterations; $i++) {
  Write-Host "Iteration $i/$Iterations"

  Push-Location $fixturePath
  try {
    if (Test-Path -LiteralPath "tmp") {
      Remove-Item -LiteralPath "tmp" -Recurse -Force -ErrorAction SilentlyContinue
    }

    $stdoutLog = Join-Path $fixturePath ("tmp/bench-{0}-stdout.log" -f $Tool)
    $stderrLog = Join-Path $fixturePath ("tmp/bench-{0}-stderr.log" -f $Tool)
    New-Item -ItemType Directory -Path (Join-Path $fixturePath "tmp") -Force | Out-Null

    $proc = Start-Process -FilePath $commandSpec.FilePath -ArgumentList $commandSpec.Arguments -WorkingDirectory $fixturePath -PassThru -WindowStyle Hidden -RedirectStandardOutput $stdoutLog -RedirectStandardError $stderrLog
    try {
      $elapsed = Wait-Health -Url $HealthUrl -TimeoutSeconds $TimeoutSeconds
      $samples += [math]::Round($elapsed, 2)
      Write-Host "  startup_ms=$([math]::Round($elapsed, 2))"
    } finally {
      Stop-ProcessTree -ParentId $proc.Id
      Start-Sleep -Milliseconds 300
    }
  } finally {
    Pop-Location
  }
}

$avg = ($samples | Measure-Object -Average).Average
$min = ($samples | Measure-Object -Minimum).Minimum
$max = ($samples | Measure-Object -Maximum).Maximum

$result = [ordered]@{
  timestamp = (Get-Date).ToString("o")
  tool = $Tool
  fixture = $Fixture
  iterations = $Iterations
  average_ms = [math]::Round($avg, 2)
  min_ms = [math]::Round($min, 2)
  max_ms = [math]::Round($max, 2)
  samples_ms = $samples
}

Write-Host "Startup benchmark result (ms):"
Write-Host "  avg=$($result.average_ms) min=$($result.min_ms) max=$($result.max_ms)"

if ($OutJson -ne "") {
  $json = $result | ConvertTo-Json -Depth 5
  Set-Content -Path $OutJson -Value $json -Encoding UTF8
  Write-Host "Wrote JSON report: $OutJson"
}

if ($OutMd -ne "") {
  $md = @(
    "# Benchmark Compare Report",
    "",
    "- Timestamp: $($result.timestamp)",
    "- Tool: $($result.tool)",
    "- Fixture: $($result.fixture)",
    "- Iterations: $($result.iterations)",
    "- Average (ms): $($result.average_ms)",
    "- Min (ms): $($result.min_ms)",
    "- Max (ms): $($result.max_ms)",
    "",
    "## Samples (ms)",
    "",
    ($samples -join ", ")
  ) -join "`n"
  Set-Content -Path $OutMd -Value $md -Encoding UTF8
  Write-Host "Wrote Markdown report: $OutMd"
}
