# Cloud189 存储策略诊断脚本
# 用于检查天翼云盘存储策略的配置是否正确

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "天翼云盘存储策略诊断工具" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# 检查数据库文件
$dbFile = "cloudreve.db"
if (-not (Test-Path $dbFile)) {
    Write-Host "❌ 错误: 找不到数据库文件 $dbFile" -ForegroundColor Red
    Write-Host "   请在Cloudreve根目录下运行此脚本" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ 找到数据库文件: $dbFile" -ForegroundColor Green
Write-Host ""

# 检查sqlite3命令
$sqlite3 = Get-Command sqlite3 -ErrorAction SilentlyContinue
if (-not $sqlite3) {
    Write-Host "❌ 错误: 未找到sqlite3命令" -ForegroundColor Red
    Write-Host "   请安装SQLite3: https://www.sqlite.org/download.html" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ SQLite3 已安装" -ForegroundColor Green
Write-Host ""

# 检查Cloud189存储策略
Write-Host "检查天翼云盘存储策略..." -ForegroundColor Cyan
$policies = sqlite3 $dbFile "SELECT id, name, type, access_key FROM storage_policies WHERE type='cloud189';"

if ([string]::IsNullOrEmpty($policies)) {
    Write-Host "❌ 未找到天翼云盘存储策略" -ForegroundColor Red
    Write-Host "   请先在管理后台创建天翼云盘存储策略" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "创建步骤：" -ForegroundColor Yellow
    Write-Host "1. 登录管理后台" -ForegroundColor Yellow
    Write-Host "2. 进入 存储策略 页面" -ForegroundColor Yellow
    Write-Host "3. 点击 添加存储策略" -ForegroundColor Yellow
    Write-Host "4. 选择 天翼云盘" -ForegroundColor Yellow
    Write-Host "5. 填写账号密码并创建" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ 找到天翼云盘存储策略:" -ForegroundColor Green
$policyLines = $policies -split "`n"
foreach ($line in $policyLines) {
    if (-not [string]::IsNullOrWhiteSpace($line)) {
        $parts = $line -split "\|"
        if ($parts.Length -ge 3) {
            $policyId = $parts[0]
            $policyName = $parts[1]
            $policyType = $parts[2]
            $accessKey = if ($parts.Length -ge 4) { $parts[3] } else { "" }
            
            Write-Host "   ID: $policyId" -ForegroundColor White
            Write-Host "   名称: $policyName" -ForegroundColor White
            Write-Host "   类型: $policyType" -ForegroundColor White
            Write-Host "   账号: $accessKey" -ForegroundColor White
            Write-Host ""
            
            # 保存策略ID用于后续检查
            $global:cloud189PolicyId = $policyId
        }
    }
}

# 检查用户组配置
Write-Host "检查用户组配置..." -ForegroundColor Cyan
$groups = sqlite3 $dbFile "SELECT g.id, g.name, g.storage_policy_groups FROM groups g;"

if ([string]::IsNullOrEmpty($groups)) {
    Write-Host "❌ 未找到用户组" -ForegroundColor Red
    exit 1
}

$foundCloud189Group = $false
$groupLines = $groups -split "`n"
foreach ($line in $groupLines) {
    if (-not [string]::IsNullOrWhiteSpace($line)) {
        $parts = $line -split "\|"
        if ($parts.Length -ge 3) {
            $groupId = $parts[0]
            $groupName = $parts[1]
            $policyId = $parts[2]
            
            if ($policyId -eq $global:cloud189PolicyId) {
                Write-Host "✓ 用户组 '$groupName' (ID: $groupId) 使用天翼云盘策略" -ForegroundColor Green
                $foundCloud189Group = $true
            } else {
                Write-Host "  用户组 '$groupName' (ID: $groupId) 使用策略ID: $policyId" -ForegroundColor Gray
            }
        }
    }
}

Write-Host ""

if (-not $foundCloud189Group) {
    Write-Host "⚠ 警告: 没有用户组使用天翼云盘存储策略" -ForegroundColor Yellow
    Write-Host "   这可能导致上传文件时出现 'Unknown policy type' 错误" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "解决方法：" -ForegroundColor Yellow
    Write-Host "1. 登录管理后台" -ForegroundColor Yellow
    Write-Host "2. 进入 用户组管理 页面" -ForegroundColor Yellow
    Write-Host "3. 编辑需要使用天翼云盘的用户组" -ForegroundColor Yellow
    Write-Host "4. 在 存储策略 下拉框中选择天翼云盘策略" -ForegroundColor Yellow
    Write-Host "5. 保存设置" -ForegroundColor Yellow
    Write-Host ""
}

# 检查用户
Write-Host "检查用户配置..." -ForegroundColor Cyan
$users = sqlite3 $dbFile "SELECT u.id, u.email, u.group_users FROM users u LIMIT 5;"

if ([string]::IsNullOrEmpty($users)) {
    Write-Host "❌ 未找到用户" -ForegroundColor Red
} else {
    Write-Host "✓ 找到用户:" -ForegroundColor Green
    $userLines = $users -split "`n"
    foreach ($line in $userLines) {
        if (-not [string]::IsNullOrWhiteSpace($line)) {
            $parts = $line -split "\|"
            if ($parts.Length -ge 3) {
                $userId = $parts[0]
                $userEmail = $parts[1]
                $groupId = $parts[2]
                
                Write-Host "   用户: $userEmail (ID: $userId, 用户组ID: $groupId)" -ForegroundColor White
            }
        }
    }
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "诊断完成" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

if ($foundCloud189Group) {
    Write-Host "✓ 配置正确！可以开始使用天翼云盘存储" -ForegroundColor Green
    Write-Host ""
    Write-Host "测试步骤：" -ForegroundColor Cyan
    Write-Host "1. 以配置了天翼云盘策略的用户组的用户身份登录" -ForegroundColor White
    Write-Host "2. 尝试上传一个小文件（< 10MB）" -ForegroundColor White
    Write-Host "3. 检查文件是否成功上传" -ForegroundColor White
} else {
    Write-Host "⚠ 配置不完整，请按照上述提示完成配置" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "如需更多帮助，请查看:" -ForegroundColor Cyan
Write-Host "- CLOUD189_USER_GUIDE_CN.md (用户指南)" -ForegroundColor White
Write-Host "- CLOUD189_TROUBLESHOOTING.md (故障排查)" -ForegroundColor White
Write-Host ""
