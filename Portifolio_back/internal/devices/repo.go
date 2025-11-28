package devices

import (
	"context"
	"errors"
)

// ErrNotFound is returned when a device is not found.
var ErrNotFound = errors.New("device not found")

// Querier is an interface for database operations (compatible with pgxpool.Pool wrapper).
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) error
	QueryRow(ctx context.Context, sql string, args ...any) interface{ Scan(dest ...any) error }
}

// RowsQuerier extends Querier with Query support for listing.
type RowsQuerier interface {
	Querier
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

// Rows is an interface for query result iteration.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

// Repo provides database operations for devices.
type Repo struct {
	q RowsQuerier
}

// NewRepo creates a new device repository.
func NewRepo(q RowsQuerier) *Repo {
	return &Repo{q: q}
}

// Create inserts a new device and returns its ID.
func (r *Repo) Create(ctx context.Context, d *Device) (int64, error) {
	sql := `INSERT INTO device (user_id, name, room, type, status, power_state, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int64
	status := d.Status
	if status == "" {
		status = "offline"
	}
	powerState := false // Default to off
	if d.PowerState != nil {
		powerState = *d.PowerState
	}
	err := r.q.QueryRow(ctx, sql, d.UserID, d.Name, d.Room, d.Type, status, powerState, d.Metadata).Scan(&id)
	return id, err
}

// GetByID returns a device by ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*Device, error) {
	sql := `SELECT id, user_id, name, room, type, status, power_state, metadata, created_at, last_seen
			FROM device WHERE id = $1`
	var d Device
	err := r.q.QueryRow(ctx, sql, id).Scan(
		&d.ID, &d.UserID, &d.Name, &d.Room, &d.Type, &d.Status, &d.PowerState, &d.Metadata, &d.CreatedAt, &d.LastSeen,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// ListByUser returns all devices for a user.
func (r *Repo) ListByUser(ctx context.Context, userID int64) ([]Device, error) {
	sql := `SELECT id, user_id, name, room, type, status, power_state, metadata, created_at, last_seen
			FROM device WHERE user_id = $1 ORDER BY created_at DESC`
	rws, err := r.q.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rws.Close()

	var out []Device
	for rws.Next() {
		var d Device
		if err := rws.Scan(&d.ID, &d.UserID, &d.Name, &d.Room, &d.Type, &d.Status, &d.PowerState, &d.Metadata, &d.CreatedAt, &d.LastSeen); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rws.Err()
}

// Update updates a device's fields.
func (r *Repo) Update(ctx context.Context, userID, deviceID int64, req *UpdateDeviceRequest) error {
	sql := `UPDATE device SET
			name = COALESCE($3, name),
			room = COALESCE($4, room),
			type = COALESCE($5, type),
			status = COALESCE($6, status),
			power_state = COALESCE($7, power_state),
			metadata = COALESCE($8, metadata)
			WHERE id = $1 AND user_id = $2`
	return r.q.Exec(ctx, sql, deviceID, userID, req.Name, req.Room, req.Type, req.Status, req.PowerState, req.Metadata)
}

// UpdateLastSeen updates the last_seen timestamp.
func (r *Repo) UpdateLastSeen(ctx context.Context, deviceID int64) error {
	sql := `UPDATE device SET last_seen = NOW() WHERE id = $1`
	return r.q.Exec(ctx, sql, deviceID)
}

// UpdatePowerState updates the power_state of a device.
func (r *Repo) UpdatePowerState(ctx context.Context, deviceID int64, powerState bool) error {
	sql := `UPDATE device SET power_state = $2 WHERE id = $1`
	return r.q.Exec(ctx, sql, deviceID, powerState)
}

// UpdateStatus updates the status (online/offline) of a device.
func (r *Repo) UpdateStatus(ctx context.Context, deviceID int64, status string) error {
	sql := `UPDATE device SET status = $2 WHERE id = $1`
	return r.q.Exec(ctx, sql, deviceID, status)
}

// Delete removes a device by ID for a user.
func (r *Repo) Delete(ctx context.Context, userID, deviceID int64) error {
	sql := `DELETE FROM device WHERE id = $1 AND user_id = $2`
	return r.q.Exec(ctx, sql, deviceID, userID)
}
