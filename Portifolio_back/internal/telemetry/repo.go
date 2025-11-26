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
	if limit > 1000 {
		limit = 1000 // Cap at 1000 to prevent resource exhaustion
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
	if limit > 1000 {
		limit = 1000 // Cap at 1000 to prevent resource exhaustion
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

// UserOwnsDevice checks if a user owns a device.
func (r *Repo) UserOwnsDevice(ctx context.Context, userID, deviceID int64) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM device WHERE id = $1 AND user_id = $2)`
	var exists bool
	err := r.q.QueryRow(ctx, sql, deviceID, userID).Scan(&exists)
	return exists, err
}

// ListByDeviceForUser returns telemetry for a device owned by a specific user.
func (r *Repo) ListByDeviceForUser(ctx context.Context, userID, deviceID int64, limit int) ([]Telemetry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	sql := `SELECT t.id, t.device_id, t.power, t.voltage, t.current, t.timestamp
			FROM telemetry t
			JOIN device d ON t.device_id = d.id
			WHERE t.device_id = $1 AND d.user_id = $2
			ORDER BY t.timestamp DESC
			LIMIT $3`
	rws, err := r.q.Query(ctx, sql, deviceID, userID, limit)
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

// UpdateDeviceLastSeenAndStatus updates the device last_seen and status when telemetry is received.
func (r *Repo) UpdateDeviceLastSeenAndStatus(ctx context.Context, deviceID int64) error {
	sql := `UPDATE device SET last_seen = NOW(), status = 'online' WHERE id = $1`
	return r.q.Exec(ctx, sql, deviceID)
}

// GetSummaryByDevice returns aggregated telemetry for a device over a time period.
func (r *Repo) GetSummaryByDevice(ctx context.Context, userID, deviceID int64, period string) (*TelemetrySummary, error) {
	var interval string
	switch period {
	case "day":
		interval = "1 day"
	case "week":
		interval = "7 days"
	case "month":
		interval = "30 days"
	default:
		interval = "1 day"
		period = "day"
	}

	sql := `SELECT 
				COUNT(*) as total_records,
				COALESCE(AVG(t.power), 0) as avg_power,
				COALESCE(MAX(t.power), 0) as max_power,
				COALESCE(MIN(t.power), 0) as min_power,
				COALESCE(AVG(t.voltage), 0) as avg_voltage,
				COALESCE(AVG(t.current), 0) as avg_current,
				MIN(t.timestamp) as start_time,
				MAX(t.timestamp) as end_time
			FROM telemetry t
			JOIN device d ON t.device_id = d.id
			WHERE t.device_id = $1 AND d.user_id = $2 AND t.timestamp >= NOW() - INTERVAL '` + interval + `'`

	var summary TelemetrySummary
	var avgVoltage, avgCurrent float64
	summary.DeviceID = deviceID
	summary.Period = period

	err := r.q.QueryRow(ctx, sql, deviceID, userID).Scan(
		&summary.TotalRecords,
		&summary.AvgPower,
		&summary.MaxPower,
		&summary.MinPower,
		&avgVoltage,
		&avgCurrent,
		&summary.StartTime,
		&summary.EndTime,
	)
	if err != nil {
		return nil, err
	}

	// Calculate estimated energy consumption (kWh)
	// Energy = Power * Time (hours)
	if summary.TotalRecords > 0 {
		duration := summary.EndTime.Sub(summary.StartTime).Hours()
		summary.TotalEnergy = (summary.AvgPower * duration) / 1000 // Convert Wh to kWh
		if avgVoltage > 0 {
			summary.AvgVoltage = &avgVoltage
		}
		if avgCurrent > 0 {
			summary.AvgCurrent = &avgCurrent
		}
	}

	return &summary, nil
}

// GetLatestByDevice returns the most recent telemetry reading for a device.
func (r *Repo) GetLatestByDevice(ctx context.Context, deviceID int64) (*Telemetry, error) {
	sql := `SELECT id, device_id, power, voltage, current, timestamp
			FROM telemetry
			WHERE device_id = $1
			ORDER BY timestamp DESC
			LIMIT 1`
	var t Telemetry
	err := r.q.QueryRow(ctx, sql, deviceID).Scan(&t.ID, &t.DeviceID, &t.Power, &t.Voltage, &t.Current, &t.Timestamp)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetLatestByUserDevices returns the latest telemetry for all devices of a user (one per device).
func (r *Repo) GetLatestByUserDevices(ctx context.Context, userID int64) ([]Telemetry, error) {
	sql := `SELECT DISTINCT ON (t.device_id) t.id, t.device_id, t.power, t.voltage, t.current, t.timestamp
			FROM telemetry t
			JOIN device d ON t.device_id = d.id
			WHERE d.user_id = $1
			ORDER BY t.device_id, t.timestamp DESC`
	rws, err := r.q.Query(ctx, sql, userID)
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
