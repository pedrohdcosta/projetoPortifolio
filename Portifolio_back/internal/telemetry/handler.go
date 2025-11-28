package telemetry

import (
	"log"
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
	g.GET("/latest", h.ListLatest)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}

// RegisterDeviceTelemetryRoutes registers device-specific telemetry routes.
func (h *Handler) RegisterDeviceTelemetryRoutes(r *gin.RouterGroup) {
	// These routes are registered under /api/devices/:id/telemetry
	r.GET("/:id/telemetry", h.ListByDevice)
	r.GET("/:id/telemetry/summary", h.GetSummary)
	r.GET("/:id/telemetry/latest", h.GetLatest)
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

// ListLatest returns the latest telemetry reading for each device of the user.
func (h *Handler) ListLatest(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	data, err := h.Repo.GetLatestByUserDevices(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch latest telemetry"})
		return
	}
	if data == nil {
		data = []Telemetry{}
	}
	c.JSON(http.StatusOK, data)
}

// ListByDevice returns telemetry for a specific device (via /api/devices/:id/telemetry).
func (h *Handler) ListByDevice(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	limit := 100
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	data, err := h.Repo.ListByDeviceForUser(c.Request.Context(), userID, deviceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch telemetry"})
		return
	}
	if data == nil {
		data = []Telemetry{}
	}
	c.JSON(http.StatusOK, data)
}

// GetSummary returns aggregated telemetry summary for a device.
// Query params: period (day|week|month, default: day)
func (h *Handler) GetSummary(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	period := c.DefaultQuery("period", "day")

	summary, err := h.Repo.GetSummaryByDevice(c.Request.Context(), userID, deviceID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch telemetry summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetLatest returns the most recent telemetry reading for a device.
func (h *Handler) GetLatest(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	// First verify user owns the device
	owns, err := h.Repo.UserOwnsDevice(c.Request.Context(), userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify device ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "device not found or not owned by user"})
		return
	}

	telemetry, err := h.Repo.GetLatestByDevice(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no telemetry data found"})
		return
	}

	c.JSON(http.StatusOK, telemetry)
}

// Create adds a new telemetry record and updates device status.
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

	// Update device last_seen and status to 'online'
	if err := h.Repo.UpdateDeviceLastSeenAndStatus(c.Request.Context(), req.DeviceID); err != nil {
		log.Printf("Warning: failed to update device status for device %d: %v", req.DeviceID, err)
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
