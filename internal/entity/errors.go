package entity

import "errors"

var (
	ErrUnclosedReception  = errors.New("есть незакрытая приемка")
	ErrNoOpenReception    = errors.New("нет открытых приемок")
	ErrNoProductsToDelete = errors.New("нет товаров для удаления")
	ErrInternalServErr    = errors.New("internal server error")
)
