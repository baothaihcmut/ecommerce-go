package exception

import "google.golang.org/grpc/codes"

type errorDetail struct {
	Msg  string
	Code codes.Code
}

var errorMap map[error]errorDetail = map[error]errorDetail{}

func MapException(err error) (string, codes.Code) {
	if detail, exist := errorMap[err]; exist {
		return detail.Msg, detail.Code
	}
	return "Internal server error", codes.Internal
}
