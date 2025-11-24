package telemetry

import "time"

// Telemetry represents a single telemetry reading from a device.
type Telemetry struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"device_id"`
	Power     float64   `json:"power"`      // Watts
	Voltage   *float64  `json:"voltage"`    // Volts (optional)
	Current   *float64  `json:"current"`    // Amps (optional)
	Timestamp time.Time `json:"timestamp"`
}

// CreateTelemetryRequest represents the payload for creating telemetry data.
type CreateTelemetryRequest struct {
	DeviceID int64    `json:"device_id" binding:"required"`
	Power    float64  `json:"power" binding:"required"`
	Voltage  *float64 `json:"voltage,omitempty"`
	Current  *float64 `json:"current,omitempty"`
}
