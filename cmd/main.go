package main

import (
	"log"
	"wb-l0/config"
	"wb-l0/internal/service/handlers/broker"
	restapi "wb-l0/internal/service/handlers/rest-api"
	"wb-l0/internal/service/usecase"
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

	services, err := usecase.NewServices(conf)
	// ???
	stan, err := broker.NewBroker(conf, services)
	if err != nil {
		log.Fatalf("can't connect stan: %v", err)
	}
	err = stan.Subscribe(conf.Stan.Topic, stan.Handler)
	if err != nil {
		log.Fatalf("can't subscribe on topic: %s. %v", conf.Stan.Topic, err)
	}

	server, err := restapi.NewServer(conf, services)
	if err != nil {
		log.Fatalf("Unable to start the rest server: %e", err)
	}
	err = server.StartServer(":3000")
	if err != nil {
		log.Fatalf("Unable to start the rest server: %e", err)
	}
}
