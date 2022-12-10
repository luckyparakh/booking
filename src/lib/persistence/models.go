package persistence

type Event struct {
	// ID        string `bson:"_id,omitempty" json:"id,omitempty"` // Should not be taken as user input
	Name      string
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}

type Location struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"` //Should not be taken as user input
	Name      string
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
