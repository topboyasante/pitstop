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

// Auth Events
type AuthenticationSuccessful struct {
	BaseEvent
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	AvatarURL  string `json:"avatar_url"`
	Locale     string `json:"locale"`
}

func NewAuthenticationSuccessful(provider,
	providerID, email, firstName, lastName,
	avatarURL, locale string) *AuthenticationSuccessful {
	return &AuthenticationSuccessful{
		BaseEvent: BaseEvent{
			Name:      "oauth.authentication_successful",
			Timestamp: time.Now(),
		},
		Provider:   provider,
		ProviderID: providerID,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		AvatarURL:  avatarURL,
		Locale:     locale,
	}
}
