package listener

import (
	"booking/src/contracts"
	"booking/src/lib/msgqueue"
	"booking/src/lib/persistence"
	"booking/src/lib/persistence/mongolayer"
	"log"
)

type IEventListener interface {
	ProcessEvent() error
	CreateEvent(msgqueue.Event) error
}
type EventListener struct {
	qConn msgqueue.EventListener
	dConn persistence.DatabaseHandler
}

func NewEventListener(qConn msgqueue.EventListener, dConn persistence.DatabaseHandler) IEventListener {
	return &EventListener{
		qConn: qConn,
		dConn: dConn,
	}
}
func (e *EventListener) ProcessEvent() error {
	log.Println("Processing Events")
	rCh, eCh, err := e.qConn.Listen("eventCreated")
	if err != nil {
		log.Printf("Error while listen messages: %v\n", err)
		return err
	}
	for {
		select {
		case <-eCh:
			log.Printf("Error recevied while listening to messages:%v\n", err)
			return err
		case event := <-rCh:
			log.Printf("event %v", event)
			e.CreateEvent(event)
		}
	}
}
func (l *EventListener) CreateEvent(evt msgqueue.Event) error {
	switch e := evt.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %v", e.ID, e)
		pe := persistence.Event{
			ID:       e.ID,
			Name:     e.Name,
			Capacity: e.Capacity,
		}
		l.dConn.AddEvent(&pe, mongolayer.BOOKING)
	default:
		log.Printf("Unknown type %t", e)
	}

	return nil
}
