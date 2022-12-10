package mqlayer

import (
	"booking/src/lib/msgqueue"
	"booking/src/lib/msgqueue/ampq"
)

func NewMqLayerEmitter(mqProvider, conn string) (msgqueue.EventEmitter, error) {
	switch mqProvider {
	case "rmq":
		return ampq.NewAMQPEventEmitter(conn)
	}
	return nil, nil
}
