package tapo

import (
	"context"
	"fmt"
	"strings"

	tapo "github.com/tess1o/tapo-go"
)

// Connection holds credentials to reach a Tapo device.
type Connection struct {
	IP       string
	Username string
	Password string
}

// createClient creates a SmartPlug client with enhanced error handling.
// Returns a more detailed error message to help with troubleshooting.
func createClient(ctx context.Context, conn Connection) (*tapo.SmartPlug, error) {
	if conn.IP == "" {
		return nil, fmt.Errorf("device IP is required")
	}
	if conn.Username == "" {
		return nil, fmt.Errorf("username/email is required (use your Tapo Cloud account email)")
	}
	if conn.Password == "" {
		return nil, fmt.Errorf("password is required (use your Tapo Cloud account password)")
	}

	// Create client with default retry config
	sp, err := tapo.NewSmartPlug(ctx, conn.IP, conn.Username, conn.Password, tapo.Options{
		RetryConfig: tapo.DefaultRetryConfig,
	})
	
	if err != nil {
		// Provide helpful error messages based on common issues
		errMsg := err.Error()
		if strings.Contains(errMsg, "403") || strings.Contains(errMsg, "handshake") {
			return nil, fmt.Errorf("authentication failed (403): verify credentials and try these fixes:\n" +
				"1. Use your Tapo Cloud EMAIL (not username) in the 'username' field\n" +
				"2. Use your Tapo Cloud PASSWORD (the one you use in the mobile app)\n" +
				"3. Ensure device firmware is up to date\n" +
				"4. Check if device IP '%s' is reachable from the backend\n" +
				"Original error: %v", conn.IP, err)
		}
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "connection refused") {
			return nil, fmt.Errorf("cannot reach device at IP '%s': check network connectivity and firewall settings. Error: %v", conn.IP, err)
		}
		return nil, fmt.Errorf("failed to create Tapo client: %w", err)
	}

	return sp, nil
}

// ReadPower reads the current power consumption (Watts) from the device.
// Uses tapo-go's SmartPlug client to fetch current power.
func ReadPower(ctx context.Context, conn Connection) (float64, error) {
	sp, err := createClient(ctx, conn)
	if err != nil {
		return 0, err
	}

	cp, err := sp.GetCurrentPower(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to read power from device: %w", err)
	}

	// CurrentPower.Result.CurrentPower is documented as an int value.
	// Return as float64 for consistency with telemetry storage.
	return float64(cp.Result.CurrentPower), nil
}

// SetPower sets the relay state (on/off) on the device using tapo-go.
func SetPower(ctx context.Context, conn Connection, on bool) error {
	sp, err := createClient(ctx, conn)
	if err != nil {
		return err
	}

	if on {
		if _, err := sp.TurnOn(ctx); err != nil {
			return fmt.Errorf("failed to turn device ON: %w", err)
		}
		return nil
	}

	if _, err := sp.TurnOff(ctx); err != nil {
		return fmt.Errorf("failed to turn device OFF: %w", err)
	}
	return nil
}
