# Pitstop Database Models - Twitter X StackOverflow

## Core Models for MVP

### User
```go
type User struct {
    ID            uint           `gorm:"primarykey" json:"id"`
    Username      string         `gorm:"uniqueIndex;not null;size:50" json:"username" validate:"required,min=3,max=50"`
    Email         string         `gorm:"uniqueIndex;not null;size:255" json:"email" validate:"required,email"`
    Password      string         `gorm:"not null" json:"-" validate:"required,min=8"`
    Bio           string         `gorm:"size:160" json:"bio"` // Twitter-style bio
    Location      string         `gorm:"size:100" json:"location"`
    Reputation    int            `gorm:"default:0" json:"reputation"`
    FollowerCount int            `gorm:"default:0" json:"follower_count"`
    FollowingCount int           `gorm:"default:0" json:"following_count"`
    PostCount     int            `gorm:"default:0" json:"post_count"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    Cars          []Car          `gorm:"foreignKey:UserID" json:"cars,omitempty"`
    Posts         []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
    Questions     []Question     `gorm:"foreignKey:UserID" json:"questions,omitempty"`
    Answers       []Answer       `gorm:"foreignKey:UserID" json:"answers,omitempty"`
    Votes         []Vote         `gorm:"foreignKey:UserID" json:"votes,omitempty"`
    Likes         []Like         `gorm:"foreignKey:UserID" json:"likes,omitempty"`
}
```

### Car
```go
type Car struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    UserID    uint           `gorm:"not null" json:"user_id"`
    Make      string         `gorm:"not null;size:50" json:"make" validate:"required"`
    Model     string         `gorm:"not null;size:50" json:"model" validate:"required"`
    Year      int            `gorm:"not null" json:"year" validate:"required,min=1900,max=2030"`
    Nickname  string         `gorm:"size:50" json:"nickname"` // "My Baby", "The Beast", etc.
    Color     string         `gorm:"size:30" json:"color"`
    Mileage   int            `json:"mileage"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Questions []Question     `gorm:"foreignKey:CarID" json:"questions,omitempty"`
    Posts     []Post         `gorm:"foreignKey:CarID" json:"posts,omitempty"`
}
```

### Post (Twitter-style)
```go
type Post struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    UserID      uint           `gorm:"not null" json:"user_id"`
    CarID       *uint          `json:"car_id"` // Optional car reference
    Content     string         `gorm:"size:280;not null" json:"content" validate:"required,min=1,max=280"`
    PostType    string         `gorm:"default:general" json:"post_type"` // general, question, tip, showoff
    MediaURLs   string         `gorm:"size:1000" json:"media_urls"` // JSON array of media URLs
    Location    string         `gorm:"size:100" json:"location"`
    LikeCount   int            `gorm:"default:0" json:"like_count"`
    RepostCount int            `gorm:"default:0" json:"repost_count"`
    ReplyCount  int            `gorm:"default:0" json:"reply_count"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Car         *Car           `gorm:"foreignKey:CarID" json:"car,omitempty"`
    Likes       []Like         `gorm:"foreignKey:PostID" json:"likes,omitempty"`
    Reposts     []Repost       `gorm:"foreignKey:PostID" json:"reposts,omitempty"`
}
```

### Question (StackOverflow-style)
```go
type Question struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    UserID      uint           `gorm:"not null" json:"user_id"`
    CarID       *uint          `json:"car_id"` // Optional specific car reference
    Title       string         `gorm:"not null;size:255" json:"title" validate:"required,min=10,max=255"`
    Content     string         `gorm:"type:text;not null" json:"content" validate:"required,min=20"`
    Tags        string         `gorm:"size:500" json:"tags"` // JSON array as string
    Views       int            `gorm:"default:0" json:"views"`
    Status      string         `gorm:"default:open" json:"status"` // open, answered, closed
    AnswerCount int            `gorm:"default:0" json:"answer_count"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Car         *Car           `gorm:"foreignKey:CarID" json:"car,omitempty"`
    Answers     []Answer       `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}
