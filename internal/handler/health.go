package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler health check handler
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type healthResponse struct {
	OK  bool   `json:"ok"`
	DB  string `json:"db"`
	Now string `json:"now"`
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check service health and DB connection
// @Tags         system
// @Produce      json
// @Success      200  {object}  healthResponse
// @Router       /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := "down"
	if sqlDB, err := h.db.DB(); err == nil {
		if err := sqlDB.Ping(); err == nil {
			status = "up"
		}
	}

	c.JSON(http.StatusOK, healthResponse{
		OK:  true,
		DB:  status,
		Now: time.Now().UTC().Format(time.RFC3339),
	})
}
