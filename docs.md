# Pitstop API Documentation for Frontend Integration

This document provides comprehensive endpoint documentation for integrating with the Pitstop backend API.

## Base URL
```
http://localhost:8080/api/v1
```

## Response Format
All API endpoints return responses in the following structured format:

```json
{
  "success": true/false,
  "message": "Human-readable message",
  "data": { ... },           // Present on success
  "error": { ... },          // Present on failure
  "meta": { ... },           // Pagination info (when applicable)
  "timestamp": "2023-12-01T10:30:00Z"
}
```

---

## Authentication Endpoints

### 1. Initiate Google OAuth
Get the Google OAuth authorization URL to redirect users for login.

**Endpoint:** `GET /auth/google`

**Request:**
```http
GET /api/v1/auth/google
```

**Response:**
```json
{
  "success": true,
  "message": "Google OAuth URL generated successfully",
  "data": {
    "auth_url": "https://accounts.google.com/o/oauth2/auth?client_id=...&state=..."
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const response = await fetch('/api/v1/auth/google');
const result = await response.json();
if (result.success) {
  window.location.href = result.data.auth_url;
}
```

---

### 2. Google OAuth Callback
Handles the OAuth callback after user authorization (used internally by the backend).

**Endpoint:** `GET /auth/google/callback`

**Note:** This endpoint is called directly by Google and handles the OAuth callback flow. Your frontend doesn't need to call this directly.

---

### 3. Exchange Authorization Code for Tokens
Exchange the OAuth authorization code for JWT access and refresh tokens.

**Endpoint:** `POST /auth/exchange`

**Request Body:**
```json
{
  "code": "authorization_code_from_oauth_callback",
  "state": "state_token_from_oauth_callback"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Authorization code exchanged successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_at": 1640995200
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const exchangeTokens = async (code, state) => {
  const response = await fetch('/api/v1/auth/exchange', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ code, state }),
  });
  
  const result = await response.json();
  if (result.success) {
    localStorage.setItem('access_token', result.data.access_token);
    localStorage.setItem('refresh_token', result.data.refresh_token);
  }
};
```

---

