package begger

import (
	"time"
)

type RetryOptions struct {
	// Max number of retry, if the request fails.
	MaxAttempt int

	/*
		Time to wait between two successive requests. For exponential backoff,
		it's the wait interval in the first attempt.
	*/
	WaitInterval time.Duration

	/*
		Exponential backoff. This value will be multiplied with the `WaitInterval`
		after each retry. And the updated `WaitInterval` will be used in the
		next request. Value of `1` indicates simple retry with fixed wait interval.
	*/
	BackoffRate *float64
}

func (ro *RetryOptions) MaxAttemptValue() int {
	if ro.MaxAttempt < 0 {
		return 0
	}
	return ro.MaxAttempt
}

func (ro *RetryOptions) WaitIntervalValue() time.Duration {
	if ro.WaitInterval < 0 {
		return 0
	}
	return ro.WaitInterval
}

func (ro *RetryOptions) BackoffRateValue() float64 {
	if ro.BackoffRate == nil || *ro.BackoffRate <= float64(0) {
		return float64(1)
	}
	return *ro.BackoffRate
}
