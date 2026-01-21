# Cloudreve v4 with Cloud189 Driver

## Release Information

This is a custom build of Cloudreve v4 with **China Telecom Cloud189 (天翼云盘)** storage driver support.

### New Features

- ✨ **Cloud189 Storage Driver** - Full support for China Telecom Cloud189 cloud storage
  - File upload (with chunked upload support)
  - File download
  - File deletion
  - File listing
  - Automatic login and session management

### Supported Storage Drivers

- Local Storage
- Aliyun OSS
- Tencent COS
- Qiniu
- Upyun
- OneDrive
- S3 Compatible
- Kingsoft KS3
- Huawei OBS
- **Cloud189 (NEW)** ⭐

## Download

### Pre-compiled Binaries

| Platform | Architecture | File | Size |
|----------|-------------|------|------|
| Windows | amd64 | cloudreve-windows-amd64.exe | ~61 MB |
| Linux | amd64 | cloudreve-linux-amd64 | ~60 MB |

## Installation

### Windows

1. Download `cloudreve-windows-amd64.exe`
2. Run the executable
3. Access `http://localhost:5212` in your browser
4. Follow the setup wizard

### Linux

```bash
# Download
wget https://github.com/xuebiya/cloudreve/releases/download/v4.0.0-cloud189/cloudreve-linux-amd64

# Make executable
chmod +x cloudreve-linux-amd64

# Run
./cloudreve-linux-amd64
```

## Configuration

### Adding Cloud189 Storage Policy

1. Log in to Cloudreve admin panel
2. Navigate to "Storage Policies"
3. Click "Add Storage Policy"
4. Fill in the following information:
   - **Policy Name**: Cloud189
   - **Storage Type**: cloud189
   - **AccessKey**: Your Cloud189 account (phone number or email)
   - **SecretKey**: Your Cloud189 password
   - **Is Private**: Yes

### Database Configuration (Alternative)

```sql
INSERT INTO storage_policies (name, type, access_key, secret_key, is_private, settings, created_at, updated_at)
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

## Usage Notes

### Security

- Account credentials are stored in the database
- Ensure your database is properly secured
- Use strong passwords for both Cloudreve and Cloud189 accounts

### Limitations

- Folder operations (create, move, rename) are not yet supported
- Thumbnail generation is not supported
- Media metadata extraction is not supported
- File path management is simplified (root directory only)

### Performance

- Default chunk size: 10MB
- Supports automatic re-login on session expiration
- Optimized for Chinese network environment

## Building from Source

```bash
# Clone repository
git clone https://github.com/xuebiya/cloudreve.git
cd cloudreve

# Build
go build -ldflags "-s -w" -o cloudreve ./main.go
```

## Technical Details

### Cloud189 Driver Implementation

The Cloud189 driver is ported from the [OpenList](https://github.com/OpenListTeam/OpenList) project and adapted for Cloudreve v4 architecture.

**Key Features:**
- RSA encryption for login credentials
- HMAC-SHA1 signature for upload requests
- AES encryption for request parameters
- MD5 verification for uploaded files
- Automatic session management

**Files:**
- `pkg/filemanager/driver/cloud189/cloud189.go` - Main driver implementation
- `pkg/filemanager/driver/cloud189/types.go` - Type definitions
- `pkg/filemanager/driver/cloud189/util.go` - Utility functions
- `pkg/filemanager/driver/cloud189/login.go` - Login logic

## Troubleshooting

### Login Failed

- Verify your Cloud189 account credentials
- Ensure your account doesn't have 2FA enabled
- Check if your account is locked due to too many login attempts

### Upload Failed

- Check the log files for detailed error messages
- Verify network connectivity to Cloud189 servers
- Ensure your account has sufficient storage space

### File Not Found

- The driver uses folder IDs instead of paths
- Currently only supports root directory operations

## Credits

- **Cloudreve**: [https://github.com/cloudreve/Cloudreve](https://github.com/cloudreve/Cloudreve)
- **OpenList**: [https://github.com/OpenListTeam/OpenList](https://github.com/OpenListTeam/OpenList)
- **Cloud189 Driver**: Ported from OpenList by xuebiya

## License

This project follows the Cloudreve project license.

## Support

For issues and questions:
- GitHub Issues: [https://github.com/xuebiya/cloudreve/issues](https://github.com/xuebiya/cloudreve/issues)
- Original Cloudreve: [https://docs.cloudreve.org/](https://docs.cloudreve.org/)

## Changelog

### v4.0.0-cloud189 (2026-01-21)

- ✨ Added Cloud189 storage driver support
- ✅ File upload with chunked upload
- ✅ File download
- ✅ File deletion
- ✅ File listing
- ✅ Automatic login and session management
- 📝 Complete documentation
- 🔧 Build scripts for Windows and Linux

---

**Note**: This is an unofficial build with Cloud189 driver support. For the official Cloudreve release, please visit [https://github.com/cloudreve/Cloudreve](https://github.com/cloudreve/Cloudreve)