### 4. Refresh Access Token
Get a new access token using the refresh token when the current token expires.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Tokens refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_at": 1640997000
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const refreshToken = async () => {
  const refreshToken = localStorage.getItem('refresh_token');
  const response = await fetch('/api/v1/auth/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  
  const result = await response.json();
  if (result.success) {
    localStorage.setItem('access_token', result.data.access_token);
    localStorage.setItem('refresh_token', result.data.refresh_token);
    return true;
  }
  return false;
};
```

---

### 5. Get Current User Info
Retrieve the authenticated user's information.

**Endpoint:** `GET /auth/me`
**Authentication:** Required (Bearer token)

**Request:**
```http
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "User information retrieved successfully",
  "data": {
    "id": "user-uuid",
    "provider_id": "google_id_123",
    "provider": "google",
    "first_name": "John",
    "last_name": "Doe",
    "username": "john_doe_123",
    "email": "user@example.com",
    "display_name": "John Doe",
    "bio": "Software Developer and car enthusiast",
    "avatar_url": "https://lh3.googleusercontent.com/a/...",
    "locale": "en",
    "created_at": "2023-12-01T10:30:00Z",
    "updated_at": "2023-12-01T10:30:00Z"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const getCurrentUser = async () => {
  const token = localStorage.getItem('access_token');
  const response = await fetch('/api/v1/auth/me', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to get user info');
};
```

---

## Authentication Flow Example

```javascript
// 1. Redirect to Google OAuth
const startAuth = async () => {
  const response = await fetch('/api/v1/auth/google');
  const result = await response.json();
  if (result.success) {
    window.location.href = result.data.auth_url;
  }
};

// 2. Handle callback (in your OAuth callback page)
const handleCallback = async () => {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  const state = urlParams.get('state');
  
  if (code && state) {
    const response = await fetch('/api/v1/auth/exchange', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code, state }),
    });
    
    const result = await response.json();
    if (result.success) {
      localStorage.setItem('access_token', result.data.access_token);
      localStorage.setItem('refresh_token', result.data.refresh_token);
      
      // Get user info and redirect to app
      const user = await getCurrentUser();
      window.location.href = '/dashboard';
    }
  }
};
```

---

## Error Responses

### Common Authentication Errors

**401 - Unauthorized:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**401 - Invalid Token:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired token",
    "details": "token has expired"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**400 - Validation Error:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Code and state are required",
    "details": "Both 'code' and 'state' fields must be provided"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

---

## Posts Endpoints

### 1. Get All Posts
Retrieve a paginated list of all posts with user information and engagement metrics.

**Endpoint:** `GET /posts`
**Authentication:** Not required (Public)

**Query Parameters:**
- `page` (optional): Page number, default is 1
- `limit` (optional): Number of posts per page, default is 20, max is 100

**Request:**
```http
GET /api/v1/posts?page=1&limit=20
```

**Response:**
```json
{
  "success": true,
  "message": "Posts retrieved successfully",
  "data": {
    "posts": [
      {
        "id": "post-uuid-123",
        "user_id": "user-uuid-456",
        "content": "Just got my new BMW M3! Can't wait to take it for a spin 🚗",
        "user": {
          "username": "john_doe_123",
          "display_name": "John Doe",
          "avatar_url": "https://lh3.googleusercontent.com/a/..."
        },
        "comment_count": 15,
        "like_count": 42,
        "created_at": "2023-12-01T10:30:00Z",
        "updated_at": "2023-12-01T10:30:00Z"
      }
    ],
    "total_count": 150,
    "page": 1,
    "limit": 20,
    "has_next": true
  },
  "meta": {
    "pagination": {
      "current_page": 1,
      "total_pages": 8,
      "total_items": 150,
      "items_per_page": 20,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const getPosts = async (page = 1, limit = 20) => {
  const response = await fetch(`/api/v1/posts?page=${page}&limit=${limit}`);
  const result = await response.json();
  
  if (result.success) {
    return {
      posts: result.data.posts,
      pagination: result.meta.pagination
    };
  }
  throw new Error(result.error?.message || 'Failed to fetch posts');
};
```

---

### 2. Get Single Post
Retrieve a specific post by ID with user information and engagement metrics.

**Endpoint:** `GET /posts/{id}`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/posts/post-uuid-123
```

**Response:**
```json
{
  "success": true,
  "message": "Post retrieved successfully",
  "data": {
    "id": "post-uuid-123",
    "user_id": "user-uuid-456",
    "content": "Just got my new BMW M3! Can't wait to take it for a spin 🚗",
    "user": {
      "username": "john_doe_123",
      "display_name": "John Doe",
      "avatar_url": "https://lh3.googleusercontent.com/a/..."
    },
    "comment_count": 15,
    "like_count": 42,
    "created_at": "2023-12-01T10:30:00Z",
    "updated_at": "2023-12-01T10:30:00Z"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const getPost = async (postId) => {
  const response = await fetch(`/api/v1/posts/${postId}`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch post');
};
```

---

### 3. Create Post
Create a new post.

**Endpoint:** `POST /posts`
**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "user_id": "user-uuid-456",
  "content": "Just picked up my dream car! A 1967 Ford Mustang Fastback in pristine condition."
}
```

**Request:**
```http
POST /api/v1/posts
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "user_id": "user-uuid-456",
  "content": "Just picked up my dream car! A 1967 Ford Mustang Fastback in pristine condition."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Post created successfully",
  "data": {
    "id": "post-uuid-789",
    "user_id": "user-uuid-456",
    "content": "Just picked up my dream car! A 1967 Ford Mustang Fastback in pristine condition.",
    "user": {
      "username": "john_doe_123",
      "display_name": "John Doe",
      "avatar_url": "https://lh3.googleusercontent.com/a/..."
    },
    "comment_count": 0,
    "like_count": 0,
    "created_at": "2023-12-01T15:45:00Z",
    "updated_at": "2023-12-01T15:45:00Z"
  },
  "timestamp": "2023-12-01T15:45:00Z"
}
```

**Frontend Usage:**
```javascript
const createPost = async (userId, content) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch('/api/v1/posts', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userId,
      content: content,
    }),
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to create post');
};
```

---

## Comments Endpoints

### 1. Get Comments for a Post
Retrieve all comments for a specific post, including replies (nested structure).

**Endpoint:** `GET /posts/{post_id}/comments`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/posts/post-uuid-123/comments
```

