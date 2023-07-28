package rest_api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"os"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/internal/service/middleware"
)

type FiberServer struct {
	server  *fiber.App
	service service.IServices
	Log     zerolog.Logger
	broker  service.Broker
}

// fix nats

func NewServer(conf *config.Config, services service.IServices) (service.RestServer, error) {
	_ = conf // on future
	s := FiberServer{}
	//var err error

	s.Log = zerolog.New(os.Stderr)
	//s.service, err = usecase.NewServices(conf)
	//if err != nil {
	//	return nil, err
	//}
	s.service = services
	s.server = fiber.New()
	s.server.Get("getModel/:id", s.GetModel)
	s.server.Post("addModel", s.AddModel)
	return &s, nil
}

func (s *FiberServer) StartServer(port string) error {
	err := s.server.Listen(port)
	if err != nil {
		s.Log.Error().Timestamp().Err(err).Send()
		return &service.MyError{Code: 500, Message: err.Error()}
	}
	s.Log.Info().Str("Service", "Server api").Msg("Server was active")
	return nil
}

func (s *FiberServer) GetModel(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	model, err := s.service.GetModel(id)
	if err != nil {
		return err
	} else {
		err = ctx.Send(model)
		if err != nil {
			return &service.MyError{Code: 500, Message: err.Error()}
		}
	}
	s.Log.Info().Str("Service", "Server api").Msgf("Server received a request to search for a model: %s", id)
	return nil
}

func (s *FiberServer) AddModel(ctx *fiber.Ctx) error {
	object := ctx.Body()
	id, err := middleware.CheckModel(object)
	if err != nil {
		s.Log.Error().Timestamp().Err(err).Send()
		return err
	}
	err = s.service.SaveModel(id, object)
	if err != nil {
		s.Log.Error().Timestamp().Err(err).Send()
		return err
	}
	s.Log.Info().Str("Service", "Server api").Msgf("Server received a request to add model: %s", id)
	return nil
}
