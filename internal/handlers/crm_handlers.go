package handlers

import (
	"context"
	"errors"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/presentation_spec_repo"
	"export-service/internal/services/crm_exporter"
	"export-service/internal/services/data_presenter"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func InstallHandler(c *fiber.Ctx, crmService crm_exporter.Crm) error {
	//hubspot doesnt require install data like a token (its oauth), other CRMs may require
	installData := map[string]any{
		"workspace_id": c.Query("workspace_id"),
		"user_id":      c.Query("user_id"),
		"company":      c.Query("company"),
	}
	response, err := crmService.Install(installData)

	status := fiber.StatusOK
	returnBody := response
	if err != nil {
		return err
	}

	return c.Status(status).JSON(returnBody)
}

func OAuthCallBackHandler(c *fiber.Ctx, crmService crm_exporter.Crm) error {
	state := c.Query("state")
	unescapedState, err := url.QueryUnescape(state)
	if err != nil {
		return err
	}
	parts := strings.Split(unescapedState, "|")
	if len(parts) != 3 {
		return errors.New("invalid state parameters")
	}

	workspaceID, userID, company := parts[0], parts[1], parts[2]

	if workspaceID == "" || userID == "" || company == "" {
		return errors.New("invalid state parameters")
	}

	_, err = crmService.OAuthCallback(c, workspaceID, userID, company)

	status := fiber.StatusNoContent
	if err != nil {
		status = fiber.StatusInternalServerError
		return c.Status(status).JSON(err)
	}

	return c.SendStatus(status)
}

func GetPipelinesHandler(c *fiber.Ctx, crmService crm_exporter.Crm, client any) error {

	pipelines, err := crmService.GetPipelines(client)

	status := fiber.StatusOK
	returnBody := pipelines
	if err != nil {
		return err
	}

	return c.Status(status).JSON(returnBody)
}

func TestLeadHandler(c *fiber.Ctx, crmService crm_exporter.Crm, client any, p *presentation_spec_repo.PgPresentationSpecRepository) error {
	getSpecParms := ports.PresentationSpecQueryParams{
		UserEmail:   c.Query("user_email"),
		UserCompany: c.Query("company"),
		Service:     "crm_" + c.Params("crm"),
		DataSource:  c.Query("base"),
	}

	var configs map[string]any
	if err := c.BodyParser(&configs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
	}

	spec, err := p.Get(context.Background(), getSpecParms)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	mappedLead, err := data_presenter.PresentSingle(crm_exporter.DrivaTestLead, spec.Spec)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	result, err := crmService.SendLead(client, mappedLead, map[string]any{}, configs, map[string]map[string]any{})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"lead":  result,
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func ValidateHandler(c *fiber.Ctx, crmService crm_exporter.Crm, client any) error {

	isValid := crmService.Validate(c, client)

	status := fiber.StatusOK
	returnBody := map[string]bool{
		"valid": isValid,
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
