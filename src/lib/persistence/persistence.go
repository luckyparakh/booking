package persistence

type DatabaseHandler interface {
	AddEvent(Event) (string, error)
	FindEvent(string) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
}
