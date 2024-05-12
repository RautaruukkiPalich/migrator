package migrator

import "errors"

var (
	ErrInvalidTablename = errors.New("invalid tablename")
	ErrInvalidQuery = errors.New("invalid query")
)
