# Frontend Implementation Prompt: Postfeed Feature

## Overview
Implement a social media feed feature called **Postfeed** that supports 4 types of posts: **Doubt**, **Quiz**, **Progress**, and **Resource**. The feature includes likes, comments with replies, shares, and media attachments.

---

## Backend API Base URL
```
http://localhost:8085
```

---

## Authentication
All protected endpoints require a JWT token in the `Authorization` header:
```
Authorization: Bearer <token>
```

---

## Data Models

### Postfeed Response
```typescript
interface Postfeed {
  id: string;
  user_id: string;
  author?: AuthorInfo;  // null if is_anonymous=true
  type: 'doubt' | 'quiz' | 'progress' | 'resource';
  title: string;
  content?: string;
  tags?: string[];
  media?: MediaItem[];
  visibility: 'public' | 'private' | 'followers_only';
  is_anonymous: boolean;
  metadata?: DoubtMetadata | QuizMetadata | ProgressMetadata | ResourceMetadata;
  view_count: number;
  like_count: number;
  comment_count: number;
  share_count: number;
  created_at: string;  // ISO 8601
  updated_at: string;
}

interface AuthorInfo {
  id: string;
  username: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
}

interface MediaItem {
  url: string;
  type: 'image' | 'video' | 'audio' | 'document';
  thumbnail_url?: string;
  width?: number;
  height?: number;
  duration?: number;  // seconds for video/audio
  file_name?: string;
  file_size?: number;
}
```

### Type-Specific Metadata

```typescript
// For type="doubt"
interface DoubtMetadata {
  subject?: string;
  is_solved: boolean;
  best_answer_id?: string;
}

// For type="quiz"
interface QuizMetadata {
  questions: QuizQuestion[];
  time_limit_seconds?: number;
  difficulty?: 'easy' | 'medium' | 'hard';
}

interface QuizQuestion {
  question: string;
  options: string[];
  correct_index: number;
}

// For type="progress"
interface ProgressMetadata {
  milestone?: string;
  percentage?: number;
  streak_days?: number;
  badge_earned?: string;
}

// For type="resource"
interface ResourceMetadata {
  resource_type?: 'video' | 'pdf' | 'link' | 'article';
  url?: string;
  file_url?: string;
  platform?: string;  // youtube, notion, etc.
  duration_minutes?: number;
}
```

---

## API Endpoints

### Postfeeds CRUD

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/v1/postfeeds` | ✅ | Create postfeed |
| `GET` | `/v1/postfeeds` | ❌ | List postfeeds (supports `?type=`, `?tag=`, `?page=`, `?limit=`) |
| `GET` | `/v1/postfeeds/:id` | ❌ | Get single postfeed |
| `PUT` | `/v1/postfeeds/:id` | ✅ | Update own postfeed |
| `DELETE` | `/v1/postfeeds/:id` | ✅ | Delete own postfeed |
| `GET` | `/v1/postfeeds/type/:type` | ❌ | Filter by type (doubt/quiz/progress/resource) |
| `GET` | `/v1/users/:user_id/postfeeds` | ❌ | Get user's postfeeds |

### Likes

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/v1/postfeeds/:id/like` | ❌ | Check like status `{liked: bool, like_count: int}` |
| `POST` | `/v1/postfeeds/:id/like` | ✅ | Like a postfeed |
| `DELETE` | `/v1/postfeeds/:id/like` | ✅ | Unlike a postfeed |

### Comments

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/v1/postfeeds/:id/comments` | ❌ | List comments (includes nested replies) |
| `POST` | `/v1/postfeeds/:id/comments` | ✅ | Add comment `{content, parent_id?}` |
| `PUT` | `/v1/postfeeds/:id/comments/:comment_id` | ✅ | Update own comment |
| `DELETE` | `/v1/postfeeds/:id/comments/:comment_id` | ✅ | Delete own comment |

### Shares

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/v1/postfeeds/:id/shares` | ❌ | Get share count |
| `POST` | `/v1/postfeeds/:id/share` | ✅ | Share postfeed `{message?}` |

