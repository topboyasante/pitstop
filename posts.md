# Posts API Integration Guide for Frontend

This document explains how to integrate with the Pitstop backend Posts API from your frontend application.

## Overview

The Posts API allows you to create, retrieve, and manage posts in the Pitstop platform. All posts endpoints require authentication via JWT tokens.

**Note:** All GET requests for posts (both single post and posts list) automatically include essential user information (username, display name, and avatar URL) for the post author, eliminating the need for separate user API calls.

All API responses follow a **structured format** for consistency:

```json
{
  "success": true/false,
  "message": "Human-readable message",
  "data": { ... },           // Only present on success
  "error": { ... },          // Only present on failure
  "meta": { ... },           // Pagination info (when applicable)
  "timestamp": "2023-12-01T10:30:00Z"
}
```

## Authentication Required

All Posts API endpoints require authentication. Include the JWT access token in the Authorization header:

```http
Authorization: Bearer <access_token>
```

## API Endpoints

### 1. Get All Posts

Retrieve a paginated list of all posts.

**Request:**
```http
GET /api/v1/posts?page=1&limit=20
Authorization: Bearer <access_token>
```

**Parameters:**
- `page` (optional): Page number, default is 1
- `limit` (optional): Number of posts per page, default is 20

**Response:**
```json
{
  "success": true,
  "message": "Posts retrieved successfully",
  "data": [
    {
      "id": "post-uuid-abc123",
      "content": "This is my first post!",
      "user_id": "user-uuid-123",
      "user": {
        "username": "johndoe",
        "display_name": "John Doe",
        "avatar_url": "https://example.com/avatar/john.jpg"
      },
      "created_at": "2023-12-01T10:30:00Z",
      "updated_at": "2023-12-01T10:30:00Z"
    },
    {
      "id": "post-uuid-def456",
      "content": "Another interesting post...",
      "user_id": "user-uuid-456",
      "user": {
        "username": "janesmith",
        "display_name": "Jane Smith",
        "avatar_url": "https://example.com/avatar/jane.jpg"
      },
      "created_at": "2023-12-01T09:15:00Z",
      "updated_at": "2023-12-01T09:15:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 157,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Implementation:**
```javascript
const getPosts = async (page = 1, limit = 20) => {
  try {
    const accessToken = localStorage.getItem('access_token');
    const response = await fetch(`/api/v1/posts?page=${page}&limit=${limit}`, {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
    });
    
    const result = await response.json();
    
    if (result.success) {
      return {
        posts: result.data,
        meta: result.meta,
      };
    } else {
      throw new Error(result.error?.message || 'Failed to retrieve posts');
    }
  } catch (error) {
    console.error('Error fetching posts:', error);
    throw error;
  }
};
```

### 2. Create New Post

Create a new post.

**Request:**
```http
POST /api/v1/posts
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "content": "This is my new post content!",
  "user_id": "user-uuid-123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Post created successfully", 
  "data": {
    "id": "post-uuid-xyz789",
    "content": "This is my new post content!",
    "user_id": "user-uuid-123",
    "created_at": "2023-12-01T10:30:00Z",
    "updated_at": "2023-12-01T10:30:00Z"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Implementation:**
```javascript
const createPost = async (content, userId) => {
  try {
    const accessToken = localStorage.getItem('access_token');
    const response = await fetch('/api/v1/posts', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        content: content,
        user_id: userId,
      }),
    });
    
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error?.message || 'Failed to create post');
    }
  } catch (error) {
    console.error('Error creating post:', error);
    throw error;
  }
};
```

### 3. Get Single Post

Retrieve a specific post by ID.

**Request:**
```http
GET /api/v1/posts/{id}
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Post retrieved successfully",
  "data": {
    "id": "post-uuid-abc123",
    "content": "This is the specific post content.",
    "user_id": "user-uuid-123",
    "user": {
      "username": "johndoe",
      "display_name": "John Doe",
      "avatar_url": "https://example.com/avatar/john.jpg"
    },
    "created_at": "2023-12-01T10:30:00Z",
    "updated_at": "2023-12-01T10:30:00Z"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Frontend Implementation:**
```javascript
const getPost = async (postId) => {
  try {
    const accessToken = localStorage.getItem('access_token');
    const response = await fetch(`/api/v1/posts/${postId}`, {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
      },
    });
    
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.error?.message || 'Failed to retrieve post');
    }
  } catch (error) {
    console.error('Error fetching post:', error);
    throw error;
  }
};
```


## Error Handling

### Common Error Responses

**Unauthorized (401):**
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

**Validation Error (400):**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request body",
    "details": "content field is required"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Post Not Found (404):**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Post not found"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**Rate Limit Exceeded (429):**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Try again later.",
    "details": "Request ID: req_abc123"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

### Error Handling Best Practices

```javascript
const handlePostsApiError = async (response) => {
  const result = await response.json();
  
  if (!result.success && result.error) {
    const { code, message, details } = result.error;
    
    switch (code) {
      case 'UNAUTHORIZED':
      case 'INVALID_TOKEN':
        // Auth errors - redirect to login
        localStorage.clear();
        window.location.href = '/login';
        break;
      case 'VALIDATION_ERROR':
        // Show validation error to user
        throw new Error(details || message);
      case 'NOT_FOUND':
        // Handle post not found
        throw new Error('The requested post was not found');
      case 'RATE_LIMIT_EXCEEDED':
        // Rate limit - show rate limit message
        throw new Error('Too many requests. Please try again later.');
      case 'INTERNAL_ERROR':
        // Server error - show generic message
        throw new Error('Something went wrong. Please try again.');
      default:
        throw new Error(message || 'Request failed');
    }
  }
  
  throw new Error('Request failed');
};
```

## Pagination Helper

Here's a simple pagination helper function to work with the API metadata:

```javascript
const buildPaginationInfo = (meta) => {
  if (!meta) return null;
  
  const { page, total_pages, has_next, has_prev, total } = meta;
  
  return {
    currentPage: page,
    totalPages: total_pages,
    hasNext: has_next,
    hasPrevious: has_prev,
    totalItems: total,
    canGoNext: has_next,
    canGoPrevious: has_prev,
    nextPage: has_next ? page + 1 : null,
    previousPage: has_prev ? page - 1 : null,
  };
};

