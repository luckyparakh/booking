package rest

import (
	"booking/src/lib/persistence"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type eventHandler struct {
	dbLayer persistence.DatabaseHandler
}

func NewEventHandler(dbLayer persistence.DatabaseHandler) *eventHandler {
	return &eventHandler{
		dbLayer: dbLayer,
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
	id, err := eh.dbLayer.AddEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
