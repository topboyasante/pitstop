package events

import "time"

// Event represents a domain event in the system
type Event interface {
	EventName() string
	EventTime() time.Time
	EventData() any
}

// BaseEvent provides common event functionality
type BaseEvent struct {
	Name      string    `json:"event_name"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

func (e BaseEvent) EventName() string {
	return e.Name
}

func (e BaseEvent) EventTime() time.Time {
	return e.Timestamp
}

func (e BaseEvent) EventData() any {
	return e.Data
}

// User Events
type UserRegistered struct {
	BaseEvent
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewUserRegistered(userID uint, email, username string) *UserRegistered {
	return &UserRegistered{
		BaseEvent: BaseEvent{
			Name:      "user.registered",
			Timestamp: time.Now(),
		},
		UserID:   userID,
		Email:    email,
		Username: username,
	}
}

type UserProfileUpdated struct {
	BaseEvent
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

func NewUserProfileUpdated(userID uint, username string) *UserProfileUpdated {
	return &UserProfileUpdated{
		BaseEvent: BaseEvent{
			Name:      "user.profile_updated",
			Timestamp: time.Now(),
		},
		UserID:   userID,
		Username: username,
	}
}