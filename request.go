package begger

import (
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Client *http.Client
	Url    Url
	/*
		If nil, it's assumed that it's a GET request with the
		only header {"Content-Type": "application/json"}.
	*/
	RequestInfo *RequestInfo
	Retry       *RetryOptions
	Log         *logrus.Logger
}

// WARNING: This method is not responsible for closing the `response body`.
func (r *Request) Do() (*http.Response, *Error) {
	url := r.Url.Get()
	if r.Log != nil {
		r.Log.Infof("Url: %s", url)
	}

	var maxRetry int
	var waitInterval time.Duration
	var backoffRate float64
	if r.Retry != nil {
		maxRetry = r.Retry.CleanMaxAttempt()
		waitInterval = r.Retry.CleanWaitInterval()
		backoffRate = r.Retry.CleanBackoffRate()
	}
	if r.Log != nil {
		r.Log.Infof(
			"Retry details -> MaxRetry: %d | WaitInterval: %d | BackoffRate: %f",
			maxRetry, waitInterval, backoffRate,
		)
	}

	var request *http.Request
	var err error
	if r.RequestInfo != nil {
		request, err = http.NewRequest(r.RequestInfo.HTTPMethod, url, r.RequestInfo.Body)
		for key, value := range r.RequestInfo.Headers {
			request.Header.Set(key, value)
		}
	} else {
		request, err = http.NewRequest(http.MethodGet, url, nil)
		request.Header.Set("Content-Type", "application/json")
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
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
			Message:    err.Error(),
		}
	}
	return response, nil
}

type RequestInfo struct {
	HTTPMethod string
	Body       io.Reader
	Headers    Headers
}
