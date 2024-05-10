package migrator

import "errors"

var (
	ErrInvalidDriver             = errors.New("not a valid driver")
	ErrFailedToCreateConnection  = errors.New("failed to create connection db")
	ErrFailedToConnectDB         = errors.New("failed to connect db")
	ErrInvalidTablename          = errors.New("invalid tablename")
	ErrParseRow                  = errors.New("cant parse row")
	ErrFailedToMarshal           = errors.New("failed to marshal")
	ErrFailedToSendKafkaMessages = errors.New("failed to send kafka messages")
)
