package main

import (
	"log"
	"wb-l0/config"
	restapi "wb-l0/internal/service/handlers/rest-api"
)

func main() {
	viperConf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.ParseConfig(viperConf)
	if err != nil {
		log.Fatal(err)
	}

	server, err := restapi.NewServer(conf)
	if err != nil {
		log.Fatalf("Unable to start the rest server: %e", err)
	}
	err = server.StartServer("8080")
	if err != nil {
		log.Fatalf("Unable to start the rest server: %e", err)
	}
}
