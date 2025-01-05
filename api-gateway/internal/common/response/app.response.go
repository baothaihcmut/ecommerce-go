package response

type AppResponse[T any] struct {
	Success  bool     `json:"sucess"`
	Messages []string `json:"messages"`
	Data     T        `json:"data"`
}

func InitResponse[T any](success bool, messages []string, data T) AppResponse[T] {
	return AppResponse[T]{
		Success:  success,
		Messages: messages,
		Data:     data,
	}
}
