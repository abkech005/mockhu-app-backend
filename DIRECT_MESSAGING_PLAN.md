# 💬 Direct Messaging (DM) Feature - Complete Implementation Plan

**Feature:** Direct Messaging System  
**Priority:** P1 (Critical Path)  
**Status:** Planning Phase  
**Estimated Time:** Not Bounded (Flexible)  
**Created:** November 29, 2024

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Core Requirements](#core-requirements)
3. [Technical Architecture](#technical-architecture)
4. [Database Schema Design](#database-schema-design)
5. [API Endpoints](#api-endpoints)
6. [Business Logic & Rules](#business-logic--rules)
7. [Privacy & Security](#privacy--security)
8. [Real-Time Delivery Strategy](#real-time-delivery-strategy)
9. [File Structure](#file-structure)
10. [Implementation Phases](#implementation-phases)
11. [Testing Strategy](#testing-strategy)
12. [Performance Considerations](#performance-considerations)
13. [Future Enhancements](#future-enhancements)

---

## 📖 Overview

### Feature Description
A complete 1-on-1 direct messaging system that allows students to communicate privately with text messages and images, with real-time delivery, message history, and unread count tracking.

### Key Goals
- ✅ Enable private 1-on-1 conversations between students
- ✅ Support text messages
- ✅ Support image messages (JPEG, PNG, GIF, WebP - up to 5 per message)
- ✅ Support file attachments (PDF, DOCX, XLSX, PPTX, TXT, ZIP, etc.)
- ✅ Real-time message delivery and notifications
- ✅ Unread message tracking and counts
- ✅ Respect user privacy settings (`who_can_message`)
- ✅ Block/unblock functionality
- ✅ Message history with pagination
- ✅ File upload with security validation
- ✅ Typing indicators (optional)
- ✅ Message read receipts (optional)
- ✅ Message deletion (optional)

### Supported Attachments

| Type | Formats | Max Size | Quantity |
|------|---------|----------|----------|
| **Images** | JPEG, PNG, GIF, WebP | 10MB each | Up to 5 per message |
| **Documents** | PDF, DOCX, XLSX, PPTX, TXT | 25MB each | Up to 5 per message |
| **Archives** | ZIP, RAR, 7Z | 50MB each | Up to 5 per message |
| **Data Files** | CSV, JSON, XML | 10MB each | Up to 5 per message |

**Total per message:** Maximum 5 files, combined size limit 100MB

### Success Metrics
- Message delivery latency < 1 second
- Support for 10,000+ concurrent users
- 99.9% message delivery success rate
- Support conversations with 1000+ messages efficiently

---

## 🎯 Core Requirements

### Must Have (MVP)
1. **Conversations**
   - Create/get conversations between two users
   - List all conversations for a user (inbox)
   - Get conversation details with participants

2. **Messages**
   - Send text messages
   - Send image messages (up to 5 images per message)
   - Send file attachments (PDF, DOCX, XLSX, PPT, TXT, ZIP, etc.)
   - Get message history (paginated)
   - Display sender information with each message

3. **Unread Tracking**
   - Track unread message count per conversation
   - Mark messages as read
   - Show total unread count across all conversations
   - Mark conversation as read (all messages)

4. **Privacy Controls**
   - Respect `who_can_message` setting (everyone/followers/none)
   - Check if user can send message to another user
   - Block users from messaging

5. **Real-Time Delivery**
   - WebSocket connection for real-time updates
   - Push notifications for new messages (future)
   - Online/offline status tracking

6. **Message Status**
   - Sent
   - Delivered
   - Read (optional for MVP)

### Nice to Have (Future)
- Message editing
- Message reactions (like, love, etc.)
- Reply to specific messages
- Forward messages
- Voice messages
- Video messages
- Message search
- Delete for everyone
- Typing indicators
- Read receipts toggle (privacy)
- Group messaging

---

## 🏗️ Technical Architecture

### Architecture Pattern
Following the existing DDD (Domain-Driven Design) pattern:
```
Repository Layer → Service Layer → Handler Layer → Routes Layer
```

### Technology Stack
- **Database:** PostgreSQL (existing)
- **Real-Time:** WebSocket (gorilla/websocket)
- **File Storage:** Local storage (same as avatars, S3-ready)
- **Image Processing:** github.com/disintegration/imaging (existing)
- **Authentication:** JWT Bearer tokens (existing)

### Key Components

#### 1. Conversation Manager
- Manages conversation creation and retrieval
- Ensures one conversation between two users
- Handles conversation listing and pagination

#### 2. Message Manager
- Handles message sending, retrieval, deletion
- Image upload and storage
- Message pagination and filtering

#### 3. Unread Counter
- Tracks unread messages per conversation
- Calculates total unread count
- Marks messages/conversations as read

#### 4. Privacy Checker
- Validates if user can message another user
- Checks blocked status
- Respects privacy settings

#### 5. WebSocket Manager
- Manages WebSocket connections
- Routes messages to recipients
- Handles connection lifecycle (connect/disconnect)
- Connection pooling and cleanup

#### 6. Notification Manager (Future)
- Push notification integration
- Email notification for messages (optional)

---

## 🗄️ Database Schema Design

### Table 1: `conversations`
Stores 1-on-1 conversation metadata between two users.

```sql
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Participants (always exactly 2 users)
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Last message info (denormalized for performance)
    last_message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    last_message_text TEXT,
    last_message_sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    last_message_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT unique_conversation UNIQUE (user1_id, user2_id),
    CONSTRAINT different_users CHECK (user1_id != user2_id),
    CONSTRAINT ordered_users CHECK (user1_id < user2_id)
);

-- Indexes
CREATE INDEX idx_conversations_user1 ON conversations(user1_id);
CREATE INDEX idx_conversations_user2 ON conversations(user2_id);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC);
CREATE INDEX idx_conversations_updated_at ON conversations(updated_at DESC);
```

**Design Notes:**
- `user1_id < user2_id` ensures consistent ordering (prevents duplicate conversations)
- Denormalized last message for efficient inbox listing
- Separate indexes on both user columns for quick lookups

---

### Table 2: `messages`
Stores individual messages within conversations.

```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Conversation reference
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    
    -- Message content
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_type VARCHAR(20) NOT NULL CHECK (message_type IN ('text', 'image', 'file')),
    content TEXT, -- Text content, image caption, or file description
    
    -- Media attachments (JSON array for multiple files)
    attachments JSONB, -- Array of {type, url, filename, size, mime_type}
    
    -- Message status
    status VARCHAR(20) NOT NULL DEFAULT 'sent' CHECK (status IN ('sent', 'delivered', 'read')),
    
    -- Read tracking
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP,
    
    -- Soft delete
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT attachments_required CHECK (
        (message_type IN ('image', 'file') AND attachments IS NOT NULL) OR 
        (message_type = 'text')
    )
);

-- Indexes
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_is_read ON messages(is_read) WHERE is_read = FALSE;
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
```

**Design Notes:**
- Composite index on `(conversation_id, created_at DESC)` for efficient message history queries
- Partial index on `is_read = FALSE` for quick unread count queries
- Soft delete support for "delete for me" functionality
- Message status tracking for future read receipts
- **JSONB attachments field**: Stores array of file metadata for flexible multi-file support
  ```json
  [
    {
      "type": "image",
      "url": "/uploads/messages/2024/11/uuid.jpg",
      "filename": "photo.jpg",
      "size": 2048576,
      "mime_type": "image/jpeg"
    },
    {
      "type": "file",
      "url": "/uploads/messages/2024/11/uuid.pdf",
      "filename": "document.pdf",
      "size": 512000,
      "mime_type": "application/pdf"
    }
  ]
  ```

---

### Table 3: `blocked_users`
Tracks users who have blocked each other.

```sql
CREATE TABLE blocked_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Who blocked whom
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Metadata
    reason TEXT, -- Optional reason for blocking
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE(blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);

-- Indexes
CREATE INDEX idx_blocked_users_blocker ON blocked_users(blocker_id);
CREATE INDEX idx_blocked_users_blocked ON blocked_users(blocked_id);
```

**Design Notes:**
- One-way blocking (A blocks B, but B can still see A's messages until B also blocks A)
- Composite unique constraint prevents duplicate blocks
- Fast lookup for blocking checks

---

### Table 4: `conversation_participants` (Optional - For Future Group Chat)
For future scalability to support group chats.

```sql
-- OPTIONAL: For future group chat support
CREATE TABLE conversation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Participant status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    left_at TIMESTAMP,
    
    -- Unread count (per user)
    unread_count INTEGER NOT NULL DEFAULT 0,
    last_read_message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    last_read_at TIMESTAMP,
    
    -- Timestamps
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE(conversation_id, user_id)
);

-- Indexes
CREATE INDEX idx_participants_conversation ON conversation_participants(conversation_id);
CREATE INDEX idx_participants_user ON conversation_participants(user_id);
CREATE INDEX idx_participants_unread ON conversation_participants(unread_count) WHERE unread_count > 0;
```

**Design Notes:**
- Not needed for MVP (1-on-1 only)
- Enables future group chat functionality
- Stores per-user unread counts

---

### Migration Files Needed

```
migrations/
├── 000015_create_conversations.up.sql
├── 000015_create_conversations.down.sql
├── 000016_create_messages.up.sql
├── 000016_create_messages.down.sql
├── 000017_create_blocked_users.up.sql
├── 000017_create_blocked_users.down.sql
```

---

## 🔌 API Endpoints

### 1. Conversations

#### `POST /v1/conversations`
**Description:** Create or get existing conversation with another user  
**Auth:** Required  
**Request Body:**
```json
{
  "recipient_id": "uuid"
}
```
**Response:** 201 Created / 200 OK
```json
{
  "success": true,
  "data": {
    "id": "conversation-uuid",
    "participant": {
      "id": "user-uuid",
      "username": "johndoe",
      "full_name": "John Doe",
      "avatar_url": "/avatars/xxx.jpg"
    },
    "last_message": {
      "id": "message-uuid",
      "content": "Hey, how are you?",
      "sender_id": "user-uuid",
      "created_at": "2024-11-29T10:00:00Z"
    },
    "unread_count": 3,
    "created_at": "2024-11-28T10:00:00Z",
    "updated_at": "2024-11-29T10:00:00Z"
  }
}
```

---

#### `GET /v1/conversations`
**Description:** Get all conversations (inbox) for authenticated user  
**Auth:** Required  
**Query Params:**
- `page` (default: 1)
- `limit` (default: 20, max: 50)
- `unread_only` (boolean, default: false)

**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "conversations": [
      {
        "id": "conversation-uuid",
        "participant": {
          "id": "user-uuid",
          "username": "johndoe",
          "full_name": "John Doe",
          "avatar_url": "/avatars/xxx.jpg",
          "is_online": true
        },
        "last_message": {
          "id": "message-uuid",
          "content": "See you tomorrow!",
          "sender_id": "user-uuid",
          "message_type": "text",
          "created_at": "2024-11-29T15:30:00Z"
        },
        "unread_count": 0,
        "updated_at": "2024-11-29T15:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45,
      "total_pages": 3,
      "total_unread": 12
    }
  }
}
```

---

#### `GET /v1/conversations/:conversationId`
**Description:** Get conversation details  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "id": "conversation-uuid",
    "participant": {
      "id": "user-uuid",
      "username": "johndoe",
      "full_name": "John Doe",
      "avatar_url": "/avatars/xxx.jpg",
      "is_online": true,
      "last_seen": "2024-11-29T15:30:00Z"
    },
    "unread_count": 3,
    "is_blocked": false,
    "created_at": "2024-11-28T10:00:00Z"
  }
}
```

---

#### `DELETE /v1/conversations/:conversationId`
**Description:** Delete conversation (soft delete, removes from user's inbox)  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "message": "Conversation deleted successfully"
}
```

---

### 2. Messages

#### `POST /v1/conversations/:conversationId/messages`
**Description:** Send a message (text or image)  
**Auth:** Required  
**Request Body (Text):**
```json
{
  "message_type": "text",
  "content": "Hey! How's your project going?"
}
```
**Request Body (Image):**
```json
{
  "message_type": "image",
  "content": "Check out these photos!",
  "attachments": [
    {
      "type": "image",
      "url": "/uploads/messages/2024/11/uuid1.jpg",
      "filename": "sunset.jpg",
      "size": 2048576,
      "mime_type": "image/jpeg"
    },
    {
      "type": "image",
      "url": "/uploads/messages/2024/11/uuid2.jpg",
      "filename": "beach.jpg",
      "size": 1536000,
      "mime_type": "image/jpeg"
    }
  ]
}
```
**Request Body (File):**
```json
{
  "message_type": "file",
  "content": "Here's the assignment PDF",
  "attachments": [
    {
      "type": "file",
      "url": "/uploads/messages/2024/11/uuid.pdf",
      "filename": "CS101_Assignment.pdf",
      "size": 512000,
      "mime_type": "application/pdf"
    }
  ]
}
```
**Response:** 201 Created
```json
{
  "success": true,
  "data": {
    "id": "message-uuid",
    "conversation_id": "conversation-uuid",
    "sender_id": "user-uuid",
    "message_type": "text",
    "content": "Hey! How's your project going?",
    "attachments": null,
    "status": "sent",
    "is_read": false,
    "created_at": "2024-11-29T16:00:00Z"
  }
}
```

---

#### `GET /v1/conversations/:conversationId/messages`
**Description:** Get message history for a conversation  
**Auth:** Required  
**Query Params:**
- `page` (default: 1)
- `limit` (default: 50, max: 100)
- `before_id` (UUID, for cursor-based pagination)
- `after_id` (UUID, for real-time updates)

**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": "message-uuid",
        "sender": {
          "id": "user-uuid",
          "username": "johndoe",
          "full_name": "John Doe",
          "avatar_url": "/avatars/xxx.jpg"
        },
        "message_type": "text",
        "content": "Great! Let's meet tomorrow.",
        "attachments": null,
        "is_read": true,
        "read_at": "2024-11-29T16:05:00Z",
        "created_at": "2024-11-29T16:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 50,
      "has_more": true,
      "next_cursor": "message-uuid"
    }
  }
}
```

---

#### `DELETE /v1/messages/:messageId`
**Description:** Delete a message (soft delete, "delete for me")  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "message": "Message deleted successfully"
}
```

---

### 3. Unread Management

#### `POST /v1/conversations/:conversationId/read`
**Description:** Mark conversation as read (all messages)  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "message": "Conversation marked as read"
}
```

---

#### `POST /v1/messages/:messageId/read`
**Description:** Mark single message as read  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "message": "Message marked as read"
}
```

---

#### `GET /v1/conversations/unread-count`
**Description:** Get total unread message count across all conversations  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "total_unread": 15,
    "unread_conversations": 5
  }
}
```

---

### 4. Privacy & Blocking

#### `POST /v1/users/:userId/block`
**Description:** Block a user from messaging  
**Auth:** Required  
**Request Body:**
```json
{
  "reason": "Optional reason for blocking"
}
```
**Response:** 200 OK
```json
{
  "success": true,
  "message": "User blocked successfully"
}
```

---

#### `DELETE /v1/users/:userId/block`
**Description:** Unblock a user  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "message": "User unblocked successfully"
}
```

---

#### `GET /v1/users/blocked`
**Description:** Get list of blocked users  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "blocked_users": [
      {
        "id": "user-uuid",
        "username": "spammer123",
        "full_name": "Spammer User",
        "blocked_at": "2024-11-20T10:00:00Z"
      }
    ]
  }
}
```

---

#### `GET /v1/users/:userId/can-message`
**Description:** Check if current user can message another user  
**Auth:** Required  
**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "can_message": true,
    "reason": null
  }
}
```
OR
```json
{
  "success": true,
  "data": {
    "can_message": false,
    "reason": "User has blocked you"
  }
}
```

---

### 5. File Upload

#### `POST /v1/messages/upload`
**Description:** Upload files (images or documents) for messaging  
**Auth:** Required  
**Request:** multipart/form-data
- `files` (multiple files, max 5 files)
- `type` (string: "image" or "file")

**Supported File Types:**
- **Images:** JPEG, PNG, GIF, WebP (max 10MB each)
- **Documents:** PDF, DOCX, XLSX, PPTX, TXT (max 25MB each)
- **Archives:** ZIP, RAR (max 50MB each)
- **Other:** CSV, JSON, XML (max 10MB each)

**Response:** 200 OK
```json
{
  "success": true,
  "data": {
    "attachments": [
      {
        "type": "image",
        "url": "/uploads/messages/2024/11/uuid1.jpg",
        "thumbnail_url": "/uploads/messages/thumbnails/uuid1.jpg",
        "filename": "photo.jpg",
        "size": 2048576,
        "mime_type": "image/jpeg"
      },
      {
        "type": "file",
        "url": "/uploads/messages/2024/11/uuid2.pdf",
        "filename": "document.pdf",
        "size": 512000,
        "mime_type": "application/pdf"
      }
    ]
  }
}
```

---

### 6. WebSocket

#### `WS /v1/ws/messages`
**Description:** WebSocket endpoint for real-time messaging  
**Auth:** Required (JWT in query param or header)  
**Protocol:** WebSocket

**Client → Server Messages:**

**1. Send Message**
```json
{
  "type": "message",
  "conversation_id": "uuid",
  "message_type": "text",
  "content": "Hello!"
}
```

**2. Typing Indicator**
```json
{
  "type": "typing",
  "conversation_id": "uuid",
  "is_typing": true
}
```

**3. Mark as Read**
```json
{
  "type": "read",
  "message_id": "uuid"
}
```

**Server → Client Messages:**

**1. New Message**
```json
{
  "type": "new_message",
  "data": {
    "id": "message-uuid",
    "conversation_id": "conversation-uuid",
    "sender": {
      "id": "user-uuid",
      "username": "johndoe",
      "full_name": "John Doe",
      "avatar_url": "/avatars/xxx.jpg"
    },
    "message_type": "text",
    "content": "Hello!",
    "created_at": "2024-11-29T16:00:00Z"
  }
}
```

**2. Message Read**
```json
{
  "type": "message_read",
  "data": {
    "message_id": "uuid",
    "conversation_id": "uuid",
    "read_by": "user-uuid",
    "read_at": "2024-11-29T16:05:00Z"
  }
}
```

**3. Typing Indicator**
```json
{
  "type": "user_typing",
  "data": {
    "conversation_id": "uuid",
    "user_id": "user-uuid",
    "is_typing": true
  }
}
```

**4. User Online Status**
```json
{
  "type": "user_status",
  "data": {
    "user_id": "uuid",
    "is_online": true,
    "last_seen": "2024-11-29T16:00:00Z"
  }
}
```

---

### Summary: API Endpoints Count

| Category | Endpoints | Total |
|----------|-----------|-------|
| **Conversations** | GET, POST, GET/:id, DELETE/:id | 4 |
| **Messages** | POST, GET, DELETE/:id | 3 |
| **Unread** | POST read, POST message read, GET count | 3 |
| **Privacy** | POST block, DELETE block, GET blocked, GET can-message | 4 |
| **Upload** | POST upload (images/files) | 1 |
| **WebSocket** | WS endpoint | 1 |
| **TOTAL** | | **16 endpoints** |

---

## 🔐 Business Logic & Rules

### Conversation Rules

1. **Conversation Creation**
   - One conversation between two users only
   - Automatically created when first message is sent
   - Use `user1_id < user2_id` ordering to prevent duplicates
   - Cannot create conversation with self

2. **Conversation Listing**
   - Order by `updated_at` (most recent first)
   - Show unread count per conversation
   - Filter by `unread_only` if requested
   - Soft-deleted conversations excluded

3. **Conversation Access**
   - Only participants can access conversation
   - Return 403 if user is not a participant

---

### Message Rules

1. **Message Sending**
   - Check if sender can message recipient (privacy + blocking)
   - **Text messages:** 1-10,000 characters
   - **Image messages:** Up to 5 images, max 10MB each, formats: JPEG, PNG, GIF, WebP
   - **File messages:** Up to 5 files, max 25MB each for documents, 50MB for archives
   - **Supported file types:**
     - Documents: PDF, DOCX, XLSX, PPTX, TXT
     - Archives: ZIP, RAR, 7Z
     - Data: CSV, JSON, XML
   - **File validation:**
     - Verify MIME type server-side
     - Scan file extension
     - Check file size limits
     - Reject executable files (.exe, .sh, .bat)
   - Create conversation if doesn't exist
   - Store attachments metadata in JSONB field
   - Update conversation's `last_message_at` and denormalized fields
   - Increment recipient's unread count
   - Trigger WebSocket notification to recipient

2. **Message Retrieval**
   - Only conversation participants can view messages
   - Paginate with cursor-based pagination (for efficiency)
   - Order by `created_at DESC` (newest first for loading)
   - Exclude soft-deleted messages

3. **Message Deletion**
   - Only sender can delete their own message
   - Soft delete (set `is_deleted = true`)
   - Don't update conversation's last message if deleted
   - Future: "Delete for everyone" (hard delete within 1 hour)

---

### Unread Tracking

1. **Marking as Read**
   - Only recipient can mark message as read
   - Set `is_read = true` and `read_at = NOW()`
   - Decrement unread count
   - Trigger WebSocket "read receipt" to sender

2. **Auto-Mark Read**
   - When user fetches messages, auto-mark unread messages as read
   - When user opens conversation (frontend calls read endpoint)

3. **Unread Count**
   - Count unread messages per conversation
   - Total unread across all conversations
   - Efficient queries using partial indexes

---

### Privacy & Blocking

1. **Privacy Settings**
   - `who_can_message` values:
     - `everyone`: Anyone can message
     - `followers`: Only followers can start conversation
     - `none`: No one can start new conversation (existing conversations work)

2. **Blocking Logic**
   - If A blocks B:
     - B cannot send new messages to A
     - B cannot see A's online status
     - Existing conversation hidden for B
     - A can still see old messages (read-only)
   
3. **Can Message Check**
   - Check if user is blocked
   - Check recipient's `who_can_message` setting
   - Check if users are already in conversation (bypass privacy for existing conversations)

---

### Real-Time Delivery

1. **WebSocket Connection**
   - Authenticate via JWT token
   - Store connection in memory (user_id → connection map)
   - Heartbeat/ping every 30 seconds
   - Auto-disconnect after 60 seconds of inactivity

2. **Message Broadcasting**
   - When message sent, check if recipient is online
   - If online, send via WebSocket immediately
   - If offline, store for later delivery (database only)
   - Set message status to "delivered" when recipient receives

3. **Online Status**
   - User is "online" if WebSocket connected
   - Update `last_seen` timestamp on disconnect
   - Broadcast online/offline status to relevant conversations

---

## 🔒 Privacy & Security

### Security Measures

1. **Authentication & Authorization**
   - All endpoints require JWT authentication
   - Verify user owns conversation before access
   - Rate limiting on message sending (10 messages/minute per user)

2. **Input Validation**
   - Sanitize text content (prevent XSS)
   - Validate image file types and sizes
   - Limit message length (10,000 chars)
   - Validate UUIDs format

3. **SQL Injection Prevention**
   - Use parameterized queries (pgx with placeholders)
   - Never concatenate user input into SQL

4. **Image Security**
   - Validate MIME types server-side
   - Re-encode images to strip EXIF data
   - Limit file size (10MB)
   - Store with random UUID filenames

5. **WebSocket Security**
   - Authenticate WebSocket connections with JWT
   - Validate all incoming WebSocket messages
   - Implement rate limiting on WebSocket messages
   - Disconnect abusive connections

---

### Privacy Controls

1. **Message Privacy**
   - Messages only visible to participants
   - Deleted messages not recoverable (soft delete in DB)
   - No screenshots prevention (client-side, future)

2. **Blocking**
   - Blocked users cannot see online status
   - Blocked users cannot send messages
   - Blocking is one-way (A blocks B ≠ B blocks A)

3. **Read Receipts (Optional)**
   - User can disable read receipts in settings (future)
   - Default: enabled

---

## 📎 File Handling & Security

### Supported File Types & Limits

| Category | File Types | Max Size | Storage Path |
|----------|------------|----------|--------------|
| **Images** | JPEG, PNG, GIF, WebP | 10MB each | `/storage/messages/images/` |
| **Documents** | PDF, DOCX, XLSX, PPTX, TXT | 25MB each | `/storage/messages/files/` |
| **Archives** | ZIP, RAR, 7Z | 50MB each | `/storage/messages/files/` |
| **Data** | CSV, JSON, XML | 10MB each | `/storage/messages/files/` |

**Per-Message Limits:**
- Maximum 5 files per message
- Combined size limit: 100MB per message
- File name length: 255 characters max

---

### File Validation Strategy

#### 1. MIME Type Validation
```go
allowedMIMEs := map[string]bool{
    // Images
    "image/jpeg": true,
    "image/png": true,
    "image/gif": true,
    "image/webp": true,
    
    // Documents
    "application/pdf": true,
    "application/msword": true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    "application/vnd.ms-excel": true,
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
    "application/vnd.ms-powerpoint": true,
    "application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
    "text/plain": true,
    
    // Archives
    "application/zip": true,
    "application/x-rar-compressed": true,
    "application/x-7z-compressed": true,
    
    // Data
    "text/csv": true,
    "application/json": true,
    "application/xml": true,
    "text/xml": true,
}
```

#### 2. File Extension Blacklist
```go
// Reject dangerous file types
blacklistedExtensions := []string{
    ".exe", ".bat", ".cmd", ".sh", ".ps1",
    ".msi", ".jar", ".apk", ".app", ".dmg",
    ".scr", ".vbs", ".js", ".wsf", ".hta",
}
```

#### 3. Magic Bytes Verification (Optional)
Verify actual file type by reading file signature:
```go
// Example: Check if file is actually a JPEG
func isJPEG(data []byte) bool {
    return len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
}
```

#### 4. Filename Sanitization
```go
// Remove special characters, limit length
func sanitizeFilename(name string) string {
    // Remove path separators
    name = strings.ReplaceAll(name, "/", "_")
    name = strings.ReplaceAll(name, "\\", "_")
    
    // Limit length
    if len(name) > 255 {
        name = name[:255]
    }
    
    return name
}
```

---

### File Storage Strategy

#### Storage Structure
```
storage/messages/
├── images/
│   ├── 2024/11/uuid1.jpg          # Original image (resized)
│   ├── 2024/11/uuid2.png
│   └── thumbnails/
│       ├── uuid1.jpg              # Thumbnail (400x400)
│       └── uuid2.jpg
└── files/
    ├── 2024/11/uuid3.pdf
    ├── 2024/11/uuid4.docx
    └── 2024/11/uuid5.zip
```

#### File Naming Convention
- Use UUID v4 for unique filenames
- Preserve original extension
- Store original filename in database
- Format: `{uuid}.{extension}`

#### Image Processing
1. **Resize large images:**
   - Max dimensions: 1920x1920
   - Maintain aspect ratio
   - Compress to 85% quality
   - Strip EXIF metadata (privacy)

2. **Generate thumbnails:**
   - Fixed size: 400x400 (crop to square)
   - Used for message previews
   - Format: JPEG at 75% quality

#### Document Storage
- Store original file without modification
- No processing required for PDFs, DOCX, etc.
- Serve with proper `Content-Type` headers
- Set `Content-Disposition: attachment` for downloads

---

### File Download & Serving

#### Secure File Access
1. **Authentication required:** All file downloads require JWT token
2. **Permission check:** Verify user is conversation participant
3. **Signed URLs (Optional):** Generate temporary signed URLs (expire in 1 hour)
4. **Rate limiting:** Limit downloads to 100 per minute per user

#### Download Endpoint
```
GET /v1/messages/files/:fileId
Authorization: Bearer {token}

Response:
- Headers:
  - Content-Type: {mime_type}
  - Content-Disposition: attachment; filename="{original_filename}"
  - Content-Length: {size}
- Body: File bytes
```

---

### Security Best Practices

#### 1. Virus Scanning (Optional for MVP, Recommended for Production)
- Integrate with ClamAV or cloud service (VirusTotal API)
- Scan files on upload before saving
- Quarantine suspicious files
- Notify user if file is rejected

#### 2. Content Validation
- **PDFs:** Check for embedded JavaScript or malicious content
- **Office files:** Scan for macros
- **Archives:** Limit extraction depth (prevent zip bombs)

#### 3. Storage Security
- Files stored outside web root (not directly accessible)
- Serve files through authenticated endpoint
- Set proper file permissions (read-only for web server)
- Regular backups

#### 4. Bandwidth Protection
- Implement rate limiting on uploads
- Limit concurrent uploads per user (max 3)
- Throttle download speeds if needed
- Monitor storage usage per user

---

### File Metadata in Database

Stored in `messages.attachments` (JSONB):
```json
[
  {
    "id": "attachment-uuid",
    "type": "image",
    "url": "/uploads/messages/images/2024/11/uuid.jpg",
    "thumbnail_url": "/uploads/messages/images/thumbnails/uuid.jpg",
    "filename": "sunset_photo.jpg",
    "original_filename": "IMG_20241129_153045.jpg",
    "size": 2048576,
    "mime_type": "image/jpeg",
    "width": 1920,
    "height": 1080,
    "uploaded_at": "2024-11-29T15:30:45Z"
  },
  {
    "id": "attachment-uuid",
    "type": "file",
    "url": "/uploads/messages/files/2024/11/uuid.pdf",
    "filename": "assignment.pdf",
    "original_filename": "CS101_Final_Assignment.pdf",
    "size": 512000,
    "mime_type": "application/pdf",
    "uploaded_at": "2024-11-29T15:30:50Z"
  }
]
```

---

### Error Handling

| Error | HTTP Code | Message |
|-------|-----------|---------|
| File too large | 413 | "File exceeds maximum size limit (10MB for images, 25MB for documents)" |
| Invalid file type | 400 | "File type not supported. Allowed: JPEG, PNG, PDF, DOCX, etc." |
| Too many files | 400 | "Maximum 5 files per message" |
| Virus detected | 403 | "File rejected: Security scan failed" |
| Storage full | 507 | "Server storage limit reached. Please try again later." |

---

### Performance Optimization

1. **Async Processing (Future):**
   - Upload file immediately
   - Process (resize, thumbnail) in background
   - Update message with processed URLs
   - Send WebSocket notification when ready

2. **CDN Integration (Future):**
   - Serve files via CloudFront/Cloudflare
   - Reduce server bandwidth
   - Faster downloads globally
   - Automatic caching

3. **Compression:**
   - Compress images on upload (85% quality)
   - Optional: Compress documents with gzip
   - Store compressed, serve decompressed

4. **Cleanup Strategy:**
   - Delete files when message is hard-deleted
   - Delete orphaned files (no message reference)
   - Archive old files (>1 year) to cold storage

---

## ⚡ Real-Time Delivery Strategy

### Option 1: WebSocket (Recommended for MVP)

**Pros:**
- True real-time, low latency
- Bidirectional communication
- Standard technology, good library support

**Cons:**
- Requires persistent connections
- More server resources for connection management
- Complex scaling (sticky sessions or Redis pub/sub)

**Implementation:**
```
1. Client connects to WS endpoint with JWT
2. Server stores connection in memory map (user_id → ws.Conn)
3. When message sent:
   a. Save to database
   b. Check if recipient is online
   c. If online, send via WebSocket
   d. If offline, skip (will fetch on next login)
4. Client auto-reconnects on disconnect
```

**Libraries:**
- `gorilla/websocket` (popular Go WebSocket library)
- `nhooyr.io/websocket` (alternative, better API)

---

### Option 2: Server-Sent Events (SSE)

**Pros:**
- Simpler than WebSocket
- Built on HTTP, works with existing infrastructure
- Automatic reconnection

**Cons:**
- One-way only (server → client)
- Client must use polling for sending messages
- Not ideal for bidirectional chat

**Verdict:** Not recommended for chat (use WebSocket instead)

---

### Option 3: Long Polling

**Pros:**
- Works everywhere (no special protocols)
- Simple implementation
- Good fallback option

**Cons:**
- Higher latency
- More HTTP overhead
- Inefficient for high-volume messaging

**Verdict:** Use as fallback if WebSocket fails

---

### Recommended Approach: WebSocket + Database Fallback

**Strategy:**
1. **Primary:** WebSocket for real-time delivery
2. **Fallback:** Database polling for offline users
3. **Hybrid:** Check database on app open for missed messages

**Flow:**
```
User A sends message to User B:
├─ Save message to database
├─ Check if User B is online (WebSocket connected)
│  ├─ YES: Send via WebSocket immediately
│  └─ NO: Do nothing (User B fetches on next login)
└─ Return success response to User A

User B opens app:
├─ Connect to WebSocket
├─ Fetch unread conversations and messages
└─ Listen for new messages via WebSocket
```

---

### Scaling Considerations

For future scaling beyond MVP:

**Single Server (MVP):**
- In-memory map for connections: `map[userID]websocket.Conn`
- Works for < 10,000 concurrent connections

**Multi-Server (Future):**
- Use Redis Pub/Sub to broadcast messages across servers
- Store WebSocket connections in Redis
- Use sticky sessions (load balancer routes user to same server)

---

## 📁 File Structure

### Directory Structure
```
internal/app/
├── messaging/
│   ├── dto.go                    # Request/Response DTOs
│   ├── model.go                  # Domain models
│   ├── repository.go             # Repository interface
│   ├── repository_postgres.go    # PostgreSQL implementation
│   ├── service.go                # Business logic
│   ├── handler.go                # HTTP handlers
│   ├── routes.go                 # Route registration
│   ├── websocket.go              # WebSocket manager
│   └── privacy.go                # Privacy checker logic
│
internal/pkg/
├── websocket/
│   ├── manager.go                # WebSocket connection manager
│   ├── client.go                 # WebSocket client wrapper
│   └── hub.go                    # Message broadcasting hub
│
├── filehandler/
│   ├── validator.go              # File type & size validation
│   ├── processor.go              # Image processing, file handling
│   └── storage.go                # File storage operations
│
storage/
├── messages/                     # Message attachments
│   ├── images/                   # Image messages
│   │   ├── 2024/
│   │   │   └── 11/
│   │   │       └── [uuid].jpg
│   │   └── thumbnails/           # Image thumbnails (400x400)
│   │       └── [uuid].jpg
│   └── files/                    # Document/file messages
│       ├── 2024/
│       │   └── 11/
│       │       ├── [uuid].pdf
│       │       ├── [uuid].docx
│       │       └── [uuid].zip
```

---

### File Breakdown

#### `dto.go` (Request/Response DTOs)
```go
// Request DTOs
type CreateConversationRequest
type SendMessageRequest
type UploadFilesRequest

// Response DTOs
type ConversationResponse
type ConversationListResponse
type MessageResponse
type MessageListResponse
type UnreadCountResponse
type CanMessageResponse
type FileUploadResponse

// Attachment DTOs
type AttachmentMetadata struct {
    ID               string    `json:"id"`
    Type             string    `json:"type"` // "image" or "file"
    URL              string    `json:"url"`
    ThumbnailURL     string    `json:"thumbnail_url,omitempty"` // For images only
    Filename         string    `json:"filename"`
    OriginalFilename string    `json:"original_filename"`
    Size             int64     `json:"size"`
    MimeType         string    `json:"mime_type"`
    Width            int       `json:"width,omitempty"`  // For images only
    Height           int       `json:"height,omitempty"` // For images only
    UploadedAt       time.Time `json:"uploaded_at"`
}
```

#### `model.go` (Domain Models)
```go
type Conversation struct {
    ID                   string
    User1ID              string
    User2ID              string
    LastMessageID        *string
    LastMessageText      *string
    LastMessageSenderID  *string
    LastMessageAt        *time.Time
    CreatedAt            time.Time
    UpdatedAt            time.Time
}

type Message struct {
    ID             string
    ConversationID string
    SenderID       string
    MessageType    string // "text", "image", "file"
    Content        *string
    Attachments    []AttachmentMetadata // Stored as JSONB in DB
    Status         string
    IsRead         bool
    ReadAt         *time.Time
    IsDeleted      bool
    DeletedAt      *time.Time
    DeletedBy      *string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type BlockedUser struct {
    ID        string
    BlockerID string
    BlockedID string
    Reason    *string
    CreatedAt time.Time
}
```

#### `repository.go` (Interface)
```go
type ConversationRepository interface
type MessageRepository interface
type BlockRepository interface
```

#### `repository_postgres.go` (Implementation)
```go
// PostgreSQL implementations of repositories
```

#### `service.go` (Business Logic)
```go
type MessagingService interface {
    // Conversations
    CreateOrGetConversation(ctx, userID, recipientID)
    GetConversations(ctx, userID, page, limit)
    GetConversation(ctx, conversationID, userID)
    DeleteConversation(ctx, conversationID, userID)
    
    // Messages
    SendMessage(ctx, conversationID, senderID, req)
    GetMessages(ctx, conversationID, userID, page, limit)
    DeleteMessage(ctx, messageID, userID)
    
    // Unread
    MarkConversationAsRead(ctx, conversationID, userID)
    MarkMessageAsRead(ctx, messageID, userID)
    GetUnreadCount(ctx, userID)
    
    // Privacy
    CanMessage(ctx, senderID, recipientID)
    BlockUser(ctx, blockerID, blockedID)
    UnblockUser(ctx, blockerID, blockedID)
    GetBlockedUsers(ctx, userID)
}
```

#### `handler.go` (HTTP Handlers)
```go
func (h *Handler) CreateConversation(c *fiber.Ctx)
func (h *Handler) GetConversations(c *fiber.Ctx)
func (h *Handler) GetConversation(c *fiber.Ctx)
func (h *Handler) DeleteConversation(c *fiber.Ctx)
func (h *Handler) SendMessage(c *fiber.Ctx)
func (h *Handler) GetMessages(c *fiber.Ctx)
func (h *Handler) DeleteMessage(c *fiber.Ctx)
func (h *Handler) MarkConversationAsRead(c *fiber.Ctx)
func (h *Handler) MarkMessageAsRead(c *fiber.Ctx)
func (h *Handler) GetUnreadCount(c *fiber.Ctx)
func (h *Handler) CanMessage(c *fiber.Ctx)
func (h *Handler) BlockUser(c *fiber.Ctx)
func (h *Handler) UnblockUser(c *fiber.Ctx)
func (h *Handler) GetBlockedUsers(c *fiber.Ctx)
func (h *Handler) UploadFiles(c *fiber.Ctx)
func (h *Handler) DownloadFile(c *fiber.Ctx)
func (h *Handler) HandleWebSocket(c *fiber.Ctx)
```

#### `routes.go` (Route Registration)
```go
func RegisterRoutes(app *fiber.App, handler *Handler)
```

#### `websocket.go` (WebSocket Manager)
```go
type WSManager struct
func NewWSManager() *WSManager
func (m *WSManager) HandleConnection(conn *websocket.Conn, userID string)
func (m *WSManager) SendToUser(userID string, message interface{})
func (m *WSManager) BroadcastToConversation(conversationID string, message interface{})
func (m *WSManager) IsUserOnline(userID string) bool
```

#### `privacy.go` (Privacy Checker)
```go
type PrivacyChecker struct
func NewPrivacyChecker(userRepo, blockRepo) *PrivacyChecker
func (p *PrivacyChecker) CanMessage(ctx, senderID, recipientID) (bool, string)
func (p *PrivacyChecker) IsBlocked(ctx, user1ID, user2ID) bool
```

---

## 📅 Implementation Phases

### Phase 1: Database & Models (Week 1, Days 1-2)
**Goal:** Set up database schema and domain models

**Tasks:**
1. Create migration files (3 tables)
   - `000015_create_conversations`
   - `000016_create_messages`
   - `000017_create_blocked_users`

2. Run migrations and verify schema
3. Create domain models in `model.go`
4. Create DTOs in `dto.go`

**Deliverable:** Database tables created, models defined

---

### Phase 2: Repository Layer (Week 1, Days 3-4)
**Goal:** Implement data access layer

**Tasks:**
1. Define repository interfaces
2. Implement PostgreSQL repositories:
   - ConversationRepository (CRUD operations)
   - MessageRepository (CRUD operations)
   - BlockRepository (CRUD operations)

3. Write SQL queries with proper indexes
4. Handle edge cases (duplicate conversations, etc.)

**Deliverable:** Repository layer complete with all database operations

---

### Phase 3: Service Layer - Core Logic (Week 1, Days 5-7)
**Goal:** Implement business logic

**Tasks:**
1. Create `MessagingService` interface
2. Implement conversation management:
   - Create/get conversation (with deduplication)
   - List conversations (paginated, sorted)
   - Delete conversation (soft delete)

3. Implement message management:
   - Send message (with validation)
   - Get messages (paginated)
   - Delete message

4. Implement privacy checker:
   - Check `who_can_message` settings
   - Check blocking status
   - Validate message permissions

**Deliverable:** Complete service layer with business logic

---

### Phase 4: HTTP Handlers & Routes (Week 2, Days 1-2)
**Goal:** Create REST API endpoints

**Tasks:**
1. Create HTTP handlers for all endpoints
2. Request validation and error handling
3. Register routes in `routes.go`
4. Integrate with existing auth middleware
5. Test with Postman

**Deliverable:** 15 REST API endpoints working

---

### Phase 5: Unread Tracking (Week 2, Days 3-4)
**Goal:** Implement unread message counting

**Tasks:**
1. Create unread count queries (optimized)
2. Implement mark-as-read endpoints
3. Auto-mark read when fetching messages
4. Add unread count to conversation responses
5. Test edge cases (multiple devices, race conditions)

**Deliverable:** Unread tracking fully functional

---

### Phase 6: File Upload (Week 2, Days 5-6)
**Goal:** Support image and file messages

**Tasks:**
1. Create file upload endpoint (`POST /v1/messages/upload`)
2. Implement file type validation:
   - Whitelist allowed MIME types
   - Blacklist dangerous extensions (.exe, .sh, .bat)
   - Verify file signatures (magic bytes)
3. **Image handling:**
   - Resize images (max 1920x1920)
   - Generate thumbnails (400x400)
   - Store in `/storage/messages/images/`
   - Supported formats: JPEG, PNG, GIF, WebP
4. **Document handling:**
   - Store original files (no processing)
   - Store in `/storage/messages/files/`
   - Supported formats: PDF, DOCX, XLSX, PPTX, TXT, CSV, JSON, XML
5. **Archive handling:**
   - Store in `/storage/messages/files/`
   - Supported formats: ZIP, RAR, 7Z
   - Optional: Virus scanning for archives
6. Implement size limits per file type
7. Support multiple file uploads (up to 5 per message)
8. Return file metadata (URL, size, MIME type, filename)
9. Update message sending to support attachments

**Deliverable:** Image and file messages working

---

### Phase 7: WebSocket - Basic (Week 3, Days 1-3)
**Goal:** Implement real-time message delivery

**Tasks:**
1. Set up WebSocket endpoint
2. Implement connection authentication (JWT)
3. Create WebSocket manager (in-memory connection map)
4. Handle connection lifecycle (connect/disconnect)
5. Implement message broadcasting
6. Test real-time delivery

**Deliverable:** Basic WebSocket working for real-time messages

---

### Phase 8: WebSocket - Advanced (Week 3, Days 4-5)
**Goal:** Add advanced real-time features

**Tasks:**
1. Implement typing indicators
2. Implement online/offline status
3. Implement read receipts via WebSocket
4. Handle reconnection logic
5. Add heartbeat/ping mechanism

**Deliverable:** Advanced WebSocket features working

---

### Phase 9: Testing & Bug Fixes (Week 3, Days 6-7)
**Goal:** Comprehensive testing

**Tasks:**
1. Test all API endpoints with Postman
2. Test WebSocket scenarios:
   - Both users online
   - One user offline
   - Network interruptions
3. Test privacy scenarios (blocking, settings)
4. Test edge cases:
   - Empty conversations
   - Large message history
   - Concurrent message sending
5. Fix bugs and edge cases
6. Performance testing (pagination, queries)

**Deliverable:** Fully tested and bug-free messaging system

---

### Phase 10: Documentation (Week 4, Days 1-2)
**Goal:** Complete documentation

**Tasks:**
1. Create Postman collection for messaging APIs
2. Update progress report
3. Write implementation notes
4. Document WebSocket protocol
5. Add inline code comments
6. Create testing guide

**Deliverable:** Complete documentation

---

### Phase 11: Polish & Optimization (Week 4, Days 3-5)
**Goal:** Performance and UX improvements

**Tasks:**
1. Optimize database queries (add indexes if needed)
2. Implement cursor-based pagination for messages
3. Add rate limiting on message sending
4. Optimize WebSocket memory usage
5. Add logging and monitoring
6. Stress testing (100+ concurrent users)

**Deliverable:** Production-ready messaging system

---

### Phase 12: Future Features (Optional, Week 4+)
**Goal:** Nice-to-have features

**Tasks:**
1. Message editing (within 15 minutes)
2. Delete for everyone (within 1 hour)
3. Message reactions (emoji)
4. Forward messages
5. Message search
6. Voice messages (audio)
7. Push notifications (FCM/APNs)

**Deliverable:** Enhanced messaging experience

---

## 🧪 Testing Strategy

### 1. Unit Tests (Optional for MVP)
- Test service layer logic
- Test privacy checker
- Test repository queries
- Mock database connections

### 2. Integration Tests
- Test API endpoints end-to-end
- Test with real database
- Test WebSocket connections
- Test concurrent users

### 3. Manual Testing with Postman

**Test Cases:**

#### Conversations
- [ ] Create conversation with another user
- [ ] Get conversation list (paginated)
- [ ] Get conversation by ID
- [ ] Delete conversation
- [ ] Try to access other user's conversation (should fail)

#### Messages
- [ ] Send text message
- [ ] Send image message (single image)
- [ ] Send multiple images (up to 5)
- [ ] Send PDF file
- [ ] Send DOCX file
- [ ] Send ZIP file
- [ ] Send mixed attachments (images + files)
- [ ] Get message history (paginated)
- [ ] Download attached file
- [ ] Delete own message
- [ ] Try to delete other user's message (should fail)
- [ ] Send message to user who blocked you (should fail)

#### Unread Tracking
- [ ] Send message and verify unread count increases
- [ ] Mark conversation as read
- [ ] Mark single message as read
- [ ] Get total unread count

#### Privacy & Blocking
- [ ] Block a user
- [ ] Verify blocked user cannot send messages
- [ ] Unblock a user
- [ ] Check `who_can_message` settings respected
- [ ] Test `can-message` endpoint

#### WebSocket
- [ ] Connect to WebSocket
- [ ] Send message and receive in real-time
- [ ] Test typing indicators
- [ ] Test online/offline status
- [ ] Test read receipts
- [ ] Disconnect and reconnect

#### Edge Cases
- [ ] Send message to self (should fail)
- [ ] Send very long message (10,000+ chars)
- [ ] Upload large image (>10MB, should fail)
- [ ] Upload large document (>25MB, should fail)
- [ ] Upload large archive (>50MB, should fail)
- [ ] Upload invalid file type (.exe, should fail)
- [ ] Upload 6 files (should fail, max is 5)
- [ ] Upload file with malicious name (../../../etc/passwd)
- [ ] Upload file with very long filename (>255 chars)
- [ ] Upload empty file (0 bytes)
- [ ] Rapid message sending (rate limiting)
- [ ] Rapid file uploads (rate limiting)
- [ ] Paginate through 1000+ messages
- [ ] Download file without authentication (should fail)

---

### 4. Load Testing (Optional)

**Tools:** Apache Bench (ab), k6, or Artillery

**Scenarios:**
- 100 concurrent users sending messages
- 1000 concurrent WebSocket connections
- Fetch message history for 100 conversations simultaneously

**Metrics:**
- Response time (p50, p95, p99)
- Error rate
- Throughput (messages/second)
- WebSocket connection stability

---

## ⚡ Performance Considerations

### Database Optimization

1. **Indexes**
   - Composite index on `(conversation_id, created_at DESC)` for message queries
   - Index on `user1_id` and `user2_id` for conversation lookups
   - Partial index on `is_read = FALSE` for unread counts

2. **Query Optimization**
   - Use `LIMIT` and `OFFSET` for pagination (or cursor-based)
   - Avoid N+1 queries (join users table for sender info)
   - Use `COUNT(*)` with indexes for fast counts

3. **Denormalization**
   - Store last message info in `conversations` table
   - Avoids expensive JOIN on messages table for inbox listing

4. **Connection Pooling**
   - Reuse existing pgxpool (already configured)
   - Max connections: 25 (adjust based on load)

---

### WebSocket Optimization

1. **Memory Management**
   - Limit max concurrent connections per server (10,000)
   - Clean up disconnected connections immediately
   - Use sync.Map for thread-safe connection storage

2. **Message Broadcasting**
   - Only send to online users
   - Use goroutines for parallel message sending
   - Buffer messages if sending fails (retry queue)

3. **Heartbeat**
   - Ping every 30 seconds
   - Disconnect after 60 seconds of inactivity
   - Client auto-reconnects

---

### Image Storage

1. **Optimization**
   - Resize images to max 1920x1920 (reduce bandwidth)
   - Generate thumbnails (400x400) for previews
   - Use WebP format for better compression

2. **CDN (Future)**
   - Serve images via CDN (CloudFront, Cloudflare)
   - Reduce server load
   - Faster image delivery

---

### Caching (Future)

1. **Redis Cache**
   - Cache unread counts (TTL: 5 minutes)
   - Cache conversation list (TTL: 1 minute)
   - Cache blocked users list (TTL: 10 minutes)

2. **In-Memory Cache**
   - Cache online user status (already in WebSocket manager)
   - Cache privacy settings (reload on change)

---

## 🚀 Future Enhancements

### Phase 2 Features (Post-MVP)

1. **Group Messaging**
   - Create group conversations
   - Add/remove participants
   - Group admin controls
   - Use `conversation_participants` table

2. **Message Features**
   - Edit messages (within 15 minutes)
   - Delete for everyone (within 1 hour)
   - Forward messages
   - Reply to specific message
   - Message reactions (emoji)

3. **Rich Media (Advanced)**
   - Voice messages (audio)
   - Video messages (short clips)
   - Location sharing
   - Contact sharing

4. **Search**
   - Search messages within conversation
   - Full-text search with PostgreSQL FTS
   - Search across all conversations

5. **Notifications**
   - Push notifications (FCM for Android, APNs for iOS)
   - Email notifications (optional)
   - Notification preferences

6. **Privacy**
   - Read receipts toggle (hide read status)
   - Last seen privacy (hide online status)
   - Disappearing messages (auto-delete after X time)

7. **Performance**
   - Redis pub/sub for multi-server WebSocket
   - Message queue for async processing
   - CDN for image delivery

8. **Analytics**
   - Message delivery success rate
   - Average response time
   - Active conversation count

---

## 📊 Success Criteria

### Functional Requirements
- ✅ Users can send and receive text messages
- ✅ Users can send and receive image messages
- ✅ Users can see message history
- ✅ Users can see unread message counts
- ✅ Users can block/unblock other users
- ✅ Privacy settings are respected
- ✅ Real-time message delivery via WebSocket

### Non-Functional Requirements
- ✅ Message delivery latency < 1 second
- ✅ Support 10,000 concurrent connections
- ✅ API response time < 500ms (p95)
- ✅ 99.9% uptime
- ✅ No data loss (all messages persisted)

### User Experience
- ✅ Smooth real-time chat experience
- ✅ Fast message loading (pagination)
- ✅ Accurate unread counts
- ✅ Reliable image upload

---

## 📝 Notes & Considerations

### Architecture Decisions

1. **Why 1-on-1 only for MVP?**
   - Simpler implementation
   - Covers 80% of use cases
   - Group chat adds complexity (permissions, notifications)
   - Can be added later using `conversation_participants` table

2. **Why WebSocket over polling?**
   - True real-time experience
   - Lower latency
   - Less server load (no constant HTTP requests)
   - Industry standard for chat

3. **Why soft delete?**
   - Data integrity
   - Audit trail
   - "Delete for me" functionality
   - Can add "Delete for everyone" later

4. **Why denormalize last message?**
   - Faster inbox listing (no JOIN required)
   - Trade-off: slight complexity in update logic
   - Worth it for performance

5. **Why cursor-based pagination for messages?**
   - More efficient than offset pagination
   - No "missing messages" issue when new messages arrive
   - Better for real-time chat

---

### Challenges & Solutions

1. **Challenge:** Duplicate conversations
   **Solution:** Use `user1_id < user2_id` ordering + unique constraint

2. **Challenge:** Race conditions in unread count
   **Solution:** Use database transactions for atomic updates

3. **Challenge:** WebSocket scaling across multiple servers
   **Solution:** (Future) Use Redis pub/sub for message broadcasting

4. **Challenge:** Large message history loading slowly
   **Solution:** Cursor-based pagination + limit to 50 messages per page

5. **Challenge:** Image storage filling up disk
   **Solution:** (Future) Move to S3, implement cleanup for old images

---

## 🎓 Learning Resources

### WebSocket in Go
- Gorilla WebSocket: https://github.com/gorilla/websocket
- nhooyr.io WebSocket: https://github.com/nhooyr/websocket
- Tutorial: https://yalantis.com/blog/how-to-build-websockets-in-go/

### Real-Time Chat Architecture
- Building a Chat App: https://www.cometchat.com/tutorials/building-a-chat-application
- Scaling WebSockets: https://blog.container-solutions.com/scaling-websockets

### Database Design for Messaging
- Chat Schema Design: https://stackoverflow.com/questions/5731923/database-schema-for-chat-application
- Optimizing Queries: https://use-the-index-luke.com/

---

## ✅ Pre-Implementation Checklist

Before starting implementation, ensure:

- [ ] Read and understand entire plan
- [ ] Review existing codebase structure
- [ ] Set up development environment
- [ ] Install WebSocket library (`go get github.com/gorilla/websocket`)
- [ ] Create feature branch (`git checkout -b feature/direct-messaging`)
- [ ] Back up database before migrations
- [ ] Plan testing strategy
- [ ] Allocate sufficient time (realistic estimate)

---

## 📌 Summary

This plan provides a comprehensive roadmap for implementing the Direct Messaging feature in Mockhu. It covers:

- ✅ Complete database schema (3 tables with JSONB for attachments)
- ✅ 16 REST API endpoints
- ✅ WebSocket for real-time delivery
- ✅ Privacy controls and blocking
- ✅ Unread tracking
- ✅ **Image messaging (up to 5 images per message)**
- ✅ **File attachments (PDF, DOCX, XLSX, PPTX, ZIP, and more)**
- ✅ **Comprehensive file validation and security**
- ✅ **Support for multiple file types with proper size limits**
- ✅ Detailed implementation phases
- ✅ Testing strategy
- ✅ Performance considerations
- ✅ Future enhancements

**Estimated Timeline:** 3-4 weeks for full implementation

**Next Steps:**
1. Review and approve this plan
2. Create feature branch
3. Start with Phase 1 (Database & Models)
4. Implement incrementally, testing each phase
5. Document progress and learnings

---

**Good luck with the implementation! 🚀**

**Questions or need clarification? Feel free to ask!**