**Response:**
```json
{
  "success": true,
  "message": "Comments retrieved successfully",
  "data": [
    {
      "id": "comment-uuid-abc",
      "post_id": "post-uuid-123",
      "user_id": "user-uuid-789",
      "user": {
        "username": "jane_smith",
        "display_name": "Jane Smith",
        "avatar_url": "https://lh3.googleusercontent.com/a/..."
      },
      "parent_id": null,
      "parent": null,
      "replies": [
        {
          "id": "comment-uuid-def",
          "post_id": "post-uuid-123",
          "user_id": "user-uuid-456",
          "user": {
            "username": "john_doe_123",
            "display_name": "John Doe",
            "avatar_url": "https://lh3.googleusercontent.com/a/..."
          },
          "parent_id": "comment-uuid-abc",
          "content": "Thanks! I'm really excited about it.",
          "like_count": 3,
          "created_at": "2023-12-01T11:15:00Z",
          "updated_at": "2023-12-01T11:15:00Z"
        }
      ],
      "content": "Wow, that's an amazing car! Congratulations!",
      "like_count": 8,
      "created_at": "2023-12-01T11:00:00Z",
      "updated_at": "2023-12-01T11:00:00Z"
    }
  ],
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const getComments = async (postId) => {
  const response = await fetch(`/api/v1/posts/${postId}/comments`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch comments');
};
```

---

### 2. Create Comment
Add a new comment to a post.

**Endpoint:** `POST /posts/{post_id}/comments`
**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "content": "That's an absolutely beautiful car! How does it drive?"
}
```

**Request:**
```http
POST /api/v1/posts/post-uuid-123/comments
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "content": "That's an absolutely beautiful car! How does it drive?"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Comment created successfully",
  "data": {
    "id": "comment-uuid-xyz",
    "post_id": "post-uuid-123",
    "user_id": "user-uuid-789",
    "user": {
      "username": "jane_smith",
      "display_name": "Jane Smith",
      "avatar_url": "https://lh3.googleusercontent.com/a/..."
    },
    "parent_id": null,
    "parent": null,
    "replies": [],
    "content": "That's an absolutely beautiful car! How does it drive?",
    "like_count": 0,
    "created_at": "2023-12-01T16:20:00Z",
    "updated_at": "2023-12-01T16:20:00Z"
  },
  "timestamp": "2023-12-01T16:20:00Z"
}
```

**Frontend Usage:**
```javascript
const createComment = async (postId, content) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/posts/${postId}/comments`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ content }),
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to create comment');
};
```

---

### 3. Create Reply to Comment
Reply to an existing comment.

**Endpoint:** `POST /posts/{post_id}/comments/{parent_comment_id}/reply`
**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "content": "It drives like a dream! The handling is incredible."
}
```

**Request:**
```http
POST /api/v1/posts/post-uuid-123/comments/comment-uuid-abc/reply
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "content": "It drives like a dream! The handling is incredible."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Reply created successfully",
  "data": {
    "id": "comment-uuid-reply-123",
    "post_id": "post-uuid-123",
    "user_id": "user-uuid-456",
    "user": {
      "username": "john_doe_123",
      "display_name": "John Doe",
      "avatar_url": "https://lh3.googleusercontent.com/a/..."
    },
    "parent_id": "comment-uuid-abc",
    "parent": {
      "id": "comment-uuid-abc",
      "content": "That's an absolutely beautiful car! How does it drive?",
      "user": {
        "username": "jane_smith",
        "display_name": "Jane Smith"
      }
    },
    "replies": [],
    "content": "It drives like a dream! The handling is incredible.",
    "like_count": 0,
    "created_at": "2023-12-01T16:25:00Z",
    "updated_at": "2023-12-01T16:25:00Z"
  },
  "timestamp": "2023-12-01T16:25:00Z"
}
```

**Frontend Usage:**
```javascript
const replyToComment = async (postId, parentCommentId, content) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/posts/${postId}/comments/${parentCommentId}/reply`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ content }),
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to create reply');
};
```

