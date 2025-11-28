package devices

import "time"

// Device represents a user IoT device.
type Device struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Name       string     `json:"name"`
	Room       string     `json:"room,omitempty"`
	Type       string     `json:"type,omitempty"`      // e.g., "smart_plug", "sensor"
	Status     string     `json:"status,omitempty"`    // "online" or "offline"
	PowerState *bool      `json:"power_state"`         // true = on, false = off
	Metadata   string     `json:"metadata,omitempty"`  // JSON string for additional config
	CreatedAt  time.Time  `json:"created_at"`
	LastSeen   *time.Time `json:"last_seen,omitempty"`
}

// CreateDeviceRequest represents the payload for creating a device.
type CreateDeviceRequest struct {
	Name     string `json:"name" binding:"required"`
	Room     string `json:"room,omitempty"`
	Type     string `json:"type,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// UpdateDeviceRequest represents the payload for updating a device.
type UpdateDeviceRequest struct {
	Name       *string `json:"name,omitempty"`
	Room       *string `json:"room,omitempty"`
	Type       *string `json:"type,omitempty"`
	Status     *string `json:"status,omitempty"`
	PowerState *bool   `json:"power_state,omitempty"`
	Metadata   *string `json:"metadata,omitempty"`
}
