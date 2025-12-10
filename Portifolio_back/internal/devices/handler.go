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

const invalidID = "invalid id"
const NotFoundDevice = "device not found"
const unauthorizedError = "unauthorized"

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError})
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
		c.JSON(http.StatusNotFound, gin.H{"error": NotFoundDevice})
		return
	}

	// Ensure user owns this device
	if device.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": NotFoundDevice})
		return
	}

	c.JSON(http.StatusOK, device)
}

// Create adds a new device for the authenticated user.
func (h *Handler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
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
		c.JSON(http.StatusNotFound, gin.H{"error": NotFoundDevice})
		return
	}

	c.JSON(http.StatusOK, device)
}

// Delete removes a device by ID.
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
		return
	}

	// Get current device
	device, err := h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": NotFoundDevice})
		return
	}

	// Ensure user owns this device
	if device.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": NotFoundDevice})
		return
	}

	desiredPower := getDesiredPowerState(device.PowerState)
	connIP, connUser, connPass := parseTapoMetadata(device.Metadata)

	if err := h.toggleDeviceState(c, id, device, desiredPower, connIP, connUser, connPass); err != nil {
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

func getDesiredPowerState(currentPowerState *bool) bool {
	if currentPowerState != nil && *currentPowerState {
		return false
	}
	return true
}

func parseTapoMetadata(metadata string) (ip, username, password string) {
	if metadata == "" {
		return "", "", ""
	}

	var meta map[string]any
	if err := json.Unmarshal([]byte(metadata), &meta); err != nil {
		return "", "", ""
	}

	tapoData, ok := meta["tapo"].(map[string]any)
	if !ok {
		return "", "", ""
	}

	if ipVal, ok := tapoData["ip"].(string); ok {
		ip = ipVal
	}
	if userVal, ok := tapoData["username"].(string); ok {
		username = userVal
	}
	if passVal, ok := tapoData["password"].(string); ok {
		password = passVal
	}

	return ip, username, password
}

func (h *Handler) toggleDeviceState(c *gin.Context, id int64, device *Device, desiredPower bool, connIP, connUser, connPass string) error {
	if connIP != "" && connUser != "" && connPass != "" {
		return h.togglePhysicalDevice(c, id, desiredPower, connIP, connUser, connPass)
	}
	return h.toggleDatabaseStatus(c, id, device.Status)
}

func (h *Handler) togglePhysicalDevice(c *gin.Context, id int64, desiredPower bool, connIP, connUser, connPass string) error {
	// conn := integrations_tapo.Connection{IP: connIP, Username: connUser, Password: connPass}
	connIP = "192.168.237.52"
	connUser = "pedrohdcosta@gmail.com"
	connPass = "gdnz8cti1"
	conn := integrations_tapo.Connection{IP: connIP, Username: connUser, Password: connPass}
	if err := integrations_tapo.SetPower(c.Request.Context(), conn, desiredPower); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to control physical device", "detail": err.Error()})
		return err
	}

	if err := h.Repo.UpdatePowerState(c.Request.Context(), id, desiredPower); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device power state"})
		return err
	}

	if err := h.Repo.UpdateStatus(c.Request.Context(), id, "online"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device status"})
		return err
	}

	return nil
}

func (h *Handler) toggleDatabaseStatus(c *gin.Context, id int64, currentStatus string) error {
	newStatus := "online"
	if currentStatus == "online" {
		newStatus = "offline"
	}

	if err := h.Repo.UpdateStatus(c.Request.Context(), id, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to toggle device status"})
		return err
	}

	return nil
}

// ReadPower performs a live read of power consumption from the device if configured.
func (h *Handler) ReadPower(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	device, err := h.getDeviceForUser(c, userID)
	if err != nil {
		return
	}

	connIP, connUser, connPass := parseTapoMetadata(device.Metadata)
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

func (h *Handler) getDeviceForUser(c *gin.Context, userID int64) (*Device, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidID})
		return nil, err
	}

	device, err := h.Repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return nil, err
	}

	if device.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "device not found or not owned by user"})
		return nil, err
	}

	return device, nil
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
