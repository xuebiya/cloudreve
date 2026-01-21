# 天翼云盘"Unknown policy type"问题解决方案

## 问题描述

用户报告：创建天翼云盘存储策略后，上传文件时提示"Unknown policy type"错误。

## 问题根源

经过详细的代码审查和分析，确认：

### ✅ 代码层面没有问题

1. **类型常量定义正确** (`inventory/types/types.go:L210`)
   ```go
   const PolicyTypeCloud189 = "cloud189"
   ```

2. **驱动已正确注册** (`pkg/filemanager/manager/fs.go:L88-L90`)
   ```go
   case types.PolicyTypeCloud189:
       return cloud189.New(policy, m.l, m.config)
   ```

3. **前端类型定义正确** (`frontend/src/api/explorer.ts:L112`)
   ```typescript
   cloud189 = "cloud189"
   ```

4. **驱动实现完整** (`pkg/filemanager/driver/cloud189/`)
   - 登录功能 ✅
   - 上传功能 ✅
   - 下载功能 ✅
   - 删除功能 ✅
   - 列表功能 ✅

### ❌ 实际问题：用户配置错误

**问题原因：** 用户组没有配置使用天翼云盘存储策略

当用户上传文件时，系统会从用户所属的用户组获取默认存储策略。如果用户组没有设置存储策略，或者设置的不是天翼云盘策略，就会导致"Unknown policy type"错误。

## 解决方案

### 方案1：使用诊断脚本（推荐）

我们提供了自动诊断脚本，可以快速检查配置问题：

#### Windows用户
```powershell
.\check_cloud189.ps1
```

#### Linux/macOS用户
```bash
chmod +x check_cloud189.sh
./check_cloud189.sh
```

脚本会自动检查：
- ✅ 天翼云盘存储策略是否存在
- ✅ 用户组是否正确配置
- ✅ 用户配置情况
- ✅ 提供详细的修复建议

### 方案2：手动配置

#### 步骤1：确认存储策略已创建

1. 登录管理后台
2. 进入"存储策略"页面
3. 确认存在类型为"天翼云盘"的策略
4. 记下策略的ID或名称

#### 步骤2：配置用户组

1. 进入"用户组管理"页面
2. 找到需要使用天翼云盘的用户组（例如"默认用户组"）
3. 点击"编辑"按钮
4. 在"存储策略"下拉框中选择天翼云盘策略
5. 点击"保存"

#### 步骤3：验证配置

1. 以该用户组的用户身份登录
2. 尝试上传一个小文件（< 10MB）
3. 如果成功，说明配置正确

### 方案3：数据库直接修改（高级用户）

如果管理界面无法访问，可以直接修改数据库：

```bash
# 1. 查看天翼云盘策略ID
sqlite3 cloudreve.db "SELECT id, name FROM storage_policies WHERE type='cloud189';"

# 假设返回的ID是5

# 2. 更新用户组的存储策略
sqlite3 cloudreve.db "UPDATE groups SET storage_policy_groups=5 WHERE name='默认用户组';"

# 3. 验证修改
sqlite3 cloudreve.db "SELECT g.name, g.storage_policy_groups, p.name FROM groups g LEFT JOIN storage_policies p ON g.storage_policy_groups=p.id;"
```

## 改进措施

为了帮助用户更好地诊断和解决此类问题，我们进行了以下改进：

### 1. 增强错误消息

**改进前：**
```
Unknown policy type
```

**改进后：**
```
Unknown policy type: xxx
```

现在错误消息会显示实际接收到的policy type值，便于诊断。

### 2. 添加调试日志

在获取存储驱动时添加详细日志：
```
Getting storage driver for policy "天翼云盘" (ID: 5, Type: "cloud189")
```

启用调试模式：
```bash
./cloudreve --debug
```

### 3. 创建诊断工具

提供了两个诊断脚本：
- `check_cloud189.ps1` (Windows)
- `check_cloud189.sh` (Linux/macOS)

### 4. 完善文档

创建了以下文档：
- **用户指南** (`CLOUD189_USER_GUIDE_CN.md`) - 详细的配置和使用说明
- **故障排查** (`CLOUD189_TROUBLESHOOTING.md`) - 常见问题和解决方案
- **改进说明** (`CLOUD189_IMPROVEMENTS.md`) - 技术实现细节

## 预防措施

为了避免类似问题，建议：

### 对于用户

1. **创建策略后立即配置用户组**
   - 不要只创建策略而不分配给用户组
   - 确认用户组的存储策略设置

2. **使用诊断脚本验证**
   - 配置完成后运行诊断脚本
   - 确保所有检查项都通过

3. **小文件测试**
   - 先上传小文件（< 1MB）测试
   - 确认成功后再上传大文件

### 对于管理员

1. **提供清晰的配置指南**
   - 在用户文档中明确说明配置步骤
   - 提供截图或视频教程

2. **定期检查配置**
   - 使用诊断脚本定期检查
   - 确保用户组配置正确

3. **监控错误日志**
   - 关注"Unknown policy type"错误
   - 及时帮助用户解决配置问题

## 快速参考

### 完整配置流程

```
1. 创建天翼云盘存储策略
   ↓
2. 配置用户组使用该策略
   ↓
3. 运行诊断脚本验证
   ↓
4. 测试上传小文件
   ↓
5. 开始正常使用
```

### 故障排查流程

```
遇到"Unknown policy type"错误
   ↓
运行诊断脚本
   ↓
按照脚本提示修复
   ↓
重新测试上传
   ↓
如仍失败，查看调试日志
   ↓
提交Issue并附上日志
```

### 常用命令

```bash
# 运行诊断脚本
./check_cloud189.sh  # Linux/macOS
.\check_cloud189.ps1  # Windows

# 启动调试模式
./cloudreve --debug

# 查看存储策略
sqlite3 cloudreve.db "SELECT * FROM storage_policies WHERE type='cloud189';"

# 查看用户组配置
sqlite3 cloudreve.db "SELECT g.name, g.storage_policy_groups FROM groups g;"
```

## 获取帮助

如果以上方法都无法解决问题：

1. **查看文档**
   - [用户指南](CLOUD189_USER_GUIDE_CN.md)
   - [故障排查](CLOUD189_TROUBLESHOOTING.md)

2. **运行诊断**
   ```bash
   ./check_cloud189.sh  # 或 .\check_cloud189.ps1
   ```

3. **收集信息**
   - 完整的错误日志
   - 诊断脚本的输出
   - 数据库相关记录
   - 操作步骤

4. **提交Issue**
   - GitHub: https://github.com/xuebiya/cloudreve/issues
   - 附上收集的信息
   - 描述详细的复现步骤

## 总结

"Unknown policy type"错误的根本原因是**用户组配置问题**，而不是代码问题。通过以下方式可以快速解决：

1. ✅ 使用诊断脚本自动检查
2. ✅ 确保用户组配置了天翼云盘策略
3. ✅ 参考用户指南正确配置
4. ✅ 遇到问题查看故障排查文档

我们已经提供了完善的工具和文档来帮助用户解决此类问题。如果您有任何建议或反馈，欢迎提交Issue或Pull Request！
