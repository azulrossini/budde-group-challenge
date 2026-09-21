package apperr

type Code string

const (
	CodeValidationError Code = "VALIDATION_ERROR"
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeNotFound        Code = "NOT_FOUND"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeInternal        Code = "INTERNAL_ERROR"
)

const MsgInternalError = "internal error"

type FieldError struct {
	Field   string
	Message string
}

type Error struct {
	Code    Code
	Message string
	Details []FieldError
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func NotFound(message string) *Error {
	return New(CodeNotFound, message)
}

func BadRequest(message string) *Error {
	return New(CodeBadRequest, message)
}

func Unauthorized(message string) *Error {
	return New(CodeUnauthorized, message)
}

func Validation(message string, details []FieldError) *Error {
	return &Error{Code: CodeValidationError, Message: message, Details: details}
}

func Internal(err error) *Error {
	return &Error{Code: CodeInternal, Message: MsgInternalError, Err: err}
}