---

## Likes Endpoints

### 1. Toggle Like on Post
Like or unlike a post. Returns the new like status and updated like count.

**Endpoint:** `POST /posts/{post_id}/like`
**Authentication:** Required (Bearer token)

**Request:**
```http
POST /api/v1/posts/post-uuid-123/like
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Post like toggled successfully",
  "data": {
    "liked": true,
    "like_count": 43
  },
  "timestamp": "2023-12-01T16:30:00Z"
}
```

**Frontend Usage:**
```javascript
const togglePostLike = async (postId) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/posts/${postId}/like`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data; // { liked: true/false, like_count: number }
  }
  throw new Error(result.error?.message || 'Failed to toggle like');
};
```

---

### 2. Check User Like Status for Post
Check if the current user has liked a specific post.

**Endpoint:** `GET /posts/{post_id}/like/status`
**Authentication:** Required (Bearer token)

**Request:**
```http
GET /api/v1/posts/post-uuid-123/like/status
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "User like status retrieved successfully",
  "data": {
    "liked": true
  },
  "timestamp": "2023-12-01T16:35:00Z"
}
```

**Frontend Usage:**
```javascript
const checkUserLikedPost = async (postId) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/posts/${postId}/like/status`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data.liked;
  }
  throw new Error(result.error?.message || 'Failed to check like status');
};
```

---

### 3. Get Likes for Post
Retrieve all users who have liked a specific post.

**Endpoint:** `GET /posts/{post_id}/likes`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/posts/post-uuid-123/likes
```

**Response:**
```json
{
  "success": true,
  "message": "Post likes retrieved successfully",
  "data": {
    "likes": [
      {
        "id": "like-uuid-123",
        "post_id": "post-uuid-123",
        "user_id": "user-uuid-789",
        "user": {
          "username": "jane_smith",
          "display_name": "Jane Smith",
          "avatar_url": "https://lh3.googleusercontent.com/a/..."
        },
        "created_at": "2023-12-01T16:30:00Z"
      }
    ],
    "total_count": 43
  },
  "timestamp": "2023-12-01T16:40:00Z"
}
```

**Frontend Usage:**
```javascript
const getPostLikes = async (postId) => {
  const response = await fetch(`/api/v1/posts/${postId}/likes`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch post likes');
};
```

---

### 4. Toggle Like on Comment
Like or unlike a comment. Returns the new like status and updated like count.

**Endpoint:** `POST /posts/{post_id}/comments/{comment_id}/like`
**Authentication:** Required (Bearer token)

**Request:**
```http
POST /api/v1/posts/post-uuid-123/comments/comment-uuid-abc/like
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Comment like toggled successfully",
  "data": {
    "liked": true,
    "like_count": 9
  },
  "timestamp": "2023-12-01T16:45:00Z"
}
```

**Frontend Usage:**
```javascript
const toggleCommentLike = async (postId, commentId) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/posts/${postId}/comments/${commentId}/like`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to toggle comment like');
};
```

---

### 5. Check User Like Status for Comment
Check if the current user has liked a specific comment.

**Endpoint:** `GET /posts/{post_id}/comments/{comment_id}/like/status`
**Authentication:** Required (Bearer token)

**Request:**
```http
GET /api/v1/posts/post-uuid-123/comments/comment-uuid-abc/like/status
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "User comment like status retrieved successfully",
  "data": {
    "liked": false
  },
  "timestamp": "2023-12-01T16:50:00Z"
}
```

