package devices

import "time"

// Device represents a user device (simple model).
type Device struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Name      string     `json:"name"`
	Room      string     `json:"room,omitempty"`
	Metadata  string     `json:"metadata,omitempty"` // JSON string; you can change to map[string]interface{} if desired
	CreatedAt time.Time  `json:"created_at"`
	LastSeen  *time.Time `json:"last_seen,omitempty"`
}
