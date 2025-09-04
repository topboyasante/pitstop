# Pitstop Database Models - Simplified MVP

## Current Models (Phase 1)

### User (OAuth-based)
```go
type User struct {
    ID          string         `gorm:"primarykey" json:"id"`
    ProviderID  string         `gorm:"not null;size:100;index" json:"provider_id"`
    Provider    string         `gorm:"not null;size:50;index" json:"provider" validate:"required,oneof=google facebook github"`
    FirstName   string         `gorm:"size:255" json:"first_name" validate:"omitempty,max=255"`
    LastName    string         `gorm:"size:255" json:"last_name" validate:"omitempty,max=255"`
    Username    string         `gorm:"uniqueIndex;size:100" json:"username" validate:"omitempty,min=3,max=100,alphanum"`
    Email       string         `gorm:"uniqueIndex;not null;size:255" json:"email" validate:"required,email,max=255"`
    DisplayName string         `gorm:"size:150" json:"display_name" validate:"omitempty,max=150"`
    Bio         string         `gorm:"size:500" json:"bio" validate:"omitempty,max=500"`
    AvatarURL   string         `gorm:"size:500" json:"avatar_url" validate:"omitempty,url,max=500"`
    Locale      string         `gorm:"size:10" json:"locale" validate:"omitempty,max=10"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Simple relationship for now
    Posts       []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}
```

### Post (Simple)
```go
type Post struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    string    `gorm:"not null" json:"user_id" validate:"required"`
    Content   string    `gorm:"type:text" json:"content" validate:"required"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relationship
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

## Future Models (Not Yet Implemented)

### Phase 2: Comments/Likes
```go
type Like struct {
    ID       uint      `gorm:"primarykey" json:"id"`
    UserID   string    `gorm:"not null" json:"user_id"`
    PostID   uint      `gorm:"not null" json:"post_id"`
    CreatedAt time.Time `json:"created_at"`
    
    User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

type Comment struct {
    ID       uint      `gorm:"primarykey" json:"id"`
    UserID   string    `gorm:"not null" json:"user_id"`
    PostID   uint      `gorm:"not null" json:"post_id"`
    Content  string    `gorm:"type:text" json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
```

### Phase 3: Car Garage
```go
type Car struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    string    `gorm:"not null" json:"user_id"`
    Make      string    `gorm:"not null;size:50" json:"make"`
    Model     string    `gorm:"not null;size:50" json:"model"`
    Year      int       `gorm:"not null" json:"year"`
    Nickname  string    `gorm:"size:50" json:"nickname"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

### Phase 4: Follow System
```go
type Follow struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    FollowerID  string    `gorm:"not null" json:"follower_id"`
    FollowingID string    `gorm:"not null" json:"following_id"`
    CreatedAt   time.Time `json:"created_at"`
    
    Follower    User      `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
    Following   User      `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}
```

## Current Database Indexes

### Performance Indexes (Phase 1)
```sql
-- Users
CREATE INDEX idx_users_provider ON users(provider);
CREATE INDEX idx_users_provider_id ON users(provider_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- Posts
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

## Current Relationships

### One-to-Many (Phase 1)
- User → Posts (1 user has many posts)

## Migration Order (Phase 1)
1. Users
2. Posts