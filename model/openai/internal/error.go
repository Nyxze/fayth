package internal

import "net/http"

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

// Error implements [error] interface
func (Error) Error() string {
	return ""
}
