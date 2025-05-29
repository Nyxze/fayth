package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Represent an ApiError from an API call (e.g: Unauthorized)
type ApiError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Param      string `json:"param"`
	Type       string `json:"type"`
	StatusCode int
	Request    *http.Request
	Response   *http.Response
}

func NewErrorFromResponse(response *http.Response) (aerror ApiError) {
	aerror.Response = response
	aerror.StatusCode = response.StatusCode
	if response.Body == nil {
		return
	}
	decoder := json.NewDecoder(response.Body)
	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			aerror.Message = err.Error()
			return
		}
		if t == "error" {
			decoder.Decode(&aerror)
			return
		}
	}
	return
}

// Error implements [error] interface
func (e ApiError) Error() string {
	return fmt.Sprintf("%s %q: %d %s\nOpenAI error: %s", e.Request.Method, e.Request.URL, e.Response.StatusCode, http.StatusText(e.Response.StatusCode), e.Message)
}
