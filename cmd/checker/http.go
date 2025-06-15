package checker

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type StaticHTTPChecker struct {
	Upstream string

	Common *HTTPCommon
}

func (h *StaticHTTPChecker) Check(ctx context.Context) error {
	var attempts uint

	for {
		if h.Common.Retries > 0 && attempts >= h.Common.Retries {
			h.Common.Logger.Error().Msgf("max retries reached: %d", h.Common.Retries)
			return fmt.Errorf("max retries (%d) exceeded", h.Common.Retries)
		}

		result := h.Common.performSingleCheck(ctx, h.Upstream)
		attempts++

		if result.Success {
			h.Common.Logger.Info().Msgf("check successful in %s with status: %s",
				result.Duration.Round(time.Millisecond), result.Status)
			return nil
		}

		jitterSeconds := rand.Intn(6) + 5 // 5-10 seconds
		if result.Error != nil {
			h.Common.Logger.Info().Err(result.Error).Msgf("[%d] check failed, retrying in %ds...", attempts, jitterSeconds)
		}

		if err := waitWithJitter(ctx, jitterSeconds); err != nil {
			h.Common.Logger.Info().Msg("shutdown signal received")
			return err
		}
	}
}

func (c *HTTPCommon) performSingleCheck(ctx context.Context, upstream string) checkResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, upstream, nil)
	if err != nil {
		return checkResult{
			Success:  false,
			Duration: time.Since(start),
			Error:    fmt.Errorf("failed to create request: %w", err),
		}
	}

	resp, err := c.HTTPClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		return checkResult{
			Success:  false,
			Duration: duration,
			Error:    fmt.Errorf("request failed: %w", err),
		}
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.Logger.Warn().Err(closeErr).Msg("failed to close response body")
		}
	}()

	success := resp.StatusCode == c.StatusCode
	return checkResult{
		Success:  success,
		Duration: duration,
		Status:   resp.Status,
		Error:    nil,
	}
}
