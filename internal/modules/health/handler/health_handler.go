package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/topboyasante/pitstop/internal/core/response"
	"github.com/topboyasante/pitstop/internal/modules/health/dto"
	"gorm.io/gorm"
)

var startTime = time.Now()

// HealthHandler handles health check requests
type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// HealthCheck performs comprehensive health checks
// @Summary Health check endpoint
// @Description Check the health status of the application and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Success 503 {object} response.APIResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	services := make(map[string]dto.ServiceInfo)
	overallStatus := dto.StatusHealthy
	
	// Check database connectivity
	dbInfo := h.checkDatabase()
	services["database"] = dbInfo
	if dbInfo.Status != dto.StatusHealthy {
		overallStatus = dto.StatusUnhealthy
	}
	
	// Check Redis connectivity
	redisInfo := h.checkRedis()
	services["redis"] = redisInfo
	if redisInfo.Status != dto.StatusHealthy {
		if overallStatus == dto.StatusHealthy {
			overallStatus = dto.StatusDegraded
		}
	}
	
	// Calculate uptime
	uptime := time.Since(startTime)
	
	healthResponse := dto.HealthCheckResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0", // You might want to get this from config or build info
		Services:  services,
		Uptime:    formatUptime(uptime),
	}
	
	// Return appropriate HTTP status code
	if overallStatus == dto.StatusUnhealthy {
		c.Status(fiber.StatusServiceUnavailable)
		return response.ErrorJSON(c, fiber.StatusServiceUnavailable, "SERVICE_UNHEALTHY", "Service is unhealthy", "Check individual service statuses")
	}
	
	return response.SuccessJSON(c, &healthResponse, "Health check completed")
}

// checkDatabase verifies database connectivity
func (h *HealthHandler) checkDatabase() dto.ServiceInfo {
	start := time.Now()
	info := dto.ServiceInfo{
		LastChecked: time.Now(),
	}
	
	// Get underlying sql.DB from GORM
	sqlDB, err := h.db.DB()
	if err != nil {
		info.Status = dto.StatusUnhealthy
		info.Error = fmt.Sprintf("Failed to get database connection: %v", err)
		return info
	}
	
	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		info.Status = dto.StatusUnhealthy
		info.Error = fmt.Sprintf("Database ping failed: %v", err)
		return info
	}
	
	// Check connection stats
	stats := sqlDB.Stats()
	if stats.OpenConnections == 0 {
		info.Status = dto.StatusUnhealthy
		info.Error = "No database connections available"
		return info
	}
	
	info.Status = dto.StatusHealthy
	info.ResponseTime = fmt.Sprintf("%.2fms", float64(time.Since(start).Nanoseconds())/1e6)
	return info
}

// checkRedis verifies Redis connectivity
func (h *HealthHandler) checkRedis() dto.ServiceInfo {
	start := time.Now()
	info := dto.ServiceInfo{
		LastChecked: time.Now(),
	}
	
	if h.redis == nil {
		info.Status = dto.StatusUnhealthy
		info.Error = "Redis client not initialized"
		return info
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// Ping Redis
	pong, err := h.redis.Ping(ctx).Result()
	if err != nil {
		info.Status = dto.StatusUnhealthy
		info.Error = fmt.Sprintf("Redis ping failed: %v", err)
		return info
	}
	
	if pong != "PONG" {
		info.Status = dto.StatusUnhealthy
		info.Error = "Unexpected Redis ping response"
		return info
	}
	
	info.Status = dto.StatusHealthy
	info.ResponseTime = fmt.Sprintf("%.2fms", float64(time.Since(start).Nanoseconds())/1e6)
	return info
}

// formatUptime formats duration into human-readable uptime
func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}