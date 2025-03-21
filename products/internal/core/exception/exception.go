package exception

import "errors"

var (
	ErrUserNotShopOwnerActive = errors.New("user haven't active shop owner yet")
	ErrShopNotExist           = errors.New("shop not exist")
	ErrUserIsNotShopOwner     = errors.New("user cannot add product to this shop")
	ErrCategoryNotExist       = errors.New("category not exist")
)
