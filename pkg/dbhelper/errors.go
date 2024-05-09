package dbhelper

import "errors"

var (
	ErrInvalidDriver            = errors.New("not a valid driver")
	ErrFailedToCreateConnection = errors.New("failed to create connection db")
	ErrFailedToConnectDB        = errors.New("failed to connect db")
)