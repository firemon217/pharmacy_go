package medicine

import "errors"

var (
	ErrNotFound        = errors.New("medicine not found")
	ErrInvalidMedicine = errors.New("invalid medicine data")
	ErrInvalidCount    = errors.New("count must be greater than zero")
	ErrOutOfStock      = errors.New("not enough medicine in stock")
)
