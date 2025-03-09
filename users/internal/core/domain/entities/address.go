package entities

import "github.com/google/uuid"

type Address struct {
	UserId   uuid.UUID
	Priority int
	Street   string
	Town     string
	City     string
	Province string
}

func NewAddress(
	userId uuid.UUID,
	priority int,
	street string,
	town string,
	city string,
	province string,
) *Address {
	return &Address{
		UserId:   userId,
		Priority: priority,
		Street:   street,
		Town:     town,
		City:     city,
		Province: province,
	}
}
