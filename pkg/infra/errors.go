package infra

type ERR string

const (
	INTERNAL_ERR   ERR = "INTERNAL_ERR"
	VALIDATION_ERR ERR = "VALIDATION_ERR"
)

type ValidationError struct {
	message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

func (e *ValidationError) Error() string {
	return e.message
}
