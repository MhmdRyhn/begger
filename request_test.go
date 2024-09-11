package begger

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Request_Successful(t *testing.T) {
	transport := newRoundTripper()
	transport.On("RoundTrip", mock.Anything).Return(
		&http.Response{StatusCode: http.StatusOK}, nil,
	)

	httpClient := http.Client{Transport: transport}
	request := Request{
		Client: &httpClient,
		Components: RequestComponents{
			Url:        Url{Actual: ValueToPointer("example.com")},
			HTTPMethod: http.MethodGet,
		},
	}
	response, err := request.Do()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func Test_Request_Error(t *testing.T) {
	transport := newRoundTripper()
	transport.On("RoundTrip", mock.Anything).Return(
		&http.Response{},
		errors.New("MockError"),
	)

	httpClient := http.Client{Transport: transport}
	request := Request{
		Client: &httpClient,
		Components: RequestComponents{
			Url:        Url{Actual: ValueToPointer("example.com")},
			HTTPMethod: http.MethodGet,
		},
	}
	response, err := request.Do()
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatusCode)
}

func Test_Request_SuccessfulAfterRetry(t *testing.T) {
	transport := newRoundTripper()
	transport.On("RoundTrip", mock.Anything).Return(
		nil, errors.New("MockError"),
	).Once()
	transport.On("RoundTrip", mock.Anything).Return(
		&http.Response{StatusCode: http.StatusOK}, nil,
	)

	httpClient := http.Client{Transport: transport}
	request := Request{
		Client: &httpClient,
		Components: RequestComponents{
			Url:        Url{Actual: ValueToPointer("example.com")},
			HTTPMethod: http.MethodGet,
		},
		Retry:   &RetryOptions{MaxAttempt: 3, WaitInterval: 2 * time.Second},
		Sleeper: NewSleeperMock(),
	}
	response, err := request.Do()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func Test_Request_ErrorAfterRetry(t *testing.T) {
	transport := newRoundTripper()
	transport.On("RoundTrip", mock.Anything).Return(
		nil, errors.New("MockError"),
	).Once()
	transport.On("RoundTrip", mock.Anything).Return(
		nil, errors.New("MockError"),
	).Once()
	transport.On("RoundTrip", mock.Anything).Return(
		nil, errors.New("MockError"),
	).Once()

	httpClient := http.Client{Transport: transport}
	request := Request{
		Client: &httpClient,
		Components: RequestComponents{
			Url:        Url{Actual: ValueToPointer("example.com")},
			HTTPMethod: http.MethodGet,
		},
		Retry: &RetryOptions{
			MaxAttempt:   2,
			WaitInterval: 2000 * time.Millisecond,
			BackoffRate:  ValueToPointer(float64(1.5)),
		},
		Sleeper: NewSleeperMock(),
	}
	response, err := request.Do()
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatusCode)
}
