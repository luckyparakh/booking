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
	RestfulEPDefault    = ":8181"                    // localhost:8181 does not work with docker
	RestfulTLSEPDefault = ":9191"
	RestfulEPDefaultBk  = ":8182" // localhost:8181 does not work with docker
	Qtype               = "rmq"
	QEndpoint           = os.Getenv("AMQP_BROKER_URL")
	QName               = "event"
)

type ServiceConfig struct {
	Databasetype       dblayer.DBTYPE `json:"databasetype"`
	DBConnection       string         `json:"dbconnection"`
	RestfulEndpoint    string         `json:"restfulapi_endpoint"`
	TlsRestfulEndpoint string         `json:"restfulapi_tlsendpoint"`
	RestfulEndpointBk  string         `json:"restfulapi_endpoint_bk"`
	Qtype              string         `json:"q_type"`
	QEndpoint          string         `json:"q_endpoint"`
	QName              string         `json:"q_name"`
}

func ExtractConfiguration(fp string) (*ServiceConfig, error) {
	config := &ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		RestfulEPDefaultBk,
		Qtype,
		QEndpoint,
		QName,
	}
	configFile, err := os.Open(fp)
	if err != nil {
		fmt.Println("Error opening configuration file")
		return config, err
	}
	err = json.NewDecoder(configFile).Decode(config)
	return config, err
}
