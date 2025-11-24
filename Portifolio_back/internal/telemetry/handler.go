package telemetry

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler handles telemetry HTTP requests.
type Handler struct {
	Repo *Repo
}

// NewHandler creates a new telemetry handler.
func NewHandler(repo *Repo) *Handler {
	return &Handler{Repo: repo}
}

// RegisterRoutes registers telemetry routes on the Gin engine.
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/telemetry")
	g.GET("", h.List)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}

// List returns telemetry data for the authenticated user or a specific device.
// Query params: device_id (optional), limit (default 100)
func (h *Handler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	deviceIDStr := c.Query("device_id")
	if deviceIDStr != "" {
		deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
			return
		}
		// Use ListByDeviceForUser to ensure user owns the device
		data, err := h.Repo.ListByDeviceForUser(c.Request.Context(), userID, deviceID, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch telemetry"})
			return
		}
		if data == nil {
			data = []Telemetry{}
		}
		c.JSON(http.StatusOK, data)
		return
	}

	// List all telemetry for user's devices
	data, err := h.Repo.ListByUser(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch telemetry"})
		return
	}
	if data == nil {
		data = []Telemetry{}
	}
	c.JSON(http.StatusOK, data)
}

// Create adds a new telemetry record.
func (h *Handler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateTelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Verify user owns the device
	owns, err := h.Repo.UserOwnsDevice(c.Request.Context(), userID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify device ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "device not found or not owned by user"})
		return
	}

	t := &Telemetry{
		DeviceID:  req.DeviceID,
		Power:     req.Power,
		Voltage:   req.Voltage,
		Current:   req.Current,
		Timestamp: time.Now(),
	}

	id, err := h.Repo.Create(c.Request.Context(), t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create telemetry"})
		return
	}

	t.ID = id
	c.JSON(http.StatusCreated, t)
}

// Delete removes a telemetry record by ID.
func (h *Handler) Delete(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.Repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete telemetry"})
		return
	}

	c.Status(http.StatusNoContent)
}

// getUserID extracts user ID from Gin context (set by auth middleware).
func getUserID(c *gin.Context) (int64, bool) {
	sub := c.GetString("sub")
	if sub == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
