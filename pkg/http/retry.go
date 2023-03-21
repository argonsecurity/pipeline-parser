package http

import (
	"math"
	"net/http"
	"time"
)

type CheckForRetry func(resp *http.Response, err error) (bool, error)

func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	if resp != nil && (resp.StatusCode == 0 || resp.StatusCode >= 500) {
		return true, nil
	}
	return false, nil
}

type Backoff func(min time.Duration, max time.Duration, attemptNum int) time.Duration

func DefaultBackoff(min time.Duration, max time.Duration, attemptNum int) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		return max
	}

	return sleep
}
