package checker

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/meysam81/x/logging"
)

type HTTPChecker struct {
	client     *http.Client
	upstream   string
	retries    uint
	statusCode uint
	logger     *logging.Logger
}

type CheckResult struct {
	Success  bool
	Duration time.Duration
	Status   string
	Error    error
}

func NewHTTPChecker(upstream string, timeout, retries, statusCode uint, logger *logging.Logger) *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		upstream:   upstream,
		retries:    retries,
		statusCode: statusCode,
		logger:     logger,
	}
}

func (h *HTTPChecker) Check(ctx context.Context) error {
	var attempts uint

	for {
		if h.retries > 0 && attempts >= h.retries {
			h.logger.Error().Msgf("max retries reached: %d", h.retries)
			return fmt.Errorf("max retries (%d) exceeded", h.retries)
		}

		result := h.performSingleCheck(ctx)
		attempts++

		if result.Success {
			h.logger.Info().Msgf("check successful in %s with status: %s",
				result.Duration.Round(time.Millisecond), result.Status)
			return nil
		}

		jitterSeconds := rand.Intn(6) + 5 // 5-10 seconds
		if result.Error != nil {
			h.logger.Info().Err(result.Error).Msgf("[%d] check failed, retrying in %ds...", attempts, jitterSeconds)
		}

		if err := h.waitWithJitter(ctx, jitterSeconds); err != nil {
			return err
		}
	}
}

func (h *HTTPChecker) performSingleCheck(ctx context.Context) CheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.upstream, nil)
	if err != nil {
		return CheckResult{
			Success:  false,
			Duration: time.Since(start),
			Error:    fmt.Errorf("failed to create request: %w", err),
		}
	}

	resp, err := h.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			Success:  false,
			Duration: duration,
			Error:    fmt.Errorf("request failed: %w", err),
		}
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			h.logger.Warn().Err(closeErr).Msg("failed to close response body")
		}
	}()

	success := resp.StatusCode == int(h.statusCode)
	return CheckResult{
		Success:  success,
		Duration: duration,
		Status:   resp.Status,
		Error:    nil,
	}
}

func (h *HTTPChecker) waitWithJitter(ctx context.Context, jitterSeconds int) error {
	select {
	case <-ctx.Done():
		h.logger.Info().Msg("shutdown signal received")
		return ctx.Err()
	case <-time.After(time.Duration(jitterSeconds) * time.Second):
		return nil
	}
}
