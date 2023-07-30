package broker

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"log"
	"os"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/internal/service/middleware"
	"wb-l0/utils"
)

type Nats struct {
	Conn    stan.Conn
	service service.IServices
	Log     zerolog.Logger
}

// NewBroker Create entity service.Broker
func NewBroker(conf *config.Config, services service.IServices) (service.Broker, error) {
	Log := zerolog.New(os.Stderr)
	nc, err := utils.TryConnectNats(conf, 15, &Log)
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "client",
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		Log.Error().Err(err).Timestamp().Msg(fmt.Sprintf(
			"Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, conf.Stan.Url))
		return nil, err
	}

	return &Nats{
		Conn:    sc,
		Log:     Log,
		service: services,
	}, nil
}

// Subscribe is subscribe on topic and execute magicFunc
func (n *Nats) Subscribe(topic string, magicFunc func(msg *stan.Msg)) error {
	_, err := n.Conn.Subscribe(topic, magicFunc, stan.DurableName("almaz"))
	if err != nil {
		_ = n.Conn.Close()
		n.Log.Error().Err(err).Timestamp()
		return &service.MyError{Code: 500, Message: err.Error()}
	}
	return nil
}

// Handler process message from stan channel
func (n *Nats) Handler(msg *stan.Msg) {
	id, err := middleware.ProcessModel(msg.Data)
	if err != nil {
		n.Log.Error().Timestamp().Err(err).Send()
		return
	}

	err = n.service.SaveModel(id, msg.Data)
	if err != nil {
		n.Log.Error().Timestamp().Err(err).Send()
		return
	}
	n.Log.Info().Str("Service", "Broker").Msgf("Broker got a new model")
}
