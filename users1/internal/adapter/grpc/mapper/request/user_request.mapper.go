package request

type UserRequestMapper interface {
}

type UserRequestMapperImpl struct {
}

func NewUserRequestMapper() UserRequestMapper {
	return &UserRequestMapperImpl{}
}
