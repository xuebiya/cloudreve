# 天翼云盘存储驱动故障排查

## 问题描述
用户报告：创建天翼云盘存储策略后，上传文件时提示"Unknown policy type"错误。

## 代码验证结果

### ✅ 已验证正确的部分

1. **类型常量定义** (`inventory/types/types.go`)
   ```go
   const PolicyTypeCloud189 = "cloud189"
   ```

2. **驱动注册** (`pkg/filemanager/manager/fs.go`)
   ```go
   case types.PolicyTypeCloud189:
       return cloud189.New(policy, m.l, m.config)
   ```

3. **前端类型定义** (`frontend/src/api/explorer.ts`)
   ```typescript
   export enum PolicyType {
       cloud189 = "cloud189",
   }
   ```

4. **前端向导组件** (`frontend/src/component/Admin/StoragePolicy/Wizards/Cloud189/Cloud189Wizard.tsx`)
   - 正确设置 `type: PolicyType.cloud189`
   - 正确映射 `access_key` (用户名) 和 `secret_key` (密码)

5. **驱动实现** (`pkg/filemanager/driver/cloud189/cloud189.go`)
   - 正确实现了所有必需的接口方法
   - 登录逻辑完整

## 可能的问题原因

### 1. 用户组未设置存储策略
当用户上传文件时，系统会从用户组获取默认存储策略。如果用户组没有设置存储策略，或者设置的不是Cloud189策略，就会出现问题。

**解决方法：**
1. 进入管理后台 → 用户组管理
2. 编辑用户所在的用户组
3. 设置"存储策略"为刚创建的天翼云盘策略

### 2. 存储策略创建失败
虽然前端显示创建成功，但可能数据库保存失败。

**验证方法：**
1. 检查数据库中的 `storage_policies` 表
2. 确认 `type` 字段的值是 `"cloud189"`（不是其他值）

**SQL查询：**
```sql
SELECT id, name, type, access_key FROM storage_policies WHERE type = 'cloud189';
```

### 3. 上传时指定了错误的策略ID
如果前端在创建上传会话时指定了错误的 `policy_id`，会导致使用错误的策略。

**验证方法：**
检查浏览器开发者工具的网络请求，查看创建上传会话的请求体：
```json
{
  "uri": "/test.txt",
  "size": 1024,
  "policy_id": "xxx"  // 这个ID应该对应Cloud189策略
}
```

## 调试步骤

### 步骤1：验证策略创建
```bash
# 查看数据库中的Cloud189策略
sqlite3 cloudreve.db "SELECT * FROM storage_policies WHERE type = 'cloud189';"
```

### 步骤2：检查用户组配置
```bash
# 查看用户组的存储策略设置
sqlite3 cloudreve.db "SELECT g.name, g.storage_policy_groups FROM groups g;"
```

### 步骤3：查看日志
启动Cloudreve时添加调试日志：
```bash
./cloudreve --debug
```

查找包含以下关键词的日志：
- "Unknown policy type"
- "PrepareUpload"
- "GetStorageDriver"

### 步骤4：测试上传流程
1. 创建一个新的Cloud189存储策略
2. 记录策略ID
3. 将用户组的存储策略设置为这个新策略
4. 尝试上传文件
5. 如果失败，查看完整的错误堆栈

## 正确的配置流程

### 1. 创建Cloud189存储策略
1. 管理后台 → 存储策略 → 添加存储策略
2. 选择"天翼云盘"
3. 填写：
   - 策略名称：例如"天翼云盘"
   - 天翼云盘账号：手机号或邮箱
   - 天翼云盘密码：账号密码
4. 点击"创建"

### 2. 配置用户组
1. 管理后台 → 用户组管理
2. 编辑目标用户组（例如"默认用户组"）
3. 在"存储策略"下拉框中选择刚创建的"天翼云盘"策略
4. 保存

### 3. 测试上传
1. 以该用户组的用户身份登录
2. 上传一个小文件测试
3. 检查文件是否成功上传到天翼云盘

## OpenList参数对照

根据OpenList的实现，天翼云盘驱动需要以下参数：

| OpenList参数 | Cloudreve映射 | 说明 |
|-------------|--------------|------|
| Username | AccessKey | 天翼云盘账号（手机号/邮箱） |
| Password | SecretKey | 天翼云盘密码 |
| Cookie | - | 可选，用于验证码处理（暂未实现） |
| RootID | - | 根目录ID，默认"-11"（已硬编码） |

## 下一步改进

如果问题仍然存在，可以考虑：

1. **添加更详细的日志**
   在 `GetStorageDriver` 方法中添加日志，记录接收到的 policy type

2. **添加策略类型验证**
   在创建存储策略时验证 type 字段是否为有效值

3. **改进错误消息**
   将"Unknown policy type"改为更详细的错误消息，包含实际接收到的类型值

4. **添加Cookie支持**
   实现Cookie参数，用于处理登录验证码的情况

## 联系支持

如果以上步骤都无法解决问题，请提供：
1. 完整的错误日志
2. 数据库中storage_policies表的相关记录
3. 浏览器开发者工具中的网络请求详情
4. Cloudreve版本信息
