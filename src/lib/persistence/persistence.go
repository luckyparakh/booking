package persistence

type DatabaseHandler interface {
	AddEvent(*Event, string) (string, error)
	FindEvent(string) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllAvailableEvents() ([]Event, error)
}