**Frontend Usage:**
```javascript
const checkUserLikedComment = async (postId, commentId) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/posts/${postId}/comments/${commentId}/like/status`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data.liked;
  }
  throw new Error(result.error?.message || 'Failed to check comment like status');
};
```

---

### 6. Get Likes for Comment
Retrieve all users who have liked a specific comment.

**Endpoint:** `GET /posts/{post_id}/comments/{comment_id}/likes`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/posts/post-uuid-123/comments/comment-uuid-abc/likes
```

**Response:**
```json
{
  "success": true,
  "message": "Comment likes retrieved successfully",
  "data": {
    "likes": [
      {
        "id": "like-uuid-456",
        "post_id": "post-uuid-123",
        "user_id": "user-uuid-456",
        "user": {
          "username": "john_doe_123",
          "display_name": "John Doe",
          "avatar_url": "https://lh3.googleusercontent.com/a/..."
        },
        "created_at": "2023-12-01T16:45:00Z"
      }
    ],
    "total_count": 9
  },
  "timestamp": "2023-12-01T16:55:00Z"
}
```

**Frontend Usage:**
```javascript
const getCommentLikes = async (postId, commentId) => {
  const response = await fetch(`/api/v1/posts/${postId}/comments/${commentId}/likes`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch comment likes');
};
```

---

## Users Endpoints

### 1. Get All Users
Retrieve a paginated list of all users with their basic profile information and follower counts.

**Endpoint:** `GET /users`
**Authentication:** Not required (Public)

**Query Parameters:**
- `page` (optional): Page number, default is 1
- `limit` (optional): Number of users per page, default is 20, max is 100

