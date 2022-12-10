package ampq

import (
	"booking/src/lib"
	"booking/src/lib/msgqueue"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type amqpEventEmitter struct {
	connection *amqp.Connection
}

func NewAMQPEventEmitter(conn string) (*amqpEventEmitter, error) {
	r := lib.Retry(connectAmqp, 5, 1)
	amqpConn, err := r(conn)
	if err != nil {
		log.Println("Error while connecting to AMPQ")
		return nil, err
	}
	emitter := &amqpEventEmitter{
		connection: amqpConn,
	}
	err = emitter.setup()
	if err != nil {
		log.Println("Error while setting up AMPQ")
		return nil, err
	}
	return emitter, nil
}

func connectAmqp(ampqUrl string) (*amqp.Connection, error) {
	return amqp.Dial(ampqUrl)
}

func (a *amqpEventEmitter) setup() error {
	// In AMQP operation are not done directly on connection. Channels are used for same.
	// They are used to multiplex several virtual connections over one actual TCP connection.
	// Channels are no thread-safe.
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	// Normally, all(many) of these options should be configurable.
	// For our example, it'll probably do.
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	return err
}

func (a *amqpEventEmitter) Emit(event msgqueue.Event) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return err
	}
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.EventName()},
		ContentType: "application/json",
		Body:        jsonBody,
	}

	err = channel.Publish("events", event.EventName(), false, false, msg)
	return err
}
