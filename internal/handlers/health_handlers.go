package handlers

import (
	"context"
	"database/sql"
	"sparkle-concept-sync/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	db           *sql.DB
	redisService *services.RedisService
}

func NewHealthHandler(db *sql.DB, redisService *services.RedisService) *HealthHandler {
	return &HealthHandler{
		db:           db,
		redisService: redisService,
	}
}

// HealthCheck provides basic health status
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "sparkle-concept-sync",
		"version":   "1.0.0",
	})
}

// DetailedHealthCheck provides comprehensive health status
func (h *HealthHandler) DetailedHealthCheck(c *fiber.Ctx) error {
	health := fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "sparkle-concept-sync",
		"version":   "1.0.0",
		"checks": fiber.Map{
			"database": h.checkDatabase(),
			"redis":    h.checkRedis(),
		},
	}

	// Determine overall status
	dbHealthy := health["checks"].(fiber.Map)["database"].(fiber.Map)["status"] == "healthy"
	redisHealthy := health["checks"].(fiber.Map)["redis"].(fiber.Map)["status"] == "healthy"

	if !dbHealthy || !redisHealthy {
		health["status"] = "unhealthy"
		return c.Status(fiber.StatusServiceUnavailable).JSON(health)
	}

	return c.JSON(health)
}

func (h *HealthHandler) checkDatabase() fiber.Map {
	start := time.Now()

	err := h.db.Ping()
	duration := time.Since(start)

	if err != nil {
		return fiber.Map{
			"status":   "unhealthy",
			"error":    err.Error(),
			"duration": duration.String(),
		}
	}

	return fiber.Map{
		"status":   "healthy",
		"duration": duration.String(),
	}
}

func (h *HealthHandler) checkRedis() fiber.Map {
	start := time.Now()

	// Try to ping Redis
	err := h.redisService.Ping(context.Background())
	duration := time.Since(start)

	if err != nil {
		return fiber.Map{
			"status":   "unhealthy",
			"error":    err.Error(),
			"duration": duration.String(),
		}
	}

	return fiber.Map{
		"status":   "healthy",
		"duration": duration.String(),
	}
}
