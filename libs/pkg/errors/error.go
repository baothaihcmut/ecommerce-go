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

func NewError(err error, stackTrace string) error {
	return &Error{
		message:    err.Error(),
		stackTrace: stackTrace,
	}
}

func CaptureStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
