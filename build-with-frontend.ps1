# Cloudreve Build Script with Frontend (支持天翼云盘驱动)

Write-Host "Building Cloudreve with Cloud189 Driver and Frontend..." -ForegroundColor Green

# 检查 Node.js
Write-Host "`nChecking Node.js..." -ForegroundColor Cyan
$nodeVersion = node --version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Node.js is not installed!" -ForegroundColor Red
    Write-Host "Please install Node.js from https://nodejs.org/" -ForegroundColor Yellow
    exit 1
}
Write-Host "Node.js version: $nodeVersion" -ForegroundColor Gray

# 检查前端目录
if (!(Test-Path "../frontend-master")) {
    Write-Host "Error: Frontend directory not found!" -ForegroundColor Red
    Write-Host "Please ensure frontend-master directory exists in parent folder" -ForegroundColor Yellow
    exit 1
}

# 构建前端
Write-Host "`nBuilding frontend..." -ForegroundColor Cyan
Push-Location ../frontend-master

# 安装依赖
Write-Host "Installing frontend dependencies..." -ForegroundColor Gray
npm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to install frontend dependencies" -ForegroundColor Red
    Pop-Location
    exit 1
}

# 构建前端
Write-Host "Building frontend assets..." -ForegroundColor Gray
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to build frontend" -ForegroundColor Red
    Pop-Location
    exit 1
}

Pop-Location

# 创建 assets.zip
Write-Host "`nCreating assets.zip..." -ForegroundColor Cyan
$assetsPath = "application/statics"
if (!(Test-Path $assetsPath)) {
    New-Item -ItemType Directory -Path $assetsPath | Out-Null
}

# 压缩前端构建产物
Push-Location ../frontend-master/dist
Compress-Archive -Path * -DestinationPath "../../cloudreve-master/$assetsPath/assets.zip" -Force
Pop-Location

Write-Host "assets.zip created successfully" -ForegroundColor Green

# 验证 assets.zip
Write-Host "`nVerifying assets.zip..." -ForegroundColor Cyan
if (Test-Path "$assetsPath/assets.zip") {
    $zipSize = (Get-Item "$assetsPath/assets.zip").Length / 1MB
    Write-Host "assets.zip size: $([math]::Round($zipSize, 2)) MB" -ForegroundColor Gray
} else {
    Write-Host "Error: assets.zip not created!" -ForegroundColor Red
    exit 1
}

# 创建输出目录
$outputDir = "build"
if (!(Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# 编译 Windows 版本
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

# 编译 Linux 版本
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
