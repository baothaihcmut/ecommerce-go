package exceptions

import "errors"

var (
	ErrParentCategoryNotExist = errors.New("parent category not exist")
	ErrCategoryNotExist       = errors.New("category not exist")
	ErrShopNotExist           = errors.New("shop not exist")
	ErrShopNotActive          = errors.New("shop is inactive")
)
