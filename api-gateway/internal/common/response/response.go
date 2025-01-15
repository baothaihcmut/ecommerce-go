package response

type AppResponse struct {
	Success  bool        `json:"sucess"`
	Messages []string    `json:"messages"`
	Data     interface{} `json:"data"`
}

func InitResponse(success bool, messages []string, data interface{}) AppResponse {
	return AppResponse{
		Success:  success,
		Messages: messages,
		Data:     data,
	}
}
