# Avatar System Design Documentation

**Project:** Mockhu Backend  
**Module:** User Avatar Management  
**Version:** 1.0  
**Date:** November 2024  
**Status:** ✅ Implemented (Local Storage) | 🔄 S3 Migration Pending

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Features](#features)
4. [API Endpoints](#api-endpoints)
5. [Image Processing](#image-processing)
6. [Storage Strategy](#storage-strategy)
7. [Error Handling](#error-handling)
8. [Security](#security)
9. [Future Enhancements](#future-enhancements)

---

## 🎯 Overview

The Avatar System allows users to upload, update, and delete profile pictures. Images are validated, processed, and stored with proper error handling and security measures.

### Key Goals
- ✅ Simple user experience (upload/delete)
- ✅ Automatic image optimization (resize/crop)
- ✅ Storage flexibility (local → S3 migration ready)
- ✅ Security (file validation, size limits)
- ✅ Performance (efficient processing)

---

## 🏗️ Architecture

### System Components

```
┌─────────────────────────────────────────────────────────┐
│                     Client (User)                        │
└────────────────────────┬────────────────────────────────┘
                         │ HTTP Multipart Form
                         ▼
┌─────────────────────────────────────────────────────────┐
│              Handler Layer (HTTP)                        │
│  • Parses multipart form                                │
│  • Extracts file data                                   │
│  • Validates authentication                             │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│              Service Layer (Business Logic)              │
│  • Coordinates avatar operations                        │
│  • Manages old avatar cleanup                           │
│  • Database updates                                     │
└────────────────────────┬────────────────────────────────┘
                         │
           ┌─────────────┴─────────────┐
           ▼                           ▼
┌──────────────────────┐    ┌──────────────────────┐
│  Avatar Package      │    │  Repository Layer    │
│  • Validate file     │    │  • Update DB         │
│  • Resize to 400x400 │    │  • Store URL         │
│  • Save to storage   │    └──────────────────────┘
│  • Return URL        │
└──────────┬───────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────┐
│              Storage Layer                               │
│  Current: Local File System (storage/avatars/)          │
│  Future:  AWS S3                                        │
└─────────────────────────────────────────────────────────┘
```

### Package Structure

```
mockhu-app-backend/
├── internal/
│   ├── app/
│   │   └── profile/
│   │       ├── handler.go           # HTTP handlers
│   │       ├── service.go           # Business logic
│   │       ├── repository.go        # DB interface
│   │       └── repository_postgres.go
│   │
│   └── pkg/
│       └── avatar/
│           └── avatar.go            # Image processing utility
│
└── storage/
    └── avatars/                     # Local storage (gitignored)
        └── {uuid}.jpg               # Stored images
```

---

## ⚙️ Features

### 1. Upload Avatar
- **Endpoint:** `POST /v1/users/me/avatar`
- **Authentication:** Required (JWT)
- **Input:** Multipart form with `avatar` field
- **Processing:**
  - Validate file type (JPEG, PNG, WebP)
  - Validate file size (max 5MB)
  - Resize to 400x400 pixels
  - Square crop (center-focused)
  - Save with unique UUID filename
  - Update database with new URL
  - Delete old avatar if exists
- **Output:** Avatar URL

### 2. Delete Avatar
- **Endpoint:** `DELETE /v1/users/me/avatar`
- **Authentication:** Required (JWT)
- **Processing:**
  - Delete file from storage
  - Clear avatar_url in database
- **Output:** Success message

### 3. View Avatar
- **Endpoint:** `GET /avatars/{filename}`
- **Authentication:** Not required (public access)
- **Processing:**
  - Serve static file
- **Output:** Image file

---

## 🔌 API Endpoints

### Upload Avatar

**Request:**
```http
POST /v1/users/me/avatar
Authorization: Bearer {JWT_TOKEN}
Content-Type: multipart/form-data

Form Data:
  avatar: [image file]
```

**Success Response (200):**
```json
{
  "avatar_url": "/avatars/123e4567-e89b-12d3-a456-426614174000.jpg",
  "message": "avatar uploaded successfully"
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `avatar file is required` | No file provided |
| 400 | `file size exceeds 5MB` | File too large |
| 400 | `invalid file type, only JPEG, PNG, and WebP allowed` | Wrong format |
| 401 | `unauthorized` | No/invalid JWT |
| 500 | `failed to upload avatar` | Server error |

---

### Delete Avatar

**Request:**
```http
DELETE /v1/users/me/avatar
Authorization: Bearer {JWT_TOKEN}
```

**Success Response (200):**
```json
{
  "message": "avatar deleted successfully"
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 401 | `unauthorized` | No/invalid JWT |
| 404 | `user not found` | User doesn't exist |
| 500 | `failed to delete avatar` | Server error |

---

### View Avatar

**Request:**
```http
GET /avatars/{filename}
```

**Success Response (200):**
- Content-Type: `image/jpeg`
- Body: Image binary data

**Error Response (404):**
- File not found

---

## 🖼️ Image Processing

### Processing Pipeline

```
1. Upload File
   ↓
2. Validate File Type
   • Check magic bytes
   • PNG: \x89PNG\r\n\x1a\n
   • JPEG: 0xFF 0xD8 0xFF
   • WebP: RIFF....WEBP
   ↓
3. Validate File Size
   • Max: 5MB (5,242,880 bytes)
   ↓
4. Decode Image
   • PNG: png.Decode()
   • JPEG: jpeg.Decode()
   • WebP: imaging.Decode()
   ↓
5. Resize & Crop
   • Target: 400x400 pixels
   • Method: imaging.Fill()
   • Algorithm: Lanczos (high quality)
   • Anchor: Center
   ↓
6. Encode & Save
   • Format: JPEG (always)
   • Quality: 90%
   • Filename: {UUID}.jpg
   ↓
7. Return URL
   • Format: /avatars/{UUID}.jpg
```

### Technical Specifications

**Supported Input Formats:**
- JPEG/JPG
- PNG
- WebP

**Output Format:**
- Always JPEG (for consistency)
- Quality: 90%
- Size: 400x400 pixels

**Processing Library:**
- `github.com/disintegration/imaging`
- Pure Go implementation
- No external dependencies (libvips, etc.)

**Resize Algorithm:**
```go
imaging.Fill(img, 400, 400, imaging.Center, imaging.Lanczos)
```
- **Fill:** Crops to exact size
- **Center:** Focuses on image center
- **Lanczos:** High-quality resampling

---

## 💾 Storage Strategy

### Current: Local File System

**Directory Structure:**
```
storage/
└── avatars/
    ├── 123e4567-e89b-12d3-a456-426614174000.jpg
    ├── 987fcdeb-51a2-43f8-b9c3-123456789abc.jpg
    └── ...
```

**Filename Format:**
- UUID v4 + `.jpg` extension
- Example: `123e4567-e89b-12d3-a456-426614174000.jpg`
- Generated by: `github.com/google/uuid`

**URL Format:**
- `/avatars/{UUID}.jpg`
- Served by Fiber static middleware
- Example: `http://localhost:8085/avatars/123e4567-e89b-12d3-a456-426614174000.jpg`

**Advantages:**
- ✅ Simple implementation
- ✅ No external dependencies
- ✅ Fast local access
- ✅ Good for development/testing

**Disadvantages:**
- ❌ Not scalable (single server)
- ❌ No CDN benefits
- ❌ Backup complexity
- ❌ Disk space management

---

### Future: AWS S3

**Migration Strategy:**

```go
// TODO: S3 Implementation
func saveToS3(img image.Image, filename string) (string, error) {
    // 1. Encode image to bytes buffer
    buf := new(bytes.Buffer)
    jpeg.Encode(buf, img, &jpeg.Options{Quality: 90})
    
    // 2. Upload to S3 bucket
    sess := session.New(&aws.Config{
        Region: aws.String("us-east-1"),
    })
    svc := s3.New(sess)
    
    _, err := svc.PutObject(&s3.PutObjectInput{
        Bucket: aws.String("mockhu-avatars"),
        Key:    aws.String(filename),
        Body:   bytes.NewReader(buf.Bytes()),
        ContentType: aws.String("image/jpeg"),
        ACL:    aws.String("public-read"),
    })
    
    // 3. Return S3 URL
    return fmt.Sprintf("https://mockhu-avatars.s3.amazonaws.com/%s", filename), nil
}

func deleteFromS3(filename string) error {
    // Delete from S3 bucket
    sess := session.New(&aws.Config{
        Region: aws.String("us-east-1"),
    })
    svc := s3.New(sess)
    
    _, err := svc.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String("mockhu-avatars"),
        Key:    aws.String(filename),
    })
    
    return err
}
```

**S3 Configuration:**

| Setting | Value |
|---------|-------|
| **Bucket Name** | `mockhu-avatars` |
| **Region** | `us-east-1` (or closest) |
| **ACL** | `public-read` |
| **Versioning** | Disabled |
| **Lifecycle** | Delete after 90 days of user deletion |
| **CORS** | Allow GET from app domains |

**CloudFront CDN:**
- Distribution: `https://cdn.mockhu.com/avatars/`
- Cache: 1 year (immutable filenames)
- SSL: Required
- Geo-restriction: None

**Cost Estimation (100K users):**
- Storage: ~40GB @ $0.023/GB = **$0.92/month**
- Requests: ~1M GET @ $0.0004/1K = **$0.40/month**
- Transfer: ~100GB @ $0.09/GB = **$9.00/month**
- **Total: ~$10.32/month**

---

## 🛡️ Security

### File Validation

**1. File Type Validation (Magic Bytes)**
```go
// Check PNG
if len(data) >= 8 && string(data[0:8]) == "\x89PNG\r\n\x1a\n" {
    return "png"
}
// Check JPEG
if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
    return "jpeg"
}
// Check WebP
if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
    return "webp"
}
```
- **Why:** File extension can be faked
- **Method:** Read actual file headers
- **Blocks:** Executable files, scripts, etc.

**2. File Size Limit**
```go
const MaxFileSize = 5 * 1024 * 1024 // 5MB

if len(fileData) > MaxFileSize {
    return "", ErrFileTooBig
}
```
- **Limit:** 5MB maximum
- **Reason:** Prevent DoS attacks, storage abuse
- **Enforcement:** Before processing

**3. Image Processing Validation**
```go
img, err := decodeImage(bytes.NewReader(fileData), fileType)
if err != nil {
    return "", fmt.Errorf("%w: %v", ErrProcessingFailed, err)
}
```
- **Validates:** File is actually a valid image
- **Blocks:** Malformed/corrupted files

### Authentication

- **Required:** JWT token in Authorization header
- **Validation:** Via middleware
- **Scope:** User can only upload/delete their own avatar

### Access Control

| Operation | Authentication | Authorization |
|-----------|----------------|---------------|
| Upload Avatar | Required | Self only |
| Delete Avatar | Required | Self only |
| View Avatar | Not Required | Public |

### File System Security

```bash
# Storage directory permissions
chmod 755 storage/avatars/

# Uploaded file permissions
chmod 644 storage/avatars/{uuid}.jpg
```

- **Directory:** Read/execute for all (755)
- **Files:** Read for all, write for owner (644)
- **No execution:** Images cannot be executed

### Injection Prevention

- **Filename:** Generated UUID (no user input)
- **Path Traversal:** Prevented by UUID-only naming
- **Content-Type:** Fixed (image/jpeg)

---

## 📊 Performance

### Metrics

**Processing Time:**
- Validate: <1ms
- Decode: 10-50ms (depends on file size)
- Resize: 20-100ms (depends on original size)
- Encode: 10-30ms
- **Total: ~50-180ms per upload**

**Storage:**
- Average avatar size: ~40-60KB (after processing)
- 100,000 users = ~4-6GB storage

**Bandwidth:**
- Upload: ~500KB-5MB per request
- Download: ~40-60KB per request
- CDN caching significantly reduces bandwidth

### Optimization

**Current:**
- ✅ JPEG encoding (90% quality, small file size)
- ✅ Fixed dimensions (400x400, predictable)
- ✅ Efficient algorithm (Lanczos)

**Future:**
- 🔄 WebP output (smaller files)
- 🔄 Progressive JPEG
- 🔄 Responsive sizes (100x100, 200x200, 400x400)
- 🔄 Image optimization service (imgix, Cloudinary)

---

## ⚠️ Error Handling

### Error Types

```go
var (
    ErrFileTooBig       = errors.New("file size exceeds 5MB")
    ErrInvalidFileType  = errors.New("invalid file type, only JPEG, PNG, and WebP allowed")
    ErrProcessingFailed = errors.New("failed to process image")
)
```

### Error Flow

```
User Upload
    ↓
File Too Large? → 400 Bad Request
    ↓
Invalid Type? → 400 Bad Request
    ↓
Processing Failed? → 500 Internal Error
    ↓
Database Failed? → 500 Internal Error
    ↓  (Rollback: Delete uploaded file)
Success → 200 OK
```

### Rollback Strategy

**On Database Update Failure:**
```go
err = s.profileRepo.UpdateAvatar(ctx, userID, avatarURL)
if err != nil {
    // Rollback: Delete uploaded file
    _ = avatar.DeleteAvatar(avatarURL)
    return nil, fmt.Errorf("failed to update avatar: %w", err)
}
```

**On Service Error:**
- New file deleted automatically (not persisted)
- Old avatar remains unchanged
- User can retry

---

## 🧪 Testing

### Manual Testing

**1. Upload Valid Image:**
```bash
TOKEN="your_jwt_token"
curl -X POST http://localhost:8085/v1/users/me/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@test_image.jpg"
```

**2. Upload Invalid Type:**
```bash
curl -X POST http://localhost:8085/v1/users/me/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@document.pdf"

# Expected: 400 Bad Request
```

**3. Upload Too Large:**
```bash
# Create 6MB file
dd if=/dev/zero of=large.jpg bs=1M count=6

curl -X POST http://localhost:8085/v1/users/me/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@large.jpg"

# Expected: 400 Bad Request
```

**4. Delete Avatar:**
```bash
curl -X DELETE http://localhost:8085/v1/users/me/avatar \
  -H "Authorization: Bearer $TOKEN"

# Expected: 200 OK
```

**5. View Avatar:**
```bash
curl http://localhost:8085/avatars/{uuid}.jpg --output downloaded.jpg

# Expected: JPEG image file
```

### Test Cases

| Test Case | Expected Result |
|-----------|----------------|
| Upload JPEG | ✅ Success |
| Upload PNG | ✅ Success |
| Upload WebP | ✅ Success |
| Upload PDF | ❌ 400 Error |
| Upload 6MB | ❌ 400 Error |
| Upload without auth | ❌ 401 Error |
| Replace existing | ✅ Old deleted |
| Delete avatar | ✅ File removed |
| View public avatar | ✅ Image served |

---

## 🚀 Future Enhancements

### Phase 1: S3 Migration (High Priority)
- [ ] AWS SDK integration
- [ ] S3 bucket setup
- [ ] CloudFront CDN
- [ ] Migration script for existing avatars
- [ ] Environment-based storage (local dev, S3 prod)

### Phase 2: Advanced Processing
- [ ] Multiple sizes (thumbnails)
  - 100x100 (chat)
  - 200x200 (cards)
  - 400x400 (profile)
- [ ] WebP output format
- [ ] Progressive JPEG
- [ ] Image optimization service

### Phase 3: Enhanced Features
- [ ] Avatar cropping tool (client-side)
- [ ] Filters/effects
- [ ] Default avatars (generated)
- [ ] Avatar history (restore previous)
- [ ] Face detection (auto-crop)

### Phase 4: Performance
- [ ] Lazy loading
- [ ] Responsive images (srcset)
- [ ] Image compression API
- [ ] Background processing (queue)

---

## 📝 Code Reference

### Key Files

| File | Purpose | Lines |
|------|---------|-------|
| `internal/pkg/avatar/avatar.go` | Image processing utility | ~200 |
| `internal/app/profile/service.go` | Avatar business logic | ~60 |
| `internal/app/profile/handler.go` | HTTP handlers | ~90 |
| `cmd/api/main.go` | Static file serving | 1 |

### Dependencies

```go
// go.mod
github.com/disintegration/imaging v1.6.2  // Image processing
github.com/google/uuid v1.3.0             // UUID generation
```

---

## 📞 Support & Maintenance

### Common Issues

**Issue:** "Avatar not showing"
- **Check:** File exists in `storage/avatars/`
- **Check:** Correct URL in database
- **Check:** Static file serving enabled

**Issue:** "Upload fails silently"
- **Check:** Disk space available
- **Check:** Directory permissions (755)
- **Check:** Server logs for errors

**Issue:** "Images look blurry"
- **Solution:** Increase JPEG quality (90 → 95)
- **Location:** `internal/pkg/avatar/avatar.go:149`

### Monitoring

**Metrics to Track:**
- Upload success rate
- Average processing time
- Storage usage
- Failed uploads (by error type)
- CDN hit rate (S3)

---

## 📚 References

- [disintegration/imaging docs](https://github.com/disintegration/imaging)
- [AWS S3 Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/optimizing-performance.html)
- [Image Optimization Guide](https://web.dev/fast/#optimize-your-images)
- [OWASP File Upload Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html)

---

**Last Updated:** November 27, 2024  
**Version:** 1.0  
**Status:** ✅ Production Ready (Local Storage) | 🔄 S3 Migration Planned

