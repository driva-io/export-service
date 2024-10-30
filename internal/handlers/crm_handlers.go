package handlers

import (
	"export-service/internal/services/crm_exporter"

	"github.com/gofiber/fiber/v2"
)

func GetPipelinesHandler(c *fiber.Ctx, crmService crm_exporter.Crm, client any) error {

	pipelines, err := crmService.GetPipelines(client)

	status := fiber.StatusOK
	returnBody := pipelines
	if err != nil {
		return err
	}

	return c.Status(status).JSON(returnBody)
}

func GetFieldsHandler(c *fiber.Ctx, crmService crm_exporter.Crm, client any) error {
	fields, err := crmService.GetFields(client)

	status := fiber.StatusOK
	returnBody := fields
	if err != nil {
		return err
	}

	return c.Status(status).JSON(returnBody)
}

func GetOwnersHandler(c *fiber.Ctx, crmService crm_exporter.Crm, client any) error {
	fields, err := crmService.GetOwners(client)

	status := fiber.StatusOK
	returnBody := fields
	if err != nil {
		return err
	}

	return c.Status(status).JSON(returnBody)
}
