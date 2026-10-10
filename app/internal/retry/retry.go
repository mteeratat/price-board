package retry

import (
	"context"
	"time"
)

// Do calls fn up to attempts times, waiting baseDelay, 2x, 4x... between tries.
// It stops early if ctx ends.
func Do(ctx context.Context, attempts int, baseDelay time.Duration, fn func() error) error {
	var err error
	for i := range attempts {
		err = fn()
		if err == nil {
			return nil
		}
		if i == attempts-1 {
			break
		}

		select {
		case <-time.After(baseDelay << i):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
