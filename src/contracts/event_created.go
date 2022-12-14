package contracts

import "time"

type EventCreatedEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LocationID string    `json:"location_id"`
	Start      time.Time `json:"start_date"`
	End        time.Time `json:"end_date"`
	Capacity   int       `json:"capacity"`
}

// Needed, in order to make 'EventCreatedEvent' implements Events interface
func (c *EventCreatedEvent) EventName() string {
	return "eventCreated"
}
