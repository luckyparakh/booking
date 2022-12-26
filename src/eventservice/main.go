package main

import (
	"booking/src/eventservice/rest"
	"booking/src/lib/configuration"
	"booking/src/lib/msgqueue/mqlayer"
	"booking/src/lib/persistence/dblayer"
	"flag"
	"log"
)

func main() {

	// Path relative to go.mod file
	configpath := flag.String("conf", "./src/lib/configuration/config.json", "Path to config json file")
	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*configpath)
	// log.Printf("config:%v %v\n", config.Qtype, config.QEndpoint)
	qh, err := mqlayer.NewMqLayerEmitter(config.Qtype, config.QEndpoint)
	if err != nil {
		log.Fatal("Error while connecting MQ layer", err)
	}
	dbh, err := dblayer.NewPersistanceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		log.Fatal("Error while connecting to persistance layer:", err)
	}
	httpErrChan, httpTlsErrChan := rest.ServeApi(config.RestfulEndpoint, config.TlsRestfulEndpoint, dbh, qh)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error", err)
	case err := <-httpTlsErrChan:
		log.Fatal("HTTPS Error", err)
	}
}
