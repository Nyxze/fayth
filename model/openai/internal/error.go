package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Represent an Error from an API call (e.g: Unauthorized)
type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Param      string `json:"param"`
	Type       string `json:"type"`
	StatusCode int
	Request    *http.Request
	Response   *http.Response
}

func NewErrorFromResponse(r io.Reader) Error {
	decoder := json.NewDecoder(r)
	apiError := Error{}
	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			apiError.Message = err.Error()
			return apiError
		}
		if t == "error" {
			decoder.Decode(&apiError)
			return apiError
		}
	}
	return apiError
}

// Error implements [error] interface
func (e Error) Error() string {
	return fmt.Sprintf("%s %q: %d %s\nOpenAI error: %s", e.Request.Method, e.Request.URL, e.Response.StatusCode, http.StatusText(e.Response.StatusCode), e.Message)
}
