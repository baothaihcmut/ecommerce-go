package errors

import "runtime"

type Error struct {
	message    string
	stackTrace string
}

func (e *Error) Error() string {
	return e.message
}
func (e *Error) StackTrace() string {
	return e.stackTrace
}

func NewError(message string, stackTrace string) error {
	return &Error{
		message:    message,
		stackTrace: stackTrace,
	}
}

func captureStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
