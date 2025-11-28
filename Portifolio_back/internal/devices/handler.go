package devices

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles device HTTP requests.
type Handler struct {
	Repo *Repo
}

// NewHandler creates a new device handler.
func NewHandler(repo *Repo) *Handler {
	return &Handler{Repo: repo}
}

// RegisterRoutes registers device routes on the Gin engine.
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/devices")
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/toggle", h.Toggle)
}

// List returns all devices for the authenticated user.
func (h *Handler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	devices, err := h.Repo.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch devices"})
		return
	}

	if devices == nil {
		devices = []Device{}
	}
	c.JSON(http.StatusOK, devices)
}

// Get returns a single device by ID.
func (h *Handler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
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

	device, err := h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	// Ensure user owns this device
	if device.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// Create adds a new device for the authenticated user.
func (h *Handler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: name is required"})
		return
	}

	device := &Device{
		UserID:   userID,
		Name:     req.Name,
		Room:     req.Room,
		Type:     req.Type,
		Status:   "offline",
		Metadata: req.Metadata,
	}

	id, err := h.Repo.Create(c.Request.Context(), device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create device"})
		return
	}

	device.ID = id
	c.JSON(http.StatusCreated, device)
}

// Update modifies an existing device.
func (h *Handler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
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

	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if err := h.Repo.Update(c.Request.Context(), userID, id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device"})
		return
	}

	// Fetch updated device
	device, err := h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// Delete removes a device by ID.
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
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

	if err := h.Repo.Delete(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete device"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Toggle switches device power state (on/off).
func (h *Handler) Toggle(c *gin.Context) {
	userID, ok := getUserID(c)
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

	// Get current device
	device, err := h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	// Ensure user owns this device
	if device.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	// Toggle power state
	var newPowerState bool
	if device.PowerState == nil || !*device.PowerState {
		newPowerState = true
	} else {
		newPowerState = false
	}

	if err := h.Repo.UpdatePowerState(c.Request.Context(), id, newPowerState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to toggle device"})
		return
	}

	// Fetch updated device
	device, err = h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated device"})
		return
	}

	c.JSON(http.StatusOK, device)
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
