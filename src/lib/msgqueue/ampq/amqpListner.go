package ampq

import (
	"booking/src/contracts"
	"booking/src/lib"
	"booking/src/lib/msgqueue"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/streadway/amqp"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
}

func NewAmqpEventListener(conn string, q string) (msgqueue.EventListener, error) {
	r := lib.Retry(connectAmqp, 5, 1)
	amqpConn, err := r(conn)
	if err != nil {
		log.Println("Error while connecting to AMPQ")
		return nil, err
	}
	listener := &amqpEventListener{
		connection: amqpConn,
		queue:      q,
	}
	err = listener.setup()
	if err != nil {
		return nil, err
	}

	return listener, err
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		log.Println("Error while creating channel for listener's setup method")
		return err
	}
	defer channel.Close()
	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err
}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		log.Println("Error while creating channel for listener's Listen method")
		return nil, nil, err
	}
	// defer channel.Close()
	for _, event := range eventNames {
		r := retrier.New(retrier.ConstantBackoff(3, 10*time.Second), nil)
		retryErr := r.Run(func() error {
			return channel.QueueBind(a.queue, event, "events", false, nil)
		})
		if retryErr != nil {
			log.Printf("Error while binding: %v\n", retryErr)
			return nil, nil, err
		}
		// if err := channel.QueueBind(a.queue, event, "events", false, nil); err != nil {
		// 	log.Printf("Error while consuming data from channel: %v\n", err)
		// 	return nil, nil, err
		// }
	}
	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		log.Println("Error while consuming data from channel")
		return nil, nil, err
	}
	events := make(chan msgqueue.Event)
	errors := make(chan error)
	go func() {
		var event msgqueue.Event
		eventNameHeader := "x-event-name"
		for msg := range msgs {
			rawEventName, ok := msg.Headers[eventNameHeader]
			if !ok {
				errors <- fmt.Errorf("message did not contain %s header", eventNameHeader)
				msg.Nack(false, false)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("%s is not string but it of %t type", eventNameHeader, rawEventName)
				msg.Nack(false, false)
				continue
			}
			switch eventName {
			case "eventCreated":
				event = &contracts.EventCreatedEvent{}
			default:
				errors <- fmt.Errorf("event type %s is not known", eventName)
				continue
			}
			if err := json.Unmarshal(msg.Body, event); err != nil {
				errors <- err
				continue
			}
			events <- event
			msg.Ack(false)
		}
	}()
	return events, errors, nil
}
