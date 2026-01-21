# Cloudreve Build Script (with Cloud189 Driver Support)

Write-Host "Building Cloudreve with Cloud189 Driver..." -ForegroundColor Green

# Create output directory
$outputDir = "build"
if (!(Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Create empty assets.zip if frontend resources don't exist
if (!(Test-Path "application/statics/assets.zip")) {
    Write-Host "Creating empty assets.zip..." -ForegroundColor Yellow
    if (!(Test-Path "application/statics")) {
        New-Item -ItemType Directory -Path "application/statics" | Out-Null
    }
    
    # Create a zip file with correct structure
    $tempDir = "temp_assets"
    New-Item -ItemType Directory -Path "$tempDir/assets/build" -Force | Out-Null
    
    # Create version.json
    $version = @{
        name = "cloudreve-frontend"
        version = "4.0.0"
    } | ConvertTo-Json
    
    Set-Content -Path "$tempDir/assets/build/version.json" -Value $version
    
    # Compress to zip with correct structure
    Compress-Archive -Path "$tempDir/assets" -DestinationPath "application/statics/assets.zip" -Force
    Remove-Item -Path $tempDir -Recurse -Force
    Write-Host "assets.zip created successfully" -ForegroundColor Green
}

# Build Windows version
Write-Host "`nBuilding Windows amd64..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -ldflags "-s -w -X 'github.com/cloudreve/Cloudreve/v4/application/constants.BackendVersion=4.0.0-cloud189'" -o "$outputDir/cloudreve-windows-amd64.exe" ./main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: Windows amd64 build completed" -ForegroundColor Green
    $size = (Get-Item "$outputDir/cloudreve-windows-amd64.exe").Length / 1MB
    Write-Host "  File size: $([math]::Round($size, 2)) MB" -ForegroundColor Gray
} else {
    Write-Host "Failed: Windows amd64 build failed" -ForegroundColor Red
    exit 1
}

# Build Linux version
Write-Host "`nBuilding Linux amd64..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -ldflags "-s -w -X 'github.com/cloudreve/Cloudreve/v4/application/constants.BackendVersion=4.0.0-cloud189'" -o "$outputDir/cloudreve-linux-amd64" ./main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: Linux amd64 build completed" -ForegroundColor Green
    $size = (Get-Item "$outputDir/cloudreve-linux-amd64").Length / 1MB
    Write-Host "  File size: $([math]::Round($size, 2)) MB" -ForegroundColor Gray
} else {
    Write-Host "Failed: Linux amd64 build failed" -ForegroundColor Red
    exit 1
}

Write-Host "`nBuild completed!" -ForegroundColor Green
Write-Host "Output directory: $outputDir" -ForegroundColor Yellow
Write-Host "`nSupported storage drivers:" -ForegroundColor Cyan
Write-Host "  - Local Storage (local)" -ForegroundColor Gray
Write-Host "  - Aliyun OSS (oss)" -ForegroundColor Gray
Write-Host "  - Tencent COS (cos)" -ForegroundColor Gray
Write-Host "  - Qiniu (qiniu)" -ForegroundColor Gray
Write-Host "  - Upyun (upyun)" -ForegroundColor Gray
Write-Host "  - OneDrive (onedrive)" -ForegroundColor Gray
Write-Host "  - S3 Compatible (s3)" -ForegroundColor Gray
Write-Host "  - Cloud189 (cloud189) [NEW]" -ForegroundColor Green