**Request:**
```http
GET /api/v1/users?page=1&limit=20
```

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": "user-uuid-123",
        "provider": "google",
        "first_name": "John",
        "last_name": "Doe",
        "username": "john_doe_123",
        "email": "john@example.com",
        "display_name": "John Doe",
        "bio": "Car enthusiast and software developer",
        "avatar_url": "https://lh3.googleusercontent.com/a/...",
        "follower_count": 45,
        "following_count": 23,
        "created_at": "2023-12-01T10:30:00Z"
      }
    ],
    "total_count": 100,
    "page": 1,
    "limit": 20,
    "has_next": true
  },
  "meta": {
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 100,
      "items_per_page": 20,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const getUsers = async (page = 1, limit = 20) => {
  const response = await fetch(`/api/v1/users?page=${page}&limit=${limit}`);
  const result = await response.json();
  
  if (result.success) {
    return {
      users: result.data.users,
      pagination: result.meta.pagination
    };
  }
  throw new Error(result.error?.message || 'Failed to fetch users');
};
```

---

### 2. Get Single User
Retrieve a specific user's profile information by ID.

**Endpoint:** `GET /users/{id}`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/users/user-uuid-123
```

**Response:**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "user-uuid-123",
    "provider": "google",
    "first_name": "John",
    "last_name": "Doe",
    "username": "john_doe_123",
    "email": "john@example.com",
    "display_name": "John Doe",
    "bio": "Car enthusiast and software developer. Love working on classic muscle cars in my spare time.",
    "avatar_url": "https://lh3.googleusercontent.com/a/...",
    "follower_count": 45,
    "following_count": 23,
    "created_at": "2023-12-01T10:30:00Z"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const getUser = async (userId) => {
  const response = await fetch(`/api/v1/users/${userId}`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch user');
};
```

---

### 3. Create User
Create a new user (typically used internally during OAuth flow).

**Endpoint:** `POST /users`
**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "provider_id": "google_id_123456789",
  "provider": "google",
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@example.com",
  "avatar_url": "https://lh3.googleusercontent.com/a/...",
  "locale": "en"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "user-uuid-456",
    "provider": "google",
    "first_name": "Jane",
    "last_name": "Smith",
    "username": null,
    "email": "jane.smith@example.com",
    "display_name": "Jane Smith",
    "bio": "",
    "avatar_url": "https://lh3.googleusercontent.com/a/...",
    "follower_count": 0,
    "following_count": 0,
    "created_at": "2023-12-01T16:00:00Z"
  },
  "timestamp": "2023-12-01T16:00:00Z"
}
```

**Frontend Usage:**
```javascript
const createUser = async (userData) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch('/api/v1/users', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to create user');
};
```

---

## Following System Endpoints

### 1. Toggle Follow User
Follow or unfollow a user. Returns the new follow status and updated follower/following counts.

**Endpoint:** `POST /users/{user_id}/follow`
**Authentication:** Required (Bearer token)

**Request:**
```http
POST /api/v1/users/user-uuid-456/follow
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "User follow toggled successfully",
  "data": {
    "is_following": true,
    "follower_count": 46,
    "following_count": 24
  },
  "timestamp": "2023-12-01T16:30:00Z"
}
```

**Frontend Usage:**
```javascript
const toggleUserFollow = async (userId) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/users/${userId}/follow`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data; // { is_following: true/false, follower_count: number, following_count: number }
  }
  throw new Error(result.error?.message || 'Failed to toggle follow');
};
```

---

### 2. Check Follow Status
Check if the current user is following a specific user.

**Endpoint:** `GET /users/{user_id}/follow/status`
**Authentication:** Required (Bearer token)

**Request:**
```http
GET /api/v1/users/user-uuid-456/follow/status
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "User follow status retrieved successfully",
  "data": {
    "is_following": true
  },
  "timestamp": "2023-12-01T16:35:00Z"
}
```

**Frontend Usage:**
```javascript
const checkFollowStatus = async (userId) => {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/v1/users/${userId}/follow/status`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  const result = await response.json();
  if (result.success) {
    return result.data.is_following;
  }
  throw new Error(result.error?.message || 'Failed to check follow status');
};
```

---

### 3. Get User's Followers
Retrieve a list of users who are following the specified user.

**Endpoint:** `GET /users/{user_id}/followers`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/users/user-uuid-456/followers
```

**Response:**
```json
{
  "success": true,
  "message": "User followers retrieved successfully",
  "data": {
    "followers": [
      {
        "id": "user-uuid-123",
        "username": "john_doe_123",
        "display_name": "John Doe",
        "avatar_url": "https://lh3.googleusercontent.com/a/...",
        "bio": "Car enthusiast and software developer"
      },
      {
        "id": "user-uuid-789",
        "username": "mike_wilson",
        "display_name": "Mike Wilson",
        "avatar_url": "https://lh3.googleusercontent.com/a/...",
        "bio": "Classic car collector"
      }
    ],
    "total_count": 46
  },
  "timestamp": "2023-12-01T16:40:00Z"
}
```

**Frontend Usage:**
```javascript
const getUserFollowers = async (userId) => {
  const response = await fetch(`/api/v1/users/${userId}/followers`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch followers');
};
```

---

### 4. Get User's Following
Retrieve a list of users that the specified user is following.

**Endpoint:** `GET /users/{user_id}/following`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/users/user-uuid-456/following
```

**Response:**
```json
{
  "success": true,
  "message": "User following list retrieved successfully",
  "data": {
    "following": [
      {
        "id": "user-uuid-111",
        "username": "sarah_jones",
        "display_name": "Sarah Jones",
        "avatar_url": "https://lh3.googleusercontent.com/a/...",
        "bio": "JDM car enthusiast"
      },
      {
        "id": "user-uuid-222",
        "username": "alex_rodriguez",
        "display_name": "Alex Rodriguez",
        "avatar_url": "https://lh3.googleusercontent.com/a/...",
        "bio": "Racing driver and car modifier"
      }
    ],
    "total_count": 24
  },
  "timestamp": "2023-12-01T16:45:00Z"
}
```

**Frontend Usage:**
```javascript
const getUserFollowing = async (userId) => {
  const response = await fetch(`/api/v1/users/${userId}/following`);
  const result = await response.json();
  
  if (result.success) {
    return result.data;
  }
  throw new Error(result.error?.message || 'Failed to fetch following list');
};
```

---

## Common Error Responses

### Posts/Users/Following Errors

**404 - Resource Not Found:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found",
    "details": "No user exists with the provided ID"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**403 - Forbidden Action:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Cannot follow yourself",
    "details": "Users are not allowed to follow themselves"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**400 - Invalid Request:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Content is required",
    "details": "Post content cannot be empty"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**429 - Rate Limited:**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests",
    "details": "You can only create 10 posts per hour"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

---

## Health Check Endpoint

### Health Status Check
Check the health and status of the application and its dependencies. This endpoint is useful for monitoring, load balancers, and container orchestration platforms.

**Endpoint:** `GET /health`
**Authentication:** Not required (Public)

**Request:**
```http
GET /api/v1/health
```

**Healthy Response (200 OK):**
```json
{
  "success": true,
  "message": "Health check completed",
  "data": {
    "status": "healthy",
    "timestamp": "2023-12-01T10:30:00Z",
    "version": "1.0.0",
    "services": {
      "database": {
        "status": "healthy",
        "response_time": "2.45ms",
        "last_checked": "2023-12-01T10:30:00Z"
      },
      "redis": {
        "status": "healthy",
        "response_time": "1.23ms",
        "last_checked": "2023-12-01T10:30:00Z"
      }
    },
    "uptime": "2h 15m 30s"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Degraded Response (200 OK - Redis issues but DB healthy):**
```json
{
  "success": true,
  "message": "Health check completed",
  "data": {
    "status": "degraded",
    "timestamp": "2023-12-01T10:30:00Z",
    "version": "1.0.0",
    "services": {
      "database": {
        "status": "healthy",
        "response_time": "3.21ms",
        "last_checked": "2023-12-01T10:30:00Z"
      },
      "redis": {
        "status": "unhealthy",
        "error": "Redis ping failed: connection refused",
        "last_checked": "2023-12-01T10:30:00Z"
      }
    },
    "uptime": "45m 12s"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Unhealthy Response (503 Service Unavailable):**
```json
{
  "success": false,
  "error": {
    "code": "SERVICE_UNHEALTHY",
    "message": "Service is unhealthy",
    "details": "Check individual service statuses"
  },
  "data": {
    "status": "unhealthy",
    "timestamp": "2023-12-01T10:30:00Z",
    "version": "1.0.0",
    "services": {
      "database": {
        "status": "unhealthy",
        "error": "Database ping failed: connection refused",
        "last_checked": "2023-12-01T10:30:00Z"
      },
      "redis": {
        "status": "unhealthy",
        "error": "Redis client not initialized",
        "last_checked": "2023-12-01T10:30:00Z"
      }
    },
    "uptime": "12s"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Usage:**
```javascript
const checkHealth = async () => {
  try {
    const response = await fetch('/api/v1/health');
    const result = await response.json();
    
    if (result.success) {
      const { status, services, uptime } = result.data;
      console.log(`API Status: ${status}`);
      console.log(`Database: ${services.database.status}`);
      console.log(`Redis: ${services.redis.status}`);
      console.log(`Uptime: ${uptime}`);
      return result.data;
    } else {
      // Service is unhealthy
      console.error('API is unhealthy:', result.error.message);
      return result.data; // Still contains health info
    }
  } catch (error) {
    console.error('Health check failed:', error);
    throw error;
  }
};

// For monitoring dashboards
const getServiceHealth = async () => {
  const health = await checkHealth();
  return {
    isHealthy: health.status === 'healthy',
    isDegraded: health.status === 'degraded',
    services: health.services,
    uptime: health.uptime
  };
};
```

### Health Status Values

- **`healthy`**: All services are functioning normally
- **`degraded`**: Core services (database) are healthy, but some non-critical services (Redis) may have issues
- **`unhealthy`**: Critical services (database) are not functioning properly

### Monitoring Integration

**Docker Health Check:**
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1
```

**Kubernetes Readiness Probe:**
```yaml
readinessProbe:
  httpGet:
    path: /api/v1/health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

**Load Balancer Health Check:**
```
Health Check URL: /api/v1/health
Expected Status: 200
Check Interval: 30s
Timeout: 5s
Healthy Threshold: 2
Unhealthy Threshold: 3
```
