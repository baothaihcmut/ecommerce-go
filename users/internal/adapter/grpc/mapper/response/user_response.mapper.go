package response

type UserResponseMapper interface {
}

type UserResponseMapperImpl struct {
}

func NewUserResponseMapper() UserResponseMapper {
	return &UserResponseMapperImpl{}
}