### Media Upload

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/v1/upload/media/request` | ✅ | Get presigned URL for R2 upload |

---

## UI Components to Build

### 1. Feed Page (`/feed`)
- Tab navigation: **All**, **Doubts**, **Quizzes**, **Progress**, **Resources**
- Infinite scroll or pagination
- Filter by tags
- Each post card shows:
  - Author avatar, name (or "Anonymous")
  - Post type badge
  - Title, content preview
  - Media thumbnails (if any)
  - Engagement bar: ❤️ likes, 💬 comments, 🔗 shares
  - Relative timestamp

### 2. Create Post Modal
- Type selector (4 types with icons)
- Common fields: title, content, tags, visibility, anonymous toggle
- Media upload with drag-drop and preview
- Type-specific fields:
  - **Doubt**: Subject dropdown
  - **Quiz**: Question builder (add/remove questions, options)
  - **Progress**: Milestone, percentage slider, streak days
  - **Resource**: URL input, resource type, platform

### 3. Post Detail Page (`/post/:id`)
- Full post content with media gallery
- Like button with animation
- Comments section with:
  - Nested replies (1 level)
  - Add comment form
  - Edit/delete own comments
- Share button

### 4. Post Type Cards
Design unique card styles for each type:
- **Doubt**: Question mark icon, "Solved" badge
- **Quiz**: Play button, question count, difficulty badge
- **Progress**: Celebration animation, progress bar
- **Resource**: Platform icon, thumbnail

---

## Media Upload Flow

```
1. User selects file
2. POST /v1/upload/media/request with {content_type, file_size}
3. Response: {upload_url, public_url}
4. PUT file to upload_url (direct to R2)
5. Add {url: public_url, type: 'image'} to media array
6. Submit postfeed with media
```

---

## State Management (Zustand/Redux)

```typescript
interface PostfeedStore {
  postfeeds: Postfeed[];
  loading: boolean;
  error: string | null;
  page: number;
  hasMore: boolean;
  
  // Actions
  fetchPostfeeds: (filters?: {type?, tag?, page?}) => Promise<void>;
  createPostfeed: (data: CreatePostfeedRequest) => Promise<Postfeed>;
  likePostfeed: (id: string) => Promise<void>;
  unlikePostfeed: (id: string) => Promise<void>;
  addComment: (postfeedId: string, content: string, parentId?: string) => Promise<void>;
}
```

---

## API Service Example

```typescript
// services/postfeedService.ts
import api from './api';

export const postfeedService = {
  list: (params?: {type?: string; tag?: string; page?: number; limit?: number}) =>
    api.get('/v1/postfeeds', { params }),
  
  getById: (id: string) =>
    api.get(`/v1/postfeeds/${id}`),
  
  create: (data: CreatePostfeedRequest) =>
    api.post('/v1/postfeeds', data),
  
  update: (id: string, data: UpdatePostfeedRequest) =>
    api.put(`/v1/postfeeds/${id}`, data),
  
  delete: (id: string) =>
    api.delete(`/v1/postfeeds/${id}`),
  
  like: (id: string) =>
    api.post(`/v1/postfeeds/${id}/like`),
  
  unlike: (id: string) =>
    api.delete(`/v1/postfeeds/${id}/like`),
  
  getComments: (id: string, page = 1, limit = 20) =>
    api.get(`/v1/postfeeds/${id}/comments`, { params: { page, limit } }),
  
  addComment: (id: string, content: string, parentId?: string) =>
    api.post(`/v1/postfeeds/${id}/comments`, { content, parent_id: parentId }),
  
  share: (id: string, message?: string) =>
    api.post(`/v1/postfeeds/${id}/share`, { message }),
  
  requestMediaUpload: (contentType: string, fileSize: number) =>
    api.post('/v1/upload/media/request', { content_type: contentType, file_size: fileSize }),
};
```

---

## Design Guidelines

- Use vibrant colors for type badges (e.g., purple for Doubt, green for Quiz)
- Glassmorphism cards with subtle shadows
- Smooth animations for likes (heart pulse), comments slide-in
- Skeleton loaders during fetch
- Empty states with illustrations
- Mobile-first responsive design

---

## Priority Order

1. Feed page with list/filter
2. Post detail with comments
3. Create post modal
4. Like/unlike functionality
5. Comment CRUD
6. Media upload
7. Share functionality
