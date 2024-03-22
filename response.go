package begger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type responseParser struct {
	response *http.Response
}

func NewResponseParser(response *http.Response) *responseParser {
	return &responseParser{response: response}
}

//
// This method's responsibility is to extract the response body from `http.Response`
// object and load into the "loadingObject". It is NOT responsible for closing the
// `response.Body`.
// "loadingObject" must be a reference to an struct object where the struct
// members contain "json" tag.
//
func (rp *responseParser) LoadBody(loadingObject any) error {
	if rp.response == nil {
		return errors.New("Nil response.")
	}
	responseBytes, err := ioutil.ReadAll(rp.response.Body)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Cannot decode response data into bytes: %s", err.Error()),
		)
	}
	err = json.Unmarshal(responseBytes, loadingObject)
	if err != nil {
		return err
	}
	return nil
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
