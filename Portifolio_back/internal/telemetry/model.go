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

// TelemetrySummary represents aggregated telemetry data for a period.
type TelemetrySummary struct {
	DeviceID     int64     `json:"device_id"`
	Period       string    `json:"period"`        // "day", "week", "month"
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	TotalRecords int       `json:"total_records"`
	AvgPower     float64   `json:"avg_power"`     // Average power in Watts
	MaxPower     float64   `json:"max_power"`     // Maximum power in Watts
	MinPower     float64   `json:"min_power"`     // Minimum power in Watts
	TotalEnergy  float64   `json:"total_energy"`  // Energy in kWh (estimated)
	AvgVoltage   *float64  `json:"avg_voltage"`   // Average voltage (optional)
	AvgCurrent   *float64  `json:"avg_current"`   // Average current (optional)
}

// DeviceWithTelemetry represents a device with its latest telemetry reading.
type DeviceWithTelemetry struct {
	DeviceID   int64    `json:"device_id"`
	DeviceName string   `json:"device_name"`
	Room       string   `json:"room,omitempty"`
	Status     string   `json:"status"`
	LatestPower *float64 `json:"latest_power,omitempty"`
	LastSeen   *time.Time `json:"last_seen,omitempty"`
}
