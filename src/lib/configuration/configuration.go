package configuration

import (
	"booking/src/lib/persistence/dblayer"
	"encoding/json"
	"fmt"
	"os"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost:8181"
)

type ServiceConfig struct {
	Databasetype    dblayer.DBTYPE `json:"databasetype"`
	DBConnection    string         `json:"dbconnection"`
	RestfulEndpoint string         `json:"restfulapi_endpoint"`
}

func ExtractConfiguration(fp string) (*ServiceConfig, error) {
	config := &ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
	}
	configFile, err := os.Open(fp)
	if err != nil {
		fmt.Println("Configuration file not found")
		return config, err
	}
	err = json.NewDecoder(configFile).Decode(config)
	return config, err
}