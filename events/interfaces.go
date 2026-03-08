package events

type EventInterface interface {
	GetName() string
	GetDate() string
	GetID() string
	GetPayload() any
}

type EventHandlerInterface interface {
	HandleEvent(event EventInterface) error
}

type EventDispatcherInterface interface {
	RegisterHandler(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	RemoveHandler(eventName string, handler EventHandlerInterface) error
	HasHandler(eventName string, handler EventHandlerInterface) bool
	ClearHandlers() error
}
