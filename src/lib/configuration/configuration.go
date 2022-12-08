package configuration

import (
	"booking/src/lib/persistence/dblayer"
	"encoding/json"
	"fmt"
	"os"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://eventsdb:27017" // Same as name of service in docker compose file
	RestfulEPDefault    = ":8180"                    // localhost:8180 does not work with docker
	RestfulTLSEPDefault = ":9191"
)

type ServiceConfig struct {
	Databasetype       dblayer.DBTYPE `json:"databasetype"`
	DBConnection       string         `json:"dbconnection"`
	RestfulEndpoint    string         `json:"restfulapi_endpoint"`
	TlsRestfulEndpoint string         `json:"restfulapi_tlsendpoint"`
}

func ExtractConfiguration(fp string) (*ServiceConfig, error) {
	config := &ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
	}
	configFile, err := os.Open(fp)
	if err != nil {
		fmt.Println("Error opening configuration file")
		return config, err
	}
	err = json.NewDecoder(configFile).Decode(config)
	return config, err
}
