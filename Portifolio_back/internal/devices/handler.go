package devices

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	integrations_tapo "github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/integrations/tapo"
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
	g.GET("/:id/read", h.ReadPower)
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

// Toggle switches device status (online/offline).
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

	// Toggle power state. If device metadata contains tapo connection info,
	// attempt to perform the operation on the real device. Otherwise, fall back
	// to toggling the status/power state in the database only.

	// Determine desired power state
	var desiredPower bool
	if device.PowerState != nil && *device.PowerState {
		desiredPower = false
	} else {
		desiredPower = true
	}

	// Try to parse metadata for tapo connection info
	var connIP, connUser, connPass string
	if device.Metadata != "" {
		// metadata expected to be a JSON object containing fields like
		// {"tapo": {"ip": "192.168.1.10", "username": "admin", "password": "pw"}}
		var meta map[string]any
		if err := json.Unmarshal([]byte(device.Metadata), &meta); err == nil {
			if t, ok := meta["tapo"].(map[string]any); ok {
				if ip, ok := t["ip"].(string); ok {
					connIP = ip
				}
				if u, ok := t["username"].(string); ok {
					connUser = u
				}
				if p, ok := t["password"].(string); ok {
					connPass = p
				}
			}
		}
	}

	// If we have connection info, attempt to call the tapo integration.
	if connIP != "" && connUser != "" && connPass != "" {
		// attempt network operation but do not fail the request entirely on error;
		// update DB only on success.
		conn := integrations_tapo.Connection{IP: connIP, Username: connUser, Password: connPass}
		if err := integrations_tapo.SetPower(c.Request.Context(), conn, desiredPower); err != nil {
			// log and return error to client
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to control physical device", "detail": err.Error()})
			return
		}

		// Update stored power_state and last_seen/status
		if err := h.Repo.UpdatePowerState(c.Request.Context(), id, desiredPower); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device power state"})
			return
		}
		// mark device online when successful
		if err := h.Repo.UpdateStatus(c.Request.Context(), id, "online"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device status"})
			return
		}
	} else {
		// No tapo metadata — toggle status in DB only (legacy behavior)
		var newStatus string
		if device.Status == "online" {
			newStatus = "offline"
		} else {
			newStatus = "online"
		}

		if err := h.Repo.UpdateStatus(c.Request.Context(), id, newStatus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to toggle device status"})
			return
		}
	}

	// Fetch updated device
	device, err = h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated device"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// ReadPower performs a live read of power consumption from the device if configured.
func (h *Handler) ReadPower(c *gin.Context) {
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

	// Get device
	device, err := h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	if device.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "device not found or not owned by user"})
		return
	}

	// parse metadata for tapo
	var connIP, connUser, connPass string
	if device.Metadata != "" {
		var meta map[string]any
		if err := json.Unmarshal([]byte(device.Metadata), &meta); err == nil {
			if t, ok := meta["tapo"].(map[string]any); ok {
				if ip, ok := t["ip"].(string); ok {
					connIP = ip
				}
				if u, ok := t["username"].(string); ok {
					connUser = u
				}
				if p, ok := t["password"].(string); ok {
					connPass = p
				}
			}
		}
	}

	if connIP == "" || connUser == "" || connPass == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device not configured for live reads"})
		return
	}

	conn := integrations_tapo.Connection{IP: connIP, Username: connUser, Password: connPass}
	power, err := integrations_tapo.ReadPower(c.Request.Context(), conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read power from device", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"power": power})
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
