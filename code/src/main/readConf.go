package main

import (
	"encoding/json"

	"os"
)

type Configuration struct {
	Elasticsearch string
	Database      string
	Server        string
}

func readConf(fileName string) (elastic, dbstr, server string, err error) {

	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return
	}
	// Elasticsearch server IP and port
	elastic = configuration.Elasticsearch

	// Mongodb connection string
	dbstr = configuration.Database

	// Channel-serivce IP and port
	server = configuration.Server

	return

}
