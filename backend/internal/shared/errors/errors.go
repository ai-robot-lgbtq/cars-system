package errors

import "fmt"

const (
	CodeOK             = 0
	CodeParamInvalid   = 10001
	CodeSystemError    = 10002
	CodeUnauthorized   = 20001
	CodeTokenExpired   = 20002
	CodeForbidden      = 20003
	CodeUserNotFound   = 30001
	CodeCarNotFound    = 31001
	CodeCarAlreadySold = 31002
	CodeOrderState     = 32001
	CodeOrderTimeout   = 32002
	CodePaymentFailed  = 33001
	CodeRefundFailed   = 33002
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code =%d message=%s", e.Code, e.Message)
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
