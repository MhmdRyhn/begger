package begger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"
)

// Implements `http.RoundTripper` interface
type RoundTripperMock struct {
	mock.Mock
}

func newRoundTripper() *RoundTripperMock {
	return &RoundTripperMock{}
}

func (m *RoundTripperMock) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// Implements `Sleeper` interface
type SleeperMock struct {
	mock.Mock
}

func NewSleeperMock() *SleeperMock {
	return &SleeperMock{}
}

func (s *SleeperMock) Sleep(d time.Duration) {
	fmt.Println(fmt.Sprintf("Skip sleeping %+v for faster unit test execution.", d))
}
