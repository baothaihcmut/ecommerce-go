package exceptions

import "errors"

var (
	ErrParentCategoryNotExist = errors.New("parent category not exist")
)
