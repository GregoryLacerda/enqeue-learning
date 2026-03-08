package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID      string
	Name    string
	Date    time.Time
	Payload any
}

func NewEvent(name string) *Event {
	return &Event{
		ID:   uuid.NewString(),
		Name: name,
		Date: time.Now(),
	}
}

func (e *Event) GetName() string {
	return e.Name
}

func (e *Event) GetDate() string {
	return e.Date.String()
}

func (e *Event) GetID() string {
	return e.ID
}

func (e *Event) GetPayload() any {
	return e.Payload
}
