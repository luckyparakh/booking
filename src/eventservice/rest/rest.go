package rest

import (
	"booking/src/lib/msgqueue"
	"booking/src/lib/persistence"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeApi(ep, tlsEp string, dbHandle persistence.DatabaseHandler, qHandler msgqueue.EventEmitter) (chan error, chan error) {
	httpErrChan := make(chan error)
	httpTlsErrChan := make(chan error)
	eh := NewEventHandler(dbHandle, qHandler)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	eventRouter := r.Group("/event")
	eventRouter.GET("", eh.AllEventHandler)
	eventRouter.POST("", eh.NewEventHandler)
	eventRouter.GET("/:searchCriteria/:search", eh.FindEventHandler)
	go func() {
		// Path to certs is relative to go.mod file
		// This will also work ./src/lib/certs/cert.pem
		httpTlsErrChan <- r.RunTLS(tlsEp, "src/lib/certs/cert.pem", "src/lib/certs/key.pem")
	}()
	go func() {
		httpErrChan <- r.Run(ep)
	}()
	return httpErrChan, httpTlsErrChan
}
