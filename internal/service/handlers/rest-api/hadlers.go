package rest_api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"wb-l0/config"
	"wb-l0/internal/service"
	"wb-l0/internal/service/middleware"
	"wb-l0/internal/service/usecase"
)

type FiberServer struct {
	server  *fiber.App
	service service.IServices
	Log     zerolog.Logger
}

func NewServer(conf *config.Config) (service.RestServer, error) {
	s := FiberServer{}
	var err error

	s.service, err = usecase.NewServices(conf)
	if err != nil {
		return nil, err
	}

	s.server = fiber.New()
	s.server.Get("getModel/:id", s.GetModel)
	s.server.Post("addModel", s.AddModel)
	return &s, nil
}

func (s *FiberServer) StartServer(port string) error {
	err := s.server.Listen(port)
	if err != nil {
		s.Log.Error().Timestamp().Err(err)
		return err
	}
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
			return err
		}
	}
	return nil
}

func (s *FiberServer) AddModel(ctx *fiber.Ctx) error {
	object := ctx.Body()
	id, err := middleware.CheckModel(object)
	if err != nil {
		s.Log.Error().Timestamp().Err(err)
		return err
	}
	err = s.service.SaveModel(id, object)
	if err != nil {
		s.Log.Error().Timestamp().Err(err)
		return err
	}
	return nil
}
