package checker

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func (c *HTTPCommon) runWithJitterBackoff(ctx context.Context, do func() error) error {
	low := c.JitterMin
	high := c.JitterMax - c.JitterMin + 1

	var attempts uint

	for {
		jitterSeconds := rand.Intn(high) + low

		err := do()
		attempts++

		tries := fmt.Sprintf("%d", attempts)
		if c.Retries > 0 {
			tries = fmt.Sprintf("%d/%d", attempts, c.Retries)
		}

		if err != nil {
			c.Logger.Info().Err(err).Msgf("[%s] failed, retrying in %ds...", tries, jitterSeconds)
		} else {
			return nil
		}

		if c.Retries > 0 {
			if attempts >= c.Retries {
				c.Logger.Error().Msgf("max retries reached: %d", c.Retries)
				return fmt.Errorf("max retries (%d) exceeded", c.Retries)
			}
		}

		if err := waitWithJitter(ctx, jitterSeconds); err != nil {
			c.Logger.Info().Msg("shutdown signal received")
			return err
		}
	}
}

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
