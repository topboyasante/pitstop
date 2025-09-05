# Authentication Integration Guide for Frontend

This document explains how to integrate with the Pitstop backend authentication system from your frontend application.

## Overview

The Pitstop API uses **Google OAuth2 authentication** with **JWT tokens** for session management. The authentication flow involves redirecting users to Google for login, then exchanging the authorization code for JWT tokens.

All API responses follow a **structured format** for consistency:

```json
{
  "success": true/false,
  "message": "Human-readable message",
  "data": { ... },           // Only present on success
  "error": { ... },          // Only present on failure
  "timestamp": "2023-12-01T10:30:00Z"
}
```

## Authentication Flow

### 1. Initiate Authentication

Start the authentication process by getting the Google OAuth URL:

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

**Frontend Implementation:**
```javascript
// Get OAuth URL and redirect user
const initiateAuth = async () => {
  try {
    const response = await fetch('/api/v1/auth/google');
    const result = await response.json();
    
    if (result.success) {
      // Redirect user to Google OAuth
      window.location.href = result.data.auth_url;
    } else {
      console.error('Auth initiation failed:', result.error?.message);
    }
  } catch (error) {
    console.error('Failed to initiate authentication:', error);
  }
};
```

### 2. Handle OAuth Callback

After successful Google authentication, users are redirected to:
```
http://your-frontend-url.com/?code=AUTH_CODE&state=STATE_TOKEN
```

**Frontend Implementation:**
```javascript
// Handle OAuth callback (in your callback page/component)
const handleAuthCallback = async () => {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  const state = urlParams.get('state');
  
  if (code && state) {
    await exchangeCodeForTokens(code, state);
  } else {
    // Handle error - missing code or state
    console.error('Missing authorization code or state');
  }
};
```

### 3. Exchange Code for Tokens

Exchange the authorization code for JWT tokens:

**Request:**
```http
POST /api/v1/auth/exchange
Content-Type: application/json

{
  "code": "AUTHORIZATION_CODE",
  "state": "STATE_TOKEN"
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

**Frontend Implementation:**
```javascript
const exchangeCodeForTokens = async (code, state) => {
  try {
    const response = await fetch('/api/v1/auth/exchange', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ code, state }),
    });
    
    const result = await response.json();
    
    if (result.success) {
      // Store tokens securely
      localStorage.setItem('access_token', result.data.access_token);
      localStorage.setItem('refresh_token', result.data.refresh_token);
      
      // Get user info and store it
      await getCurrentUser();
      
      // Redirect to dashboard or home page
      window.location.href = '/dashboard';
    } else {
      throw new Error(result.error?.message || 'Token exchange failed');
    }
  } catch (error) {
    console.error('Token exchange error:', error);
    // Handle error (show message to user)
  }
};
```

## Token Management

### Access Tokens
- **Duration**: 30 minutes
- **Usage**: Include in `Authorization` header for API requests
- **Format**: `Bearer <access_token>`

### Refresh Tokens
- **Duration**: 30 days
- **Usage**: Used to get new access tokens when they expire

### Token Storage

**Recommended approach:**
```javascript
// Store tokens securely
const storeTokens = (accessToken, refreshToken) => {
  // Option 1: localStorage (simple but less secure)
  localStorage.setItem('access_token', accessToken);
  localStorage.setItem('refresh_token', refreshToken);
  
  // Option 2: httpOnly cookies (more secure, requires backend support)
  // Tokens set via Set-Cookie headers from backend
  
  // Option 3: Secure session storage
  sessionStorage.setItem('access_token', accessToken);
  sessionStorage.setItem('refresh_token', refreshToken);
};
```

## Making Authenticated Requests

Include the access token in the `Authorization` header:

```javascript
const makeAuthenticatedRequest = async (url, options = {}) => {
  const accessToken = localStorage.getItem('access_token');
  
  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
    },
  });
  
  if (response.status === 401) {
    // Token expired, try to refresh
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      // Retry original request
      return makeAuthenticatedRequest(url, options);
    } else {
      // Refresh failed, redirect to login
      redirectToLogin();
    }
  }
  
  return response;
};
```

## Token Refresh

When access tokens expire (401 response), use the refresh token:

**Request:**
```http
POST /api/v1/auth/refresh
Content-Type: application/json

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

**Frontend Implementation:**
```javascript
const refreshAccessToken = async () => {
  try {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) return false;
    
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    
    const result = await response.json();
    
    if (result.success) {
      // Update stored tokens
      localStorage.setItem('access_token', result.data.access_token);
      localStorage.setItem('refresh_token', result.data.refresh_token);
      
      return true;
    } else {
      // Refresh failed, clear tokens
      console.error('Token refresh failed:', result.error?.message);
      clearAuthData();
      return false;
    }
  } catch (error) {
    console.error('Token refresh error:', error);
    clearAuthData();
    return false;
  }
};
```

## Get Current User

