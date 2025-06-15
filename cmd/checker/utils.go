package checker

import (
	"context"
	"time"
)

// waitWithJitter will only return error if the context is canceled/done
func waitWithJitter(ctx context.Context, jitterSeconds int) error {
	t := time.NewTicker(time.Duration(jitterSeconds) * time.Second)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
