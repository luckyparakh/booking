package rest

import (
	"booking/src/contracts"
	"booking/src/lib/msgqueue"
	"booking/src/lib/persistence"
	"booking/src/lib/persistence/mongolayer"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type eventHandler struct {
	dbLayer  persistence.DatabaseHandler
	qHandler msgqueue.EventEmitter
}

func NewEventHandler(dbLayer persistence.DatabaseHandler, qh msgqueue.EventEmitter) *eventHandler {
	return &eventHandler{
		dbLayer:  dbLayer,
		qHandler: qh,
	}
}
func (eh *eventHandler) AllEventHandler(c *gin.Context) {
	events, err := eh.dbLayer.FindAllAvailableEvents()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &events)
}

func (eh *eventHandler) NewEventHandler(c *gin.Context) {
	var event persistence.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := eh.dbLayer.AddEvent(&event, mongolayer.EVENTS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg := contracts.EventCreatedEvent{
		ID:         id,
		Name:       event.Name,
		LocationID: event.Location.ID,
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
		Capacity:   event.Capacity,
	}
	eh.qHandler.Emit(&msg)
	c.JSON(http.StatusOK, gin.H{"id": id})
}
func (eh *eventHandler) FindEventHandler(c *gin.Context) {
	sc := c.Param("searchCriteria")
	if sc == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No search criteria found"})
		return
	}
	sk := c.Param("search")
	if sk == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No search key found"})
		return
	}
	var err error
	var event persistence.Event
	switch strings.ToLower(sc) {
	case "name":
		event, err = eh.dbLayer.FindEventByName(sk)
	case "id":
		event, err = eh.dbLayer.FindEvent(sk)
		// id, er := hex.DecodeString(sk)
		// if er == nil {
		// 	event, err = eh.dbLayer.FindEvent(id)
		// }
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, &event)
}
