package rest

import (
	"booking/src/lib/persistence"

	"github.com/gin-gonic/gin"
)

func ServeApi(ep string, dbHandle persistence.DatabaseHandler) error {

	eh := NewEventHandler(dbHandle)
	r := gin.Default()
	eventRouter := r.Group("/event")
	eventRouter.GET("", eh.AllEventHandler)
	eventRouter.POST("", eh.NewEventHandler)
	eventRouter.GET("/:searchCriteria/:search", eh.FindEventHandler)
	return r.Run(ep)
}
