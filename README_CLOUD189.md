# Cloudreve with Cloud189 Driver Support

[![Build Status](https://github.com/xuebiya/cloudreve/actions/workflows/build.yml/badge.svg)](https://github.com/xuebiya/cloudreve/actions/workflows/build.yml)
[![Test Build](https://github.com/xuebiya/cloudreve/actions/workflows/test-build.yml/badge.svg)](https://github.com/xuebiya/cloudreve/actions/workflows/test-build.yml)
[![License](https://img.shields.io/github/license/xuebiya/cloudreve)](LICENSE)
[![Release](https://img.shields.io/github/v/release/xuebiya/cloudreve)](https://github.com/xuebiya/cloudreve/releases)

> This is a custom build of [Cloudreve](https://github.com/cloudreve/Cloudreve) v4 with **China Telecom Cloud189 (天翼云盘)** storage driver support.

## ✨ Features

### Cloud189 Storage Driver (NEW)

- ✅ **File Upload** - Chunked upload support (10MB chunks)
- ✅ **File Download** - Direct download with automatic redirect handling
- ✅ **File Deletion** - Batch deletion support
- ✅ **File Listing** - Fast directory listing
- ✅ **Auto Login** - Automatic session management and re-login
- ✅ **Secure** - RSA encryption for credentials, HMAC-SHA1 signatures
- ✅ **Reliable** - MD5 verification for uploaded files

### All Original Cloudreve Features

- Multiple storage backends (Local, OSS, COS, S3, OneDrive, etc.)
- User management and group permissions
- WebDAV support
- File sharing and collaboration
- Archive and extraction
- Image thumbnails
- Video streaming
- And much more...

## 📦 Download

### Pre-compiled Binaries

Download from [Releases](https://github.com/xuebiya/cloudreve/releases) page:

| Platform | Architecture | Download |
|----------|-------------|----------|
| Windows | amd64 | [cloudreve-windows-amd64.exe](https://github.com/xuebiya/cloudreve/releases/latest) |
| Windows | arm64 | [cloudreve-windows-arm64.exe](https://github.com/xuebiya/cloudreve/releases/latest) |
| Linux | amd64 | [cloudreve-linux-amd64](https://github.com/xuebiya/cloudreve/releases/latest) |
| Linux | arm64 | [cloudreve-linux-arm64](https://github.com/xuebiya/cloudreve/releases/latest) |
| Linux | arm | [cloudreve-linux-arm](https://github.com/xuebiya/cloudreve/releases/latest) |
| macOS | amd64 | [cloudreve-darwin-amd64](https://github.com/xuebiya/cloudreve/releases/latest) |
| macOS | arm64 | [cloudreve-darwin-arm64](https://github.com/xuebiya/cloudreve/releases/latest) |

### Build from Source

```bash
# Clone repository
git clone https://github.com/xuebiya/cloudreve.git
cd cloudreve

# Build
go build -ldflags "-s -w" -o cloudreve ./main.go
```

## 🚀 Quick Start

### 1. Download and Run

**Linux/macOS**:
```bash
# Download
wget https://github.com/xuebiya/cloudreve/releases/latest/download/cloudreve-linux-amd64

# Make executable
chmod +x cloudreve-linux-amd64

# Run
./cloudreve-linux-amd64
```

**Windows**:
```cmd
# Download cloudreve-windows-amd64.exe
# Double-click to run or use command line:
cloudreve-windows-amd64.exe
```

### 2. Access Web Interface

Open your browser and navigate to:
```
http://localhost:5212
```

Default admin credentials will be displayed in the console on first run.

### 3. Configure Cloud189 Storage

1. Log in to admin panel
2. Navigate to **Storage Policies**
3. Click **Add Storage Policy**
4. Fill in the form:
   - **Policy Name**: Cloud189
   - **Storage Type**: `cloud189`
   - **AccessKey**: Your Cloud189 account (phone number or email)
   - **SecretKey**: Your Cloud189 password
   - **Is Private**: Yes (recommended)
5. Click **Save**

## 📖 Documentation

### Cloud189 Driver Configuration

#### Via Admin Panel

See Quick Start section above.

#### Via Database

```sql
INSERT INTO storage_policies (
  name, 
  type, 
  access_key, 
  secret_key, 
  is_private, 
  settings, 
  created_at, 
  updated_at
)
VALUES (
  'Cloud189',
  'cloud189',
  'your_phone_or_email',
  'your_password',
  1,
  '{}',
  NOW(),
  NOW()
);
```

### Supported Storage Drivers

- ✅ Local Storage
- ✅ Aliyun OSS
- ✅ Tencent COS
- ✅ Qiniu
- ✅ Upyun
- ✅ OneDrive
- ✅ S3 Compatible
- ✅ Kingsoft KS3
- ✅ Huawei OBS
- ✅ **Cloud189 (NEW)** ⭐

### Cloud189 Driver Limitations

- ❌ Folder operations (create, move, rename) not yet supported
- ❌ Thumbnail generation not supported
- ❌ Media metadata extraction not supported
- ⚠️ Currently only supports root directory operations

## 🔧 Advanced Configuration

### Environment Variables

```bash
# Database
export DB_TYPE=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=cloudreve
export DB_PASSWORD=your_password
export DB_NAME=cloudreve

# Server
export LISTEN=:5212
export DEBUG=false

# Run
./cloudreve
```

### Configuration File

Create `conf.ini`:

```ini
[System]
Mode = master
Listen = :5212
Debug = false

[Database]
Type = mysql
Host = localhost
Port = 3306
User = cloudreve
Password = your_password
Name = cloudreve
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Git

### Build

```bash
# Clone repository
git clone https://github.com/xuebiya/cloudreve.git
cd cloudreve

# Install dependencies
go mod download

# Build
go build -o cloudreve ./main.go
```

### Run Tests

```bash
go test ./...
```

## 📝 Changelog

### v4.0.0-cloud189 (2026-01-21)

- ✨ Added Cloud189 storage driver support
- ✅ File upload with chunked upload (10MB chunks)
- ✅ File download with automatic redirect handling
- ✅ File deletion
- ✅ File listing
- ✅ Automatic login and session management
- 🔒 RSA encryption for login credentials
- 🔐 HMAC-SHA1 signatures for upload requests
- ✔️ MD5 verification for uploaded files

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project follows the Cloudreve project license.

## 🙏 Credits

- **Cloudreve**: [https://github.com/cloudreve/Cloudreve](https://github.com/cloudreve/Cloudreve)
- **OpenList**: [https://github.com/OpenListTeam/OpenList](https://github.com/OpenListTeam/OpenList) (Cloud189 driver source)

## 📮 Support

- **Issues**: [GitHub Issues](https://github.com/xuebiya/cloudreve/issues)
- **Discussions**: [GitHub Discussions](https://github.com/xuebiya/cloudreve/discussions)
- **Original Cloudreve Docs**: [https://docs.cloudreve.org/](https://docs.cloudreve.org/)

## ⚠️ Disclaimer

This is an **unofficial build** with Cloud189 driver support. For the official Cloudreve release, please visit [https://github.com/cloudreve/Cloudreve](https://github.com/cloudreve/Cloudreve).

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=xuebiya/cloudreve&type=Date)](https://star-history.com/#xuebiya/cloudreve&Date)

---

Made with ❤️ by [xuebiya](https://github.com/xuebiya)
