package routes

import (
	"export-service/api/middlewares"
	"export-service/internal/adapters"
	"export-service/internal/core/ports"
	"export-service/internal/server"
	"export-service/internal/usecases"
	"export-service/internal/writers"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.uber.org/zap"
)

func RegisterSheetRoutes(s *server.FiberServer, uploader ports.Uploader, specRepo ports.PresentationSpecRepository, mailer ports.Mailer, logger *zap.Logger) {
	sheetRoutes := s.App.Group("/sheet/v1")
	sheetRoutes.Use(middlewares.TokenMiddleware())

	sheetRoutes.Post("/export", func(c *fiber.Ctx) error {
		ctx := getContext(c)

		logger := logger.With(apmzap.TraceContext(ctx)...)

		var req usecases.ExportRequest
		err := c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to unmarshal message", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
		}

		sheetUc := usecases.NewSheetExportUseCase(&writers.ExcelWriter{}, &adapters.HTTPDownloader{}, uploader, specRepo, mailer, logger)
		downloadUrl, err := sheetUc.Execute(req)
		if err != nil {
			logger.Error("Failed to execute use case", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": "error while executing use case"})
		}

		return c.Status(fiber.StatusCreated).JSON(map[string]any{
			"download_url": downloadUrl,
		})
	})
}
