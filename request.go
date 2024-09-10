package begger

import (
	"bytes"
	"io"
	"net/http"
	"time"
	// "github.com/sirupsen/logrus"
)

type Request struct {
	Client     *http.Client
	Components RequestComponents
	Retry      *RetryOptions
	// Log        *logrus.Logger
}

// WARNING: This method is not responsible for closing the `response body`.
func (r *Request) Do() (*http.Response, *Error) {
	url := r.Components.Url.Get()
	// if r.Log != nil {
	// 	r.Log.Infof("Url: %s", url)
	// }

	var maxRetry int
	var waitInterval time.Duration
	var backoffRate float64
	if r.Retry != nil {
		maxRetry = r.Retry.MaxAttemptValue()
		waitInterval = r.Retry.WaitIntervalValue()
		backoffRate = r.Retry.BackoffRateValue()
	}
	// if r.Log != nil {
	// 	r.Log.Infof(
	// 		"Retry details -> MaxRetry: %d | WaitInterval: %d | BackoffRate: %f",
	// 		maxRetry, waitInterval, backoffRate,
	// 	)
	// }

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
	for attempt := 0; attempt < maxRetry+1; attempt++ {
		if attempt != 0 {
			time.Sleep(waitBeforeRetry)
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
