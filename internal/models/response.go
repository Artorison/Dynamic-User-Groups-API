package models

// ResponseError represents the structure of an error response.
// @Description Error response structure.
// @example {"status": 400, "error": "Bad Request", "message": "Invalid input data"}
type ResponseError struct {
	Status  int    `json:"status,omitempty"`  // HTTP status code
	Error   string `json:"error,omitempty"`   // Error type
	Message string `json:"message,omitempty"` // Error message
}

const (
	StatusOK    = "OK"
	StatusERROR = "Error"
)

func OK(msg string) ResponseError {
	return ResponseError{
		Message: msg,
	}
}

func ResponseErr(msg string, errs ...error) ResponseError {
	if len(errs) > 0 {
		err := errs[0]
		return ResponseError{
			Error:   err.Error(),
			Message: msg,
		}
	}
	return ResponseError{
		Message: msg,
	}

}

// Response represents the structure of a successful response.
// @Description Standard response structure.
// @example {"message": "Operation successful", "data": {"key": "value"}}
type Response struct {
	Message string      `json:"message"`        // Message
	Data    interface{} `json:"data,omitempty"` // Data
}
