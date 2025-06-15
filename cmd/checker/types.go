package checker

import (
	"net/http"
	"time"

	"github.com/meysam81/x/logging"
)

type HTTPCommon struct {
	HTTPClient *http.Client
	Retries    uint
	StatusCode int
	JitterMin  int
	JitterMax  int
	Logger     *logging.Logger
}

type checkResult struct {
	Success  bool
	Duration time.Duration
	Status   string
	Error    error
}
