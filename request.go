package begger

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Request struct {
	Client     *http.Client
	Components RequestComponents
	Retry      *RetryOptions

	// If this is left as nil, the default is `time.Sleep` function.
	// Otherwise, Custom implementation like time.Sleep function.
	Sleeper Sleeper
}

// WARNING: This method is not responsible for closing the `response body`.
func (r *Request) Do() (*http.Response, *Error) {
	url := r.Components.Url.Get()

	var maxRetry int
	var waitInterval time.Duration
	var backoffRate float64
	if r.Retry != nil {
		maxRetry = r.Retry.MaxAttemptValue()
		waitInterval = r.Retry.WaitIntervalValue()
		backoffRate = r.Retry.BackoffRateValue()
	}
	// One for actual request, others for retry attemps.
	maxRetry++

	var body io.Reader
	if len(r.Components.Body) != 0 {
		body = bytes.NewBuffer(r.Components.Body)
	} else {
		body = nil
	}
	request, err := http.NewRequest(r.Components.HTTPMethod, url, body)
	for key, value := range r.Components.Headers {
		request.Header.Set(key, value)
	}

	var response *http.Response
	waitBeforeRetry := waitInterval
	for attempt := 0; attempt < maxRetry; attempt++ {
		if attempt != 0 {
			if r.Sleeper == nil {
				time.Sleep(waitBeforeRetry)
			} else {
				r.Sleeper.Sleep(waitBeforeRetry)
			}
		}
		response, err = r.Client.Do(request)
		if err == nil && response != nil {
			break
		}
		if attempt != 0 {
			// Multiply the "WaitInterval" with the "BackoffRate".
			waitBeforeRetry = time.Duration(float64(waitBeforeRetry) * backoffRate)
		}
	}
	if err != nil && response == nil {
		return nil, &Error{
			HTTPStatusCode: http.StatusInternalServerError,
			StatusName:     http.StatusText(http.StatusInternalServerError),
			Message:        err.Error(),
		}
	}
	return response, nil
}

type Sleeper interface {
	Sleep(d time.Duration)
}
