package telemetry

import (
	"context"
	"time"
)

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

// Repo provides database operations for telemetry.
type Repo struct {
	q RowsQuerier
}

// NewRepo creates a new telemetry repository.
func NewRepo(q RowsQuerier) *Repo {
	return &Repo{q: q}
}

// Create inserts a new telemetry record and returns its ID.
func (r *Repo) Create(ctx context.Context, t *Telemetry) (int64, error) {
	sql := `INSERT INTO telemetry (device_id, power, voltage, current, timestamp)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	ts := t.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	err := r.q.QueryRow(ctx, sql, t.DeviceID, t.Power, t.Voltage, t.Current, ts).Scan(&id)
	return id, err
}

// ListByDevice returns telemetry records for a device, ordered by timestamp desc.
func (r *Repo) ListByDevice(ctx context.Context, deviceID int64, limit int) ([]Telemetry, error) {
	if limit <= 0 {
		limit = 100
	}
	sql := `SELECT id, device_id, power, voltage, current, timestamp
			FROM telemetry
			WHERE device_id = $1
			ORDER BY timestamp DESC
			LIMIT $2`
	rws, err := r.q.Query(ctx, sql, deviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rws.Close()

	var out []Telemetry
	for rws.Next() {
		var t Telemetry
		if err := rws.Scan(&t.ID, &t.DeviceID, &t.Power, &t.Voltage, &t.Current, &t.Timestamp); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rws.Err()
}

// ListByUser returns telemetry records for all devices owned by a user.
func (r *Repo) ListByUser(ctx context.Context, userID int64, limit int) ([]Telemetry, error) {
	if limit <= 0 {
		limit = 100
	}
	sql := `SELECT t.id, t.device_id, t.power, t.voltage, t.current, t.timestamp
			FROM telemetry t
			JOIN device d ON t.device_id = d.id
			WHERE d.user_id = $1
			ORDER BY t.timestamp DESC
			LIMIT $2`
	rws, err := r.q.Query(ctx, sql, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rws.Close()

	var out []Telemetry
	for rws.Next() {
		var t Telemetry
		if err := rws.Scan(&t.ID, &t.DeviceID, &t.Power, &t.Voltage, &t.Current, &t.Timestamp); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rws.Err()
}

// Delete removes a telemetry record by ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	return r.q.Exec(ctx, `DELETE FROM telemetry WHERE id = $1`, id)
}
