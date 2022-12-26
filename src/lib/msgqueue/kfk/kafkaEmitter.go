package kfk

import (
	"booking/src/lib/msgqueue"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	kafka "github.com/segmentio/kafka-go"
	// "github.com/confluentinc/confluent-kafka-go/kafka"
)

const topic = "events"

type messageEnvelope struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}
type kfkEventEmitter struct {
	writer *kafka.Writer
}

func createTopic(broker1Address string) error {
	connection, err := kafka.Dial("tcp", broker1Address)

	if err != nil {
		log.Printf("Failed during dailing kafka: %v\n", err)
		return err
	}

	eventTopicConfig := kafka.TopicConfig{Topic: topic, NumPartitions: 1, ReplicationFactor: 1}

	err = connection.CreateTopics(eventTopicConfig)

	if err != nil {
		log.Printf("Failed during topic creation: %v\n", err)
		return err
	}
	log.Printf("Topic created successfully")
	return nil
}

func NewKfkEventEmitter(broker1Address string) (*kfkEventEmitter, error) {
	r := retrier.New(retrier.ExponentialBackoff(4, 5*time.Second), nil)
	retryErr := r.Run(func() error {
		return createTopic(broker1Address)
	})
	if retryErr != nil {
		log.Printf("Error while kafka binding: %v\n", retryErr)
		return nil, retryErr
	}

	log.Printf("NewKfkEventEmitter Broker Address: %v\n", broker1Address)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})
	// w:=&kafka.Writer{
	// 	Addr: kafka.TCP(broker1Address),
	// 	Topic: topic,
	// }
	return &kfkEventEmitter{
		writer: w,
	}, nil
}

func (a *kfkEventEmitter) Emit(event msgqueue.Event) error {
	payload := &messageEnvelope{
		EventName: event.EventName(),
		Payload:   event,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed during converting event into byte: %v\n", err)
		return err
	}
	err = a.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("Event Published:"),
		Value: []byte(data),
	})
	if err != nil {
		log.Printf("Failed during writing into %v topic: %v\n", topic, err)
		return err
	}
	log.Printf("Successfully added message %v into %v topic.\n", event, topic)
	return nil
}
