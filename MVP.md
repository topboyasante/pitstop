# Pitstop MVP - Simple Social Platform for Car Enthusiasts

## Vision
A simple social platform where car owners can share posts, connect with other enthusiasts, and build a community around automotive interests.

## Current MVP Scope (What's Actually Built)

### 1. User Management ✅ 
- **Google OAuth Authentication**
  - Google login/signup flow
  - JWT token-based authentication
  - Token refresh functionality
  - User profile management

### 2. Basic Posts ✅
- **Simple Post Creation**
  - Create text-based posts
  - View individual posts
  - Get paginated list of all posts
  - Basic CRUD operations

### 3. User Profiles ✅
- **Profile Information**
  - Display name, bio, avatar
  - Basic user information from OAuth
  - Get user by ID
  - List all users with pagination

## Next Phase Features (Not Yet Built)

### Phase 2: Social Interactions
- **Post Engagement**
  - Like posts
  - Comment on posts
  - User can delete their own posts

### Phase 3: Car Context
- **Car Garage**
  - Add cars to profile (make, model, year)
  - Tag posts with specific car
  - Filter posts by car type

### Phase 4: Enhanced Social
- **Follow System**
  - Follow/unfollow users
  - Timeline feed of followed users
  - User follower/following counts

## Current Database Schema

### Users
```go
type User struct {
    ID          string    // Primary key
    ProviderID  string    // OAuth provider ID
    Provider    string    // google, facebook, github
    FirstName   string
    LastName    string
    Username    string    // Unique
    Email       string    // Unique
    DisplayName string
    Bio         string
    AvatarURL   string
    Locale      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Posts
```go
type Post struct {
    ID        uint      // Primary key
    UserID    string    // Foreign key to users
    Content   string    // Post text content
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## API Endpoints Currently Available

### Authentication
- `GET /auth/google` - Get Google OAuth URL
- `GET /auth/google/callback` - Handle OAuth callback
- `POST /auth/exchange` - Exchange auth code for JWT tokens
- `POST /auth/refresh` - Refresh JWT tokens
- `GET /auth/me` - Get current user info

### Users
- `GET /users` - Get all users (paginated)
- `POST /users` - Create new user
- `GET /users/{id}` - Get user by ID

### Posts
- `GET /posts` - Get all posts (paginated)
- `POST /posts` - Create new post
- `GET /posts/{id}` - Get post by ID

## Success Metrics for Current MVP
- 50+ registered users via Google OAuth
- 100+ posts created in first month  
- Users posting at least once per week
- Basic user engagement with platform

## Deferred Features (Too Complex for Initial MVP)
- Complex social features (likes, comments, reposts)
- Car-specific functionality (garage, car tagging)
- Advanced feeds and discovery
- Real-time features
- Mobile app
- Image/video uploads
- Advanced user profiles
- Reputation systems
- Search functionality