package routes

import (
	"export-service/internal/gateways"
	"export-service/internal/handlers"
	"export-service/internal/server"
)

func RegisterServerRoutes(s *server.FiberServer, a gateways.AuthServiceGateway) {
	s.Get("/v1/health", handlers.HealthHandler)
}