Retrieve current user information:

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
    "provider": "google",
    "first_name": "John",
    "last_name": "Doe",
    "username": "john_doe_123",
    "email": "user@example.com",
    "display_name": "John Doe",
    "bio": "Software Developer",
    "avatar_url": "https://...",
    "created_at": "2023-12-01T10:30:00Z"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

## Complete Authentication Hook (React)

Here's a complete React hook for authentication:

```javascript
import { useState, useEffect, createContext, useContext } from 'react';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const accessToken = localStorage.getItem('access_token');
      if (!accessToken) {
        setLoading(false);
        return;
      }

      const response = await fetch('/api/v1/auth/me', {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
        },
      });

      const result = await response.json();

      if (result.success) {
        setUser(result.data);
      } else if (response.status === 401) {
        // Try to refresh token
        const refreshed = await refreshAccessToken();
        if (refreshed) {
          checkAuthStatus(); // Retry
        } else {
          clearAuthData();
        }
      }
    } catch (error) {
      console.error('Auth check error:', error);
    }
    setLoading(false);
  };

  const login = async () => {
    try {
      const response = await fetch('/api/v1/auth/google');
      const result = await response.json();
      
      if (result.success) {
        window.location.href = result.data.auth_url;
      } else {
        console.error('Login error:', result.error?.message);
      }
    } catch (error) {
      console.error('Login error:', error);
    }
  };

  const logout = () => {
    clearAuthData();
    setUser(null);
    window.location.href = '/';
  };

  const clearAuthData = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
  };

  const refreshAccessToken = async () => {
    try {
      const refreshToken = localStorage.getItem('refresh_token');
      if (!refreshToken) return false;

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
    } catch (error) {
      console.error('Token refresh error:', error);
      return false;
    }
  };

  const value = {
    user,
    loading,
    login,
    logout,
    isAuthenticated: !!user,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
```

## Error Handling

### Common Error Responses

**Invalid/Expired Token (401):**
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

**Missing Authorization Header (401):**
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

**Invalid Authorization Code (400):**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Failed to exchange authorization code",
    "details": "authorization code was not found"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

**State Token Mismatch (400):**
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
const handleApiError = async (response) => {
  const result = await response.json();
  
  // Check if it's a structured error response
  if (!result.success && result.error) {
    const { code, message, details } = result.error;
    
    switch (code) {
      case 'UNAUTHORIZED':
      case 'INVALID_TOKEN':
      case 'INVALID_CLAIMS':
        // Auth errors - redirect to login
        clearAuthData();
        window.location.href = '/login';
        break;
      case 'VALIDATION_ERROR':
        // Validation errors - show specific message
        throw new Error(details || message);
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
  
  // Fallback for non-structured responses
  switch (response.status) {
    case 401:
      clearAuthData();
      window.location.href = '/login';
      break;
    case 500:
      throw new Error('Something went wrong. Please try again.');
    default:
      throw new Error('Request failed');
  }
};
```

## Security Best Practices

1. **Token Storage**: Consider using httpOnly cookies for better security
2. **HTTPS**: Always use HTTPS in production
3. **Token Validation**: Check token expiry before making requests
4. **Logout**: Clear all stored tokens on logout
5. **Error Handling**: Don't expose sensitive error details to users
6. **CSRF Protection**: The backend handles CSRF via state tokens
7. **Token Refresh**: Implement automatic token refresh for better UX

## Environment Configuration

Your frontend should be configured with:

```javascript
// config.js
export const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
export const FRONTEND_URL = process.env.REACT_APP_FRONTEND_URL || 'http://localhost:3000';
```

## Testing Authentication

For testing, you can mock the authentication responses:

```javascript
// Mock authentication for testing
const mockAuthResponse = {
  access_token: 'mock-access-token',
  refresh_token: 'mock-refresh-token',
  expires_in: 1800,
  user: {
    id: 'test-user-id',
    email: 'test@example.com',
    first_name: 'Test',
    last_name: 'User',
    username: 'test_user',
    display_name: 'Test User',
    avatar_url: 'https://example.com/avatar.jpg',
    bio: 'Test user bio',
    locale: 'en'
  }
};
```

## Structured API Helper Function

For consistent handling of the structured response format, use this helper function:

```javascript
const apiCall = async (url, options = {}) => {
  try {
    const response = await fetch(url, options);
    const result = await response.json();
    
    if (result.success) {
      return result.data;
    } else {
      // Handle structured error
      const error = new Error(result.error?.message || 'API call failed');
      error.code = result.error?.code;
      error.details = result.error?.details;
      throw error;
    }
  } catch (error) {
    if (error.code) throw error; // Structured error, re-throw
    throw new Error('Network or parsing error');
  }
};

// Usage example
const getCurrentUser = async () => {
  try {
    const user = await apiCall('/api/v1/auth/me', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
      },
    });
    return user;
  } catch (error) {
    if (error.code === 'UNAUTHORIZED') {
      // Handle auth error
      clearAuthData();
      window.location.href = '/login';
    }
    throw error;
  }
};
```

This completes the frontend integration guide for the Pitstop authentication system. The backend handles all the OAuth complexity, state management, and security measures, while providing **consistent structured responses** for easy frontend integration.