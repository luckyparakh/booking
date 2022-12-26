package kfk

import (
	"booking/src/contracts"
	"booking/src/lib/msgqueue"
	"context"
	"encoding/json"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

type kfkEventListener struct {
	reader *kafka.Reader
}

const topicListen = "events"

func NewKfkEventListener(broker1Address string) (*kfkEventListener, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topicListen,
	})
	return &kfkEventListener{
		reader: r,
	}, nil
}
func (k *kfkEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	// Not using eventNames
	events := make(chan msgqueue.Event)
	errors := make(chan error)
	go func() {
		for {
			msg, err := k.reader.ReadMessage(context.Background())
			if err != nil {
				errors <- err
				continue
			}
			var body messageEnvelope
			if err := json.Unmarshal(msg.Value, &body); err != nil {
				errors <- err
				continue
			}
			var e msgqueue.Event
			switch body.EventName {
			case "eventCreated":
				e, err = convertInterfaceIntoEventCreatedEvent(body.Payload)
				if err != nil {
					errors <- err
					continue
				}

			}
			events <- e
		}

	}()
	return events, errors, nil
}

func convertInterfaceIntoEventCreatedEvent(data interface{}) (*contracts.EventCreatedEvent, error) {
	// e := &contracts.EventCreatedEvent{}
	if e, ok := data.(contracts.EventCreatedEvent); !ok {
		return nil, fmt.Errorf("cannot convert to EventCreatedEvent type")
	} else {
		return &e, nil
	}
	// m := data.(map[string]interface{})
	// if v, ok := m["id"].(string); ok {
	// 	e.ID = v
	// }
}
