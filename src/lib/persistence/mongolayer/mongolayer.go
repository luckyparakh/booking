package mongolayer

import (
	"booking/src/lib/persistence"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoLayer(url string) (*MongoDBLayer, error) {
	s, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MongoDBLayer{
		session: s,
	}, nil
}

func (m *MongoDBLayer) AddEvent(e persistence.Event) ([]byte, error) {
	s := m.getFreshSession()
	s.Close()
	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}
func (m *MongoDBLayer) FindEvent(id []byte) (persistence.Event, error) {
	s := m.getFreshSession()
	s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}
func (m *MongoDBLayer) FindEventByName(name string) (persistence.Event, error) {
	s := m.getFreshSession()
	s.Close()
	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}
func (m *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	s := m.getFreshSession()
	s.Close()
	events := []persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events)
	return events, err
}
func (m *MongoDBLayer) getFreshSession() *mgo.Session {
	return m.session.Copy()
}
