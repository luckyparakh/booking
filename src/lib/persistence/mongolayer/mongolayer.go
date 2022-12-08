package mongolayer

import (
	"booking/src/lib/persistence"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

var ctx = context.TODO()

type MongoDBLayer struct {
	client *mongo.Client
}

func NewMongoLayer(url string) (*MongoDBLayer, error) {
	clientOpts := options.Client().ApplyURI(url)
	c, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDBLayer{
		client: c,
	}, nil
}

func (m *MongoDBLayer) AddEvent(e persistence.Event) (string, error) {
	result, err := m.client.Database(DB).Collection(EVENTS).InsertOne(ctx, e)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result.InsertedID), nil
}

func (m *MongoDBLayer) FindEvent(id string) (persistence.Event, error) {
	e := persistence.Event{}
	eventId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return e, err
	}
	// https://stackoverflow.com/questions/64281675/bson-d-vs-bson-m-for-find-queries
	filter := bson.M{"_id": eventId}
	err = m.client.Database(DB).Collection(EVENTS).FindOne(ctx, filter).Decode(&e)
	if err != nil {
		return e, err
	}
	return e, nil
}
func (m *MongoDBLayer) FindEventByName(name string) (persistence.Event, error) {
	e := persistence.Event{}
	filter := bson.M{"name": name}
	err := m.client.Database(DB).Collection(EVENTS).FindOne(ctx, filter).Decode(&e)
	if err != nil {
		return e, err
	}
	return e, nil
}
func (m *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	var event persistence.Event
	var events []persistence.Event
	cur, err := m.client.Database(DB).Collection(EVENTS).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, err
}
