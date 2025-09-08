package dto

import "time"

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Services  map[string]ServiceInfo `json:"services"`
	Uptime    string                 `json:"uptime"`
}

// ServiceInfo represents the status of a service dependency
type ServiceInfo struct {
	Status      string        `json:"status"`
	ResponseTime string       `json:"response_time,omitempty"`
	Error       string        `json:"error,omitempty"`
	LastChecked time.Time     `json:"last_checked"`
}

// HealthStatus constants
const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
	StatusDegraded  = "degraded"
)