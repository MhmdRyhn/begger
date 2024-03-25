package begger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type responseParser struct {
	response *http.Response
}

func NewResponseParser(response *http.Response) *responseParser {
	return &responseParser{response: response}
}

//
// If `response` is nil, this method will return "0" as status code.
//
func (rp *responseParser) HTTPStatusCode() int {
	if rp.response == nil {
		return 0
	}
	return rp.response.StatusCode
}

/*
=
This method's responsibility is to extract the response body from `http.Response`
object and load into the "loadingObject". It is NOT responsible for closing the
`response.Body`.
"loadingObject" is usually a reference to an struct object where the struct
members contain "json" tag, or it can be a map object.
=
*/
func (rp *responseParser) LoadBody(loadingObject any) *Error {
	if rp.response == nil {
		return &Error{
			HTTPStatusCode: http.StatusInternalServerError,
			StatusName:     http.StatusText(http.StatusInternalServerError),
			Message:        "Nil response.",
		}
	}
	responseBytes, err := io.ReadAll(rp.response.Body)
	if err != nil {
		return &Error{
			HTTPStatusCode: http.StatusInternalServerError,
			StatusName:     http.StatusText(http.StatusInternalServerError),
			Message:        fmt.Sprintf("Cannot decode response data into bytes: %s", err.Error()),
		}
	}
	err = json.Unmarshal(responseBytes, loadingObject)
	if err != nil {
		return &Error{
			HTTPStatusCode: http.StatusInternalServerError,
			StatusName:     http.StatusText(http.StatusInternalServerError),
			Message:        err.Error(),
		}
	}
	return nil
}

// TODO: Method to parse response headers will be added here later.
