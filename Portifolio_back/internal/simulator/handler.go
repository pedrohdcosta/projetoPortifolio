package simulator

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// Seed random number generator for better randomness
	rand.Seed(time.Now().UnixNano())
}

// SimulatorConfig defines configuration for telemetry simulation.
type SimulatorConfig struct {
	BasePower   float64 `json:"base_power"`   // Base power in Watts (e.g., 100)
	Variation   float64 `json:"variation"`    // Random variation percentage (e.g., 0.2 = 20%)
	BaseVoltage float64 `json:"base_voltage"` // Base voltage in Volts (e.g., 220)
	Count       int     `json:"count"`        // Number of telemetry readings to generate
	IntervalSec int     `json:"interval_sec"` // Interval between readings in seconds (for historical data)
}

// SimulatedTelemetry represents a generated telemetry reading.
type SimulatedTelemetry struct {
	DeviceID  int64     `json:"device_id"`
	Power     float64   `json:"power"`
	Voltage   float64   `json:"voltage"`
	Current   float64   `json:"current"`
	Timestamp time.Time `json:"timestamp"`
}

// TelemetryCreator interface for creating telemetry records.
type TelemetryCreator interface {
	CreateTelemetry(deviceID int64, power, voltage, current float64, timestamp time.Time) (int64, error)
	UserOwnsDevice(userID, deviceID int64) (bool, error)
	UpdateDeviceStatus(deviceID int64) error
}

// Handler handles simulator HTTP requests.
type Handler struct {
	Creator TelemetryCreator
}

// NewHandler creates a new simulator handler.
func NewHandler(creator TelemetryCreator) *Handler {
	return &Handler{Creator: creator}
}

// RegisterRoutes registers simulator routes on the Gin engine.
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/simulator")
	g.POST("/generate/:device_id", h.Generate)
	g.POST("/bulk/:device_id", h.GenerateBulk)
}

// Generate creates a single simulated telemetry reading for a device.
func (h *Handler) Generate(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// Verify user owns the device
	owns, err := h.Creator.UserOwnsDevice(userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify device ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "device not found or not owned by user"})
		return
	}

	// Default config for a smart plug simulation
	config := SimulatorConfig{
		BasePower:   150,  // 150W typical for a small appliance
		Variation:   0.15, // 15% variation
		BaseVoltage: 220,  // 220V
	}

	// Allow custom config from request body (ignore errors - use defaults for invalid JSON)
	if err := c.ShouldBindJSON(&config); err != nil {
		log.Printf("Using default simulator config (invalid JSON or empty body)")
	}

	// Generate simulated reading
	telemetry := generateTelemetry(deviceID, config, time.Now())

	// Save to database
	id, err := h.Creator.CreateTelemetry(
		deviceID,
		telemetry.Power,
		telemetry.Voltage,
		telemetry.Current,
		telemetry.Timestamp,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create telemetry"})
		return
	}

	// Update device status
	if err := h.Creator.UpdateDeviceStatus(deviceID); err != nil {
		log.Printf("Warning: failed to update device status for device %d: %v", deviceID, err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        id,
		"device_id": telemetry.DeviceID,
		"power":     telemetry.Power,
		"voltage":   telemetry.Voltage,
		"current":   telemetry.Current,
		"timestamp": telemetry.Timestamp,
	})
}

// GenerateBulk creates multiple simulated telemetry readings (historical data).
func (h *Handler) GenerateBulk(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// Verify user owns the device
	owns, err := h.Creator.UserOwnsDevice(userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify device ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "device not found or not owned by user"})
		return
	}

	// Default config for bulk generation
	config := SimulatorConfig{
		BasePower:   150,
		Variation:   0.15,
		BaseVoltage: 220,
		Count:       24,  // 24 readings by default
		IntervalSec: 300, // 5 minutes between readings
	}

	// Allow custom config from request body
	if err := c.ShouldBindJSON(&config); err != nil {
		// Use defaults if no body provided
	}

	// Validate limits
	if config.Count <= 0 {
		config.Count = 24
	}
	if config.Count > 100 {
		config.Count = 100 // Max 100 readings at once
	}
	if config.IntervalSec <= 0 {
		config.IntervalSec = 300
	}

	// Generate historical readings
	now := time.Now()
	created := 0
	for i := config.Count - 1; i >= 0; i-- {
		timestamp := now.Add(-time.Duration(i*config.IntervalSec) * time.Second)
		telemetry := generateTelemetry(deviceID, config, timestamp)

		_, err := h.Creator.CreateTelemetry(
			deviceID,
			telemetry.Power,
			telemetry.Voltage,
			telemetry.Current,
			telemetry.Timestamp,
		)
		if err == nil {
			created++
		}
	}

	// Update device status
	if err := h.Creator.UpdateDeviceStatus(deviceID); err != nil {
		log.Printf("Warning: failed to update device status for device %d: %v", deviceID, err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "bulk telemetry generated",
		"device_id":       deviceID,
		"readings_created": created,
		"base_power":      config.BasePower,
		"variation":       config.Variation,
		"interval_sec":    config.IntervalSec,
	})
}

// generateTelemetry creates a simulated telemetry reading with realistic values.
func generateTelemetry(deviceID int64, config SimulatorConfig, timestamp time.Time) SimulatedTelemetry {
	// Add random variation to base power
	variation := (rand.Float64()*2 - 1) * config.Variation * config.BasePower
	power := config.BasePower + variation

	// Voltage with small variation (typically ±5%)
	voltageVariation := (rand.Float64()*2 - 1) * 0.05 * config.BaseVoltage
	voltage := config.BaseVoltage + voltageVariation

	// Calculate current: I = P / V (with zero check)
	var current float64
	if voltage != 0 {
		current = power / voltage
	}

	return SimulatedTelemetry{
		DeviceID:  deviceID,
		Power:     round2(power),
		Voltage:   round2(voltage),
		Current:   round2(current),
		Timestamp: timestamp,
	}
}

// round2 rounds a float to 2 decimal places using proper mathematical rounding.
func round2(v float64) float64 {
	return math.Round(v*100) / 100
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
