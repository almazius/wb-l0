package rest_api

import (
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	// usecase model
	server *fiber.App
}

func NewServer() FiberServer {
	s := FiberServer{}
	s.server = fiber.New()
	s.server.Get("GetModel/ip", GetModel)
	s.server.Post("", AddModel)
	return s
}

func GetModel(ctx *fiber.Ctx) error {
	return nil
}

func AddModel(ctx *fiber.Ctx) error {
	return nil
}
