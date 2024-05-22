package broker

import "errors"

var (
	ErrFailedToMarshal           = errors.New("failed to marshal")
	ErrFailedToSendKafkaMessages = errors.New("failed to send kafka messages")
)