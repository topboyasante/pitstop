package services

import "github.com/topboyasante/pitstop/internal/config"

type AuthService struct {
	config *config.Config
}

func NewAuthService(config *config.Config) *AuthService {
	return &AuthService{
		config: config,
	}
}

func (as *AuthService) Authenticate() string {
	url := as.config.OAuth.AuthCodeURL("randomstate")
	return url
}
