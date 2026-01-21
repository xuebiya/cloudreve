#!/bin/bash

# Cloud189 存储策略诊断脚本
# 用于检查天翼云盘存储策略的配置是否正确

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}==================================${NC}"
echo -e "${CYAN}天翼云盘存储策略诊断工具${NC}"
echo -e "${CYAN}==================================${NC}"
echo ""

# 检查数据库文件
DB_FILE="cloudreve.db"
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}❌ 错误: 找不到数据库文件 $DB_FILE${NC}"
    echo -e "${YELLOW}   请在Cloudreve根目录下运行此脚本${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 找到数据库文件: $DB_FILE${NC}"
echo ""

# 检查sqlite3命令
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到sqlite3命令${NC}"
    echo -e "${YELLOW}   请安装SQLite3:${NC}"
    echo -e "${YELLOW}   Ubuntu/Debian: sudo apt-get install sqlite3${NC}"
    echo -e "${YELLOW}   CentOS/RHEL: sudo yum install sqlite${NC}"
    echo -e "${YELLOW}   macOS: brew install sqlite3${NC}"
    exit 1
fi

echo -e "${GREEN}✓ SQLite3 已安装${NC}"
echo ""

# 检查Cloud189存储策略
echo -e "${CYAN}检查天翼云盘存储策略...${NC}"
POLICIES=$(sqlite3 "$DB_FILE" "SELECT id, name, type, access_key FROM storage_policies WHERE type='cloud189';")

if [ -z "$POLICIES" ]; then
    echo -e "${RED}❌ 未找到天翼云盘存储策略${NC}"
    echo -e "${YELLOW}   请先在管理后台创建天翼云盘存储策略${NC}"
    echo ""
    echo -e "${YELLOW}创建步骤：${NC}"
    echo -e "${YELLOW}1. 登录管理后台${NC}"
    echo -e "${YELLOW}2. 进入 存储策略 页面${NC}"
    echo -e "${YELLOW}3. 点击 添加存储策略${NC}"
    echo -e "${YELLOW}4. 选择 天翼云盘${NC}"
    echo -e "${YELLOW}5. 填写账号密码并创建${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 找到天翼云盘存储策略:${NC}"
CLOUD189_POLICY_ID=""
while IFS='|' read -r id name type access_key; do
    echo -e "   ID: $id"
    echo -e "   名称: $name"
    echo -e "   类型: $type"
    echo -e "   账号: $access_key"
    echo ""
    CLOUD189_POLICY_ID="$id"
done <<< "$POLICIES"

# 检查用户组配置
echo -e "${CYAN}检查用户组配置...${NC}"
GROUPS=$(sqlite3 "$DB_FILE" "SELECT g.id, g.name, g.storage_policy_groups FROM groups g;")

if [ -z "$GROUPS" ]; then
    echo -e "${RED}❌ 未找到用户组${NC}"
    exit 1
fi

FOUND_CLOUD189_GROUP=false
while IFS='|' read -r group_id group_name policy_id; do
    if [ "$policy_id" = "$CLOUD189_POLICY_ID" ]; then
        echo -e "${GREEN}✓ 用户组 '$group_name' (ID: $group_id) 使用天翼云盘策略${NC}"
        FOUND_CLOUD189_GROUP=true
    else
        echo -e "  用户组 '$group_name' (ID: $group_id) 使用策略ID: $policy_id"
    fi
done <<< "$GROUPS"

echo ""

if [ "$FOUND_CLOUD189_GROUP" = false ]; then
    echo -e "${YELLOW}⚠ 警告: 没有用户组使用天翼云盘存储策略${NC}"
    echo -e "${YELLOW}   这可能导致上传文件时出现 'Unknown policy type' 错误${NC}"
    echo ""
    echo -e "${YELLOW}解决方法：${NC}"
    echo -e "${YELLOW}1. 登录管理后台${NC}"
    echo -e "${YELLOW}2. 进入 用户组管理 页面${NC}"
    echo -e "${YELLOW}3. 编辑需要使用天翼云盘的用户组${NC}"
    echo -e "${YELLOW}4. 在 存储策略 下拉框中选择天翼云盘策略${NC}"
    echo -e "${YELLOW}5. 保存设置${NC}"
    echo ""
fi

# 检查用户
echo -e "${CYAN}检查用户配置...${NC}"
USERS=$(sqlite3 "$DB_FILE" "SELECT u.id, u.email, u.group_users FROM users u LIMIT 5;")

if [ -z "$USERS" ]; then
    echo -e "${RED}❌ 未找到用户${NC}"
else
    echo -e "${GREEN}✓ 找到用户:${NC}"
    while IFS='|' read -r user_id user_email group_id; do
        echo -e "   用户: $user_email (ID: $user_id, 用户组ID: $group_id)"
    done <<< "$USERS"
fi

echo ""
echo -e "${CYAN}==================================${NC}"
echo -e "${CYAN}诊断完成${NC}"
echo -e "${CYAN}==================================${NC}"
echo ""

if [ "$FOUND_CLOUD189_GROUP" = true ]; then
    echo -e "${GREEN}✓ 配置正确！可以开始使用天翼云盘存储${NC}"
    echo ""
    echo -e "${CYAN}测试步骤：${NC}"
    echo -e "1. 以配置了天翼云盘策略的用户组的用户身份登录"
    echo -e "2. 尝试上传一个小文件（< 10MB）"
    echo -e "3. 检查文件是否成功上传"
else
    echo -e "${YELLOW}⚠ 配置不完整，请按照上述提示完成配置${NC}"
fi

echo ""
echo -e "${CYAN}如需更多帮助，请查看:${NC}"
echo -e "- CLOUD189_USER_GUIDE_CN.md (用户指南)"
echo -e "- CLOUD189_TROUBLESHOOTING.md (故障排查)"
echo ""
