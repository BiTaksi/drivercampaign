package event

import (
	"errors"
)

type Type string

var (
	ErrEventUnexpectedType   = errors.New("event: unexpected event type")
	ErrHandlerUnexpectedType = errors.New("handler: unexpected event type")
)

type Event interface {
	Type() string
	Data() interface{}
}

type IEventFactory interface {
	Make(topic string, eventType string, data []byte) (Event, error)
}
