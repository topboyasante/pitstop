# Pitstop MVP - Twitter X StackOverflow for Car Owners

## Vision
A real-time social Q&A platform where car owners share quick questions, tips, and expertise in a Twitter-like feed with StackOverflow's knowledge structure.

## Core MVP Features

### 1. User Management
- **User Registration/Login**
  - Email/password authentication
  - Profile with bio, location, car garage (multiple cars)
  - Follow/unfollow other users
  - User reputation system

### 2. Posts (Twitter-style)
- **Create Posts**
  - Short questions/tips (280 characters + optional media)
  - Quick car problems, tips, or show-offs
  - Photo/video uploads
  - Car tags (#Honda #Maintenance #DIY)
  - Location tagging (optional)
- **Post Feed**
  - Real-time timeline of followed users + trending posts
  - Algorithmic "For You" feed
  - Filter by post type (questions, tips, show-offs)

### 3. Detailed Q&A (StackOverflow-style)
- **Long-form Questions**
  - Detailed technical questions
  - Step-by-step problem descriptions
  - Multiple photos/diagrams
  - Car specifications
- **Expert Answers**
  - Detailed solutions with explanations
  - Step-by-step guides
  - Accept best answer

### 4. Social Interactions
- **Engagement**
  - Like/heart posts and answers
  - Repost/share with comments
  - Quick replies vs detailed answers
  - Bookmark useful posts
- **Following**
  - Follow mechanics, experts, car enthusiasts
  - Notifications for activity from followed users

### 5. Discovery & Trending
- **Real-time Discovery**
  - Trending hashtags (#BrakeProblems #Tesla #DIY)
  - Popular posts in your area
  - Hot questions needing answers
- **Smart Feed**
  - Posts relevant to your car(s)
  - Questions you can answer based on expertise

### 6. Reputation & Gamification
- **Social Credit**
  - Likes, reposts, followers count
  - Expert badges for specific car brands/topics
  - Quick vs detailed answer points

## Technical MVP Requirements

### Database Schema
- Users (id, email, username, bio, reputation, location, follower_count)
- Cars (id, user_id, make, model, year, nickname)
- Posts (id, user_id, content, post_type, media_urls, location, likes, reposts)
- Questions (id, user_id, title, content, tags, car_id, status, views)
- Answers (id, question_id, user_id, content, votes, is_accepted)
- Follows (id, follower_id, following_id)
- Likes (id, user_id, post_id/answer_id, type)
- Hashtags (id, name, usage_count)

### API Endpoints
- Auth: POST /auth/register, POST /auth/login
- Feed: GET /feed/timeline, GET /feed/for-you
- Posts: GET/POST /posts, POST /posts/:id/like, POST /posts/:id/repost
- Questions: GET/POST /questions, GET /questions/:id
- Answers: POST /questions/:id/answers, POST /answers/:id/vote
- Users: GET/POST /users/:id/follow, GET /users/:id/cars
- Search: GET /search?q=query&type=posts&hashtags=maintenance

## Out of Scope (Future Features)
- Real-time WebSocket notifications
- Advanced content moderation
- Car VIN integration & automatic specs
- Marketplace/parts selling
- Live video streaming
- Advanced analytics dashboard
- Private messaging/DMs
- Verified mechanic badges
- AI-powered answer suggestions

## Success Metrics
- 100+ registered users with active profiles
- 500+ posts/questions in first month
- 70% user retention after first week
- Average 5+ interactions per post
- 50+ car garage entries
- Active hashtag usage (#CarProblems trending)