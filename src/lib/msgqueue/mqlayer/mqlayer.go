package mqlayer

import (
	"booking/src/lib/msgqueue"
	"booking/src/lib/msgqueue/ampq"
	"booking/src/lib/msgqueue/kfk"
)

func NewMqLayerEmitter(mqProvider, conn string) (msgqueue.EventEmitter, error) {
	switch mqProvider {
	case "rmq":
		return ampq.NewAMQPEventEmitter(conn)
	case "kafka":
		return kfk.NewKfkEventEmitter(conn)
	}
	return nil, nil
}

func NewMqLayerListener(mqProvider, conn string) (msgqueue.EventListener, error) {
	switch mqProvider {
	case "rmq":
		return ampq.NewAmqpEventListener(conn)
	case "kafka":
		return kfk.NewKfkEventListener(conn)
	}

	return nil, nil
}