// Example usage
const displayPosts = async (page = 1) => {
  try {
    const result = await getPosts(page);
    const paginationInfo = buildPaginationInfo(result.meta);
    
    // Display posts
    result.posts.forEach(post => {
      console.log(`Post ${post.id} by ${post.user?.display_name || post.user?.username}: ${post.content}`);
    });
    
    // Display pagination info
    if (paginationInfo) {
      console.log(`Page ${paginationInfo.currentPage} of ${paginationInfo.totalPages}`);
      console.log(`Total posts: ${paginationInfo.totalItems}`);
    }
    
    return { posts: result.posts, pagination: paginationInfo };
  } catch (error) {
    console.error('Error displaying posts:', error);
  }
};
```

## Data Validation

### Post Content Validation

```javascript
const validatePostContent = (content) => {
  const errors = [];
  
  if (!content || content.trim().length === 0) {
    errors.push('Post content is required');
  }
  
  if (content && content.length > 1000) {
    errors.push('Post content must be less than 1000 characters');
  }
  
  if (content && content.trim().length < 3) {
    errors.push('Post content must be at least 3 characters');
  }
  
  return errors;
};

// Example usage in a form submission
const handlePostSubmission = async (formData) => {
  const content = formData.get('content');
  const userId = JSON.parse(localStorage.getItem('user')).id;
  
  // Validate content
  const validationErrors = validatePostContent(content);
  if (validationErrors.length > 0) {
    console.error('Validation errors:', validationErrors);
    return { success: false, errors: validationErrors };
  }
  
  try {
    const newPost = await createPost(content, userId);
    console.log('Post created successfully:', newPost);
    return { success: true, post: newPost };
  } catch (error) {
    console.error('Failed to create post:', error);
    return { success: false, errors: [error.message] };
  }
};

// Input sanitization helper
const sanitizePostContent = (content) => {
  // Remove potentially harmful content
  return content
    .trim()
    .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '') // Remove script tags
    .replace(/<[^>]*>/g, '') // Remove all HTML tags
    .replace(/javascript:/gi, '') // Remove javascript: protocols
    .substring(0, 1000); // Limit length
};
```

## Testing Posts API

For testing, you can mock the posts responses:

```javascript
// Mock posts data for testing
const mockPostsResponse = {
  success: true,
  message: "Posts retrieved successfully",
  data: [
    {
      id: "test-post-uuid-123",
      content: "Test post content",
      user_id: "test-user-id",
      user: {
        username: "testuser",
        display_name: "Test User",
        avatar_url: "https://example.com/avatar/test.jpg"
      },
      created_at: "2023-12-01T10:30:00Z",
      updated_at: "2023-12-01T10:30:00Z"
    }
  ],
  meta: {
    page: 1,
    limit: 20,
    total: 1,
    total_pages: 1,
    has_next: false,
    has_prev: false
  },
  timestamp: "2023-12-01T10:30:00Z"
};

const mockCreatePostResponse = {
  success: true,
  message: "Post created successfully",
  data: {
    id: "test-post-uuid-456",
    content: "New test post",
    user_id: "test-user-id",
    created_at: "2023-12-01T10:30:00Z",
    updated_at: "2023-12-01T10:30:00Z"
  },
  timestamp: "2023-12-01T10:30:00Z"
};
```

## Security Best Practices

1. **Authentication**: Always include valid JWT tokens
2. **Input Validation**: Validate post content on the frontend
3. **Rate Limiting**: Handle rate limit responses gracefully
4. **Error Handling**: Don't expose sensitive error details
5. **Content Sanitization**: Sanitize post content before display
6. **HTTPS**: Always use HTTPS in production

## Environment Configuration

Configure your frontend for different environments:

```javascript
// config.js
export const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Use in your API calls
const response = await fetch(`${API_BASE_URL}/api/v1/posts`);
```

This completes the frontend integration guide for the Pitstop Posts API. The backend provides consistent structured responses and proper error handling, making frontend integration straightforward and reliable.