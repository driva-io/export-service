package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.elastic.co/apm/module/apmfiber/v2"
)

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	app := fiber.New(fiber.Config{
		ServerHeader: "export-service",
		AppName:      "export-service",
	})
	app.Use(apmfiber.Middleware())

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
		return err
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://test.sales.driva.io,https://sales.driva.io,http://localhost:3000",
		AllowCredentials: true,
	}))

	server := &FiberServer{
		App: app,
	}

	return server
}
