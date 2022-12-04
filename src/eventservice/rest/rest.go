package rest

import (
	"booking/src/lib/persistence"

	"github.com/gin-gonic/gin"
)

func ServeApi(ep, tlsEp string, dbHandle persistence.DatabaseHandler) (chan error, chan error) {
	httpErrChan := make(chan error)
	httpTlsErrChan := make(chan error)
	eh := NewEventHandler(dbHandle)
	r := gin.Default()
	eventRouter := r.Group("/event")
	eventRouter.GET("", eh.AllEventHandler)
	eventRouter.POST("", eh.NewEventHandler)
	eventRouter.GET("/:searchCriteria/:search", eh.FindEventHandler)
	go func() {
		httpTlsErrChan <- r.RunTLS(tlsEp, "../../lib/certs/cert.pem", "../../lib/certs/key.pem")
	}()
	go func() {
		httpErrChan <- r.Run(ep)
	}()
	return httpErrChan, httpTlsErrChan
}
