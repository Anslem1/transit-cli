package middleware2


// TransitError is a custom error type that includes an error code.
type TransitError struct {
	Code    int
	Message string
}

// Error implements the error interface.
func (e *TransitError) Error() string {
	return e.Message
}

// NewTransitError creates a new TransitError.
func NewTransitError(code int, message string) error {
	return &TransitError{
		Code:    code,
		Message: message,
	}
}
