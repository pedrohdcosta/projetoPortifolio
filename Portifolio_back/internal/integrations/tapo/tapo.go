package tapo

import (
	"context"
	"fmt"

	tapo "github.com/tess1o/tapo-go"
)

// Connection holds credentials to reach a Tapo device.
type Connection struct {
	IP       string
	Username string
	Password string
}

// ReadPower reads the current power consumption (Watts) from the device.
// Uses tapo-go's SmartPlug client to fetch current power.
func ReadPower(ctx context.Context, conn Connection) (float64, error) {
	sp, err := tapo.NewSmartPlug(ctx, conn.IP, conn.Username, conn.Password, tapo.Options{RetryConfig: tapo.DefaultRetryConfig})
	if err != nil {
		return 0, fmt.Errorf("tapo: failed to create smartplug client: %w", err)
	}

	cp, err := sp.GetCurrentPower(ctx)
	if err != nil {
		return 0, fmt.Errorf("tapo: GetCurrentPower failed: %w", err)
	}

	// CurrentPower.Result.CurrentPower is documented as an int value.
	// Return as float64 for consistency with telemetry storage.
	return float64(cp.Result.CurrentPower), nil
}

// SetPower sets the relay state (on/off) on the device using tapo-go.
func SetPower(ctx context.Context, conn Connection, on bool) error {
	sp, err := tapo.NewSmartPlug(ctx, conn.IP, conn.Username, conn.Password, tapo.Options{RetryConfig: tapo.DefaultRetryConfig})
	if err != nil {
		return fmt.Errorf("tapo: failed to create smartplug client: %w", err)
	}

	if on {
		if _, err := sp.TurnOn(ctx); err != nil {
			return fmt.Errorf("tapo: TurnOn failed: %w", err)
		}
		return nil
	}

	if _, err := sp.TurnOff(ctx); err != nil {
		return fmt.Errorf("tapo: TurnOff failed: %w", err)
	}
	return nil
}
