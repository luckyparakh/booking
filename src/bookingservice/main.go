package main

import (
	"booking/src/bookingservice/listener"
	"booking/src/lib/configuration"
	"booking/src/lib/msgqueue/mqlayer"
	"booking/src/lib/persistence/dblayer"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Path relative to go.mod file
	configpath := flag.String("conf", "./src/lib/configuration/config.json", "Path to config json file")
	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*configpath)
	qh, err := mqlayer.NewMqLayerListener(config.Qtype, config.QEndpoint, config.QName)
	if err != nil {
		log.Fatal("Error while connecting MQ layer", err)
	}
	dbh, err := dblayer.NewPersistanceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		log.Fatal("Error while connecting DB layer", err)
	}
	el := listener.NewEventListener(qh, dbh)
	go el.ProcessEvent()
	r := gin.Default()
	r.Run(config.RestfulEndpointBk)
}
