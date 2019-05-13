package account

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrNotEnoughBalance = errors.New("not enough balance")
)
