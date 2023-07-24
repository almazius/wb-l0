package broker

import "wb-l0/internal/service"

type Nats struct {
	Connection byte
}

func (n *Nats) GetEvent() service.Event {
	return service.Event{}
}
