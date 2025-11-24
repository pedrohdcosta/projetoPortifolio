package devices

import (
	"context"
	"database/sql"
	"fmt"
)

// Repo provides database operations for devices.
type Repo struct {
	DB *sql.DB
}

// NewRepo creates a new Repo.
func NewRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

// Create creates a new device and returns the new ID.
func (r *Repo) Create(ctx context.Context, d *Device) (int64, error) {
	var id int64
	query := `INSERT INTO device (user_id, name, room, metadata) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRowContext(ctx, query, d.UserID, d.Name, d.Room, d.Metadata).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert device: %w", err)
	}
	return id, nil
}

// ListByUser lists devices for a user.
func (r *Repo) ListByUser(ctx context.Context, userID int64) ([]Device, error) {
	query := `SELECT id, user_id, name, room, metadata, created_at, last_seen FROM device WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Device
	for rows.Next() {
		var d Device
		var metadata sql.NullString
		var lastSeen sql.NullTime
		if err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.Room, &metadata, &d.CreatedAt, &lastSeen); err != nil {
			return nil, err
		}
		if metadata.Valid {
			d.Metadata = metadata.String
		}
		if lastSeen.Valid {
			t := lastSeen.Time
			d.LastSeen = &t
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// Delete deletes a device by id and user (returns sql.ErrNoRows if not found).
func (r *Repo) Delete(ctx context.Context, userID, deviceID int64) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM device WHERE id = $1 AND user_id = $2`, deviceID, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
