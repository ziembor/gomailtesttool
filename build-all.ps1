#!/usr/bin/env pwsh
# Build script for gomailtesttool
# Builds the unified gomailtest binary (optimized)

param(
    [switch]$Verbose,
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"

function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

# Header
Write-ColorOutput "`n═══════════════════════════════════════════════════════════" "Cyan"
Write-ColorOutput "  gomailtesttool Suite - Build Script" "Cyan"
Write-ColorOutput "═══════════════════════════════════════════════════════════`n" "Cyan"

# Ensure bin directory exists
$binDir = Join-Path $PSScriptRoot "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
    Write-ColorOutput "Created bin/ directory`n" "Yellow"
}

# Read version from version.go
$versionFile = Join-Path $PSScriptRoot "internal" "common" "version" "version.go"
if (-not (Test-Path $versionFile)) {
    Write-ColorOutput "ERROR: version.go not found at $versionFile" "Red"
    exit 1
}
$versionContent = Get-Content $versionFile -Raw
if ($versionContent -match 'const Version = "([^"]+)"') {
    $version = $matches[1]
} else {
    Write-ColorOutput "ERROR: Could not extract version from version.go" "Red"
    exit 1
}
Write-ColorOutput "Version: $version`n" "Yellow"

# Build gomailtest
Write-ColorOutput "Building gomailtest..." "Cyan"

$outputFile = Join-Path $binDir "gomailtest.exe"

try {
    if ($Verbose) {
        go build -v -ldflags="-s -w" -o $outputFile ./cmd/gomailtest
    } else {
        go build -ldflags="-s -w" -o $outputFile ./cmd/gomailtest
    }

    if ($LASTEXITCODE -eq 0) {
        $size = (Get-Item $outputFile).Length / 1MB
        Write-ColorOutput "  ✓ Build successful (Size: $($size.ToString('N2')) MB)" "Green"
    } else {
        throw "Build failed with exit code $LASTEXITCODE"
    }
} catch {
    Write-ColorOutput "  ✗ Build failed: $_" "Red"
    exit 1
}

# Run version check (optional)
if (-not $SkipTests) {
    Write-ColorOutput "`nVerifying version..." "Cyan"
    $toolVersion = & $outputFile --version 2>&1
    if ($toolVersion -match $version) {
        Write-ColorOutput "  ✓ Version correct: $version" "Green"
    } else {
        Write-ColorOutput "  ⚠ Version mismatch (expected: $version, got: $toolVersion)" "Yellow"
    }
}

# Summary
Write-ColorOutput "`n═══════════════════════════════════════════════════════════" "Cyan"
Write-ColorOutput "  Build Complete!" "Green"
Write-ColorOutput "═══════════════════════════════════════════════════════════" "Cyan"

Write-ColorOutput "`nBuilt executable: bin\gomailtest.exe" "White"

Write-ColorOutput "`nUsage examples:" "Yellow"
Write-ColorOutput "  .\bin\gomailtest.exe --version" "Gray"
Write-ColorOutput "  .\bin\gomailtest.exe smtp testconnect --host smtp.example.com --port 25" "Gray"
Write-ColorOutput "  .\bin\gomailtest.exe imap testconnect --host imap.gmail.com --imaps" "Gray"
Write-ColorOutput "  .\bin\gomailtest.exe pop3 testconnect --host pop.gmail.com --pop3s" "Gray"
Write-ColorOutput "  .\bin\gomailtest.exe jmap testconnect --host jmap.fastmail.com" "Gray"
Write-ColorOutput "  .\bin\gomailtest.exe msgraph getevents`n" "Gray"

Write-ColorOutput "For more information, see README.md and docs/protocols/`n" "Cyan"