```

### Answer
```go
type Answer struct {
    ID         uint           `gorm:"primarykey" json:"id"`
    QuestionID uint           `gorm:"not null" json:"question_id"`
    UserID     uint           `gorm:"not null" json:"user_id"`
    Content    string         `gorm:"type:text;not null" json:"content" validate:"required,min=10"`
    IsAccepted bool           `gorm:"default:false" json:"is_accepted"`
    VoteScore  int            `gorm:"default:0" json:"vote_score"` // Cached vote total
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    Question   Question       `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
    User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Votes      []Vote         `gorm:"foreignKey:AnswerID" json:"votes,omitempty"`
}
```

### Vote
```go
type Vote struct {
    ID       uint      `gorm:"primarykey" json:"id"`
    UserID   uint      `gorm:"not null" json:"user_id"`
    AnswerID uint      `gorm:"not null" json:"answer_id"`
    VoteType string    `gorm:"not null;check:vote_type IN ('up', 'down')" json:"vote_type" validate:"required,oneof=up down"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relationships
    User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Answer   Answer    `gorm:"foreignKey:AnswerID" json:"answer,omitempty"`
}
```

### Like (Twitter-style)
```go
type Like struct {
    ID       uint      `gorm:"primarykey" json:"id"`
    UserID   uint      `gorm:"not null" json:"user_id"`
    PostID   *uint     `json:"post_id"`   // For liking posts
    AnswerID *uint     `json:"answer_id"` // For liking answers
    CreatedAt time.Time `json:"created_at"`
    
    // Relationships
    User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Post     *Post     `gorm:"foreignKey:PostID" json:"post,omitempty"`
    Answer   *Answer   `gorm:"foreignKey:AnswerID" json:"answer,omitempty"`
}
```

### Repost (Twitter-style)
```go
type Repost struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    PostID    uint      `gorm:"not null" json:"post_id"`
    Comment   string    `gorm:"size:280" json:"comment"` // Optional comment when reposting
    CreatedAt time.Time `json:"created_at"`
    
    // Relationships
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
```

### Follow (Twitter-style)
```go
type Follow struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    FollowerID  uint      `gorm:"not null" json:"follower_id"`  // User who follows
    FollowingID uint      `gorm:"not null" json:"following_id"` // User being followed
    CreatedAt   time.Time `json:"created_at"`
    
    // Relationships
    Follower    User      `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
    Following   User      `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}
```


## Database Indexes

### Performance Indexes
```sql
-- Users
CREATE INDEX idx_users_reputation ON users(reputation DESC);
CREATE INDEX idx_users_follower_count ON users(follower_count DESC);
CREATE INDEX idx_users_username ON users(username);

-- Cars
CREATE INDEX idx_cars_user_id ON cars(user_id);
CREATE INDEX idx_cars_make_model ON cars(make, model);

-- Posts (Twitter-style)
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_post_type ON posts(post_type);
CREATE INDEX idx_posts_like_count ON posts(like_count DESC);

-- Questions (StackOverflow-style)
CREATE INDEX idx_questions_user_id ON questions(user_id);
CREATE INDEX idx_questions_car_id ON questions(car_id);
CREATE INDEX idx_questions_status ON questions(status);
CREATE INDEX idx_questions_created_at ON questions(created_at DESC);
CREATE INDEX idx_questions_views ON questions(views DESC);

-- Answers
CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_answers_user_id ON answers(user_id);
CREATE INDEX idx_answers_is_accepted ON answers(is_accepted);
CREATE INDEX idx_answers_vote_score ON answers(vote_score DESC);

-- Social Features
CREATE UNIQUE INDEX idx_follows_follower_following ON follows(follower_id, following_id);
CREATE UNIQUE INDEX idx_votes_user_answer ON votes(user_id, answer_id);
CREATE UNIQUE INDEX idx_likes_user_post ON likes(user_id, post_id) WHERE post_id IS NOT NULL;
CREATE UNIQUE INDEX idx_likes_user_answer ON likes(user_id, answer_id) WHERE answer_id IS NOT NULL;
CREATE UNIQUE INDEX idx_reposts_user_post ON reposts(user_id, post_id);

```

## Relationships

### One-to-Many
- User → Cars (1 user has many cars)
- User → Posts (1 user has many posts)
- User → Questions (1 user has many questions)
- User → Answers (1 user has many answers)
- User → Votes (1 user has many votes)
- User → Likes (1 user has many likes)
- User → Reposts (1 user has many reposts)
- Car → Posts (1 car has many posts)
- Car → Questions (1 car has many questions)
- Question → Answers (1 question has many answers)
- Post → Likes (1 post has many likes)
- Post → Reposts (1 post has many reposts)
- Answer → Votes (1 answer has many votes)
- Answer → Likes (1 answer has many likes)

### Many-to-Many
- User ↔ User (follows relationship)

### Constraints
- Follow: Unique constraint on (follower_id, following_id) - prevent duplicate follows
- Vote: Unique constraint on (user_id, answer_id) - user can only vote once per answer
- Like: Unique constraint on (user_id, post_id) and (user_id, answer_id) - prevent duplicate likes
- Repost: Unique constraint on (user_id, post_id) - prevent duplicate reposts
- Question status: Must be 'open', 'answered', or 'closed'
- Post type: Must be 'general', 'question', 'tip', or 'showoff'
- Vote type: Must be 'up' or 'down'

## Migration Order
1. Users
2. Cars
3. Posts
4. Questions
5. Answers
6. Follows
7. Votes
8. Likes
9. Reposts
10. Add indexes