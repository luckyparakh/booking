package main

import (
	"booking/src/eventservice/rest"
	"booking/src/lib/configuration"
	"booking/src/lib/persistence/dblayer"
	"flag"
	"log"
)

func main() {
	configpath := flag.String("conf", "../lib/configuration/config.json", "Path to config json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*configpath)
	dbh, _ := dblayer.NewPersistanceLayer(config.Databasetype, config.DBConnection)
	httpErrChan, httpTlsErrChan := rest.ServeApi(config.RestfulEndpoint, config.TlsRestfulEndpoint, dbh)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error", err)
	case err := <-httpTlsErrChan:
		log.Fatal("HTTPS Error", err)
	}
}
