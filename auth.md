# Authentication Integration Guide for Frontend

This document explains how to integrate with the Pitstop backend authentication system from your frontend application.

## Overview

The Pitstop API uses **Google OAuth2 authentication** with **JWT tokens** for session management. The authentication flow involves redirecting users to Google for login, then exchanging the authorization code for JWT tokens.

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
  "auth_url": "https://accounts.google.com/o/oauth2/auth?client_id=...&state=..."
}
```

**Frontend Implementation:**
```javascript
// Get OAuth URL and redirect user
const initiateAuth = async () => {
  try {
    const response = await fetch('/api/v1/auth/google');
    const { auth_url } = await response.json();
    
    // Redirect user to Google OAuth
    window.location.href = auth_url;
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
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 1800,
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "username": "john_doe_123",
    "display_name": "John Doe",
    "avatar_url": "https://...",
    "bio": "Software Developer",
    "locale": "en"
  }
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
    
    if (response.ok) {
      const data = await response.json();
      
      // Store tokens securely
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      
      // Store user data
      localStorage.setItem('user', JSON.stringify(data.user));
      
      // Redirect to dashboard or home page
      window.location.href = '/dashboard';
    } else {
      throw new Error('Token exchange failed');
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
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 1800
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
    
    if (response.ok) {
      const data = await response.json();
      
      // Update stored tokens
      localStorage.setItem('access_token', data.access_token);
      localStorage.setItem('refresh_token', data.refresh_token);
      
      return true;
    } else {
      // Refresh failed, clear tokens
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
  "id": "user-uuid",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "username": "john_doe_123",
  "display_name": "John Doe",
  "avatar_url": "https://...",
  "bio": "Software Developer",
  "locale": "en"
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

      if (response.ok) {
        const userData = await response.json();
        setUser(userData);
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
      const { auth_url } = await response.json();
      window.location.href = auth_url;
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

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
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
  "error": "Unauthorized",
  "message": "Invalid or expired token"
}
```

**Missing Authorization Header (401):**
```json
{
  "error": "Unauthorized", 
  "message": "Missing authorization header"
}
```

**Invalid Authorization Code (400):**
```json
{
  "error": "Bad Request",
  "message": "Invalid authorization code"
}
```

**State Token Mismatch (400):**
```json
{
  "error": "Bad Request",
  "message": "Invalid state token"
}
```

### Error Handling Best Practices

```javascript
const handleApiError = (response) => {
  switch (response.status) {
    case 401:
      // Unauthorized - redirect to login
      clearAuthData();
      window.location.href = '/login';
      break;
    case 400:
      // Bad request - show error message
      return response.json().then(data => {
        throw new Error(data.message);
      });
    case 500:
      // Server error - show generic message
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

This completes the frontend integration guide for the Pitstop authentication system. The backend handles all the OAuth complexity, state management, and security measures, making frontend integration straightforward.