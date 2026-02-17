package errors

import "errors"

var ErrorAlreadyRegistered = errors.New("handler already registered for this event")

func NewAlreadyRegisteredError(eventName string) error {
	return ErrorAlreadyRegistered
}
