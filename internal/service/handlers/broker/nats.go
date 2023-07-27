package broker

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"log"
	"os"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/internal/service/middleware"
)

type Nats struct {
	Conn    stan.Conn
	service service.IServices
	Log     zerolog.Logger
}

func NewBroker(conf *config.Config, services service.IServices) (service.Broker, error) {
	Log := zerolog.New(os.Stderr)
	nc, err := nats.Connect(conf.Stan.Url)
	if err != nil {
		Log.Error().Err(err).Timestamp()
		return nil, err
	}
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

func (n *Nats) Subscribe(topic string, magicFunc func(msg *stan.Msg)) error {
	_, err := n.Conn.Subscribe(topic, magicFunc, stan.StartWithLastReceived())
	if err != nil {
		_ = n.Conn.Close()
		n.Log.Error().Err(err).Timestamp()
		return err
	}
	return nil
}

func (n *Nats) Handler(msg *stan.Msg) {
	id, err := middleware.CheckModel(msg.Data)
	if err != nil {
		n.Log.Error().Timestamp().Err(err).Send()
		return
	}

	err = n.service.SaveModel(id, msg.Data)
	if err != nil {
		n.Log.Error().Timestamp().Err(err).Send()
		return
	}
}

//func NewBroker1() {
//	nc, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer nc.Close()
//
//	sc, err := stan.Connect("test-cluster", "client",
//		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
//			log.Fatalf("Connection lost, reason: %v", reason)
//		}))
//	if err != nil {
//		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, nats.DefaultURL)
//	}
//
//	//wg := sync.WaitGroup{}
//	//wg.Add(1)
//	//wg.Wait()
//
//	_, err = sc.Subscribe("updates", func(msg *stan.Msg) {
//		fmt.Printf("Received a message: %s\n", string(msg.Data))
//	}, stan.StartWithLastReceived())
//	if err != nil {
//		sc.Close()
//		log.Fatal(err)
//	}
//	//wg.Done()
//	//
//	//sub.Unsubscribe()
//	//sc.Close()
//
//}
