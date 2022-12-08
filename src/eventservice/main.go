package main

import (
	"booking/src/eventservice/rest"
	"booking/src/lib/configuration"
	"booking/src/lib/persistence/dblayer"
	"flag"
	"log"
)

func main() {
	// Path relative to go.mod file
	configpath := flag.String("conf", "./src/lib/configuration/config.json", "Path to config json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*configpath)
	dbh, err := dblayer.NewPersistanceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		log.Fatal("Error while connecting to persistance layer:", err)
	}
	httpErrChan, httpTlsErrChan := rest.ServeApi(config.RestfulEndpoint, config.TlsRestfulEndpoint, dbh)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error", err)
	case err := <-httpTlsErrChan:
		log.Fatal("HTTPS Error", err)
	}
}
