package app_errors

import "encoding/json"

// Used to properly encapsulate an error withing a JSON object
type ErrorResponse struct {
	Error interface{} `json:"error"`
}

func (er *ErrorResponse) ToJSON() ([]byte, error) {
	return json.Marshal(er)
}

func WrapError(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err,
	}
}
