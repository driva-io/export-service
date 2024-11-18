package handlers

import (
	"export-service/internal/services/crm_exporter"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New(session.Config{
	CookieSecure:   true,           
	CookieSameSite: "None",         
	CookieHTTPOnly: true,          
	Expiration:     24 * time.Hour, 
})


func InstallHandler(c *fiber.Ctx, crmService crm_exporter.Crm) error {
	//hubspot doesnt require install data like a token (its oauth), other CRMs may require
	response, err := crmService.Install(nil)

	status := fiber.StatusOK
	returnBody := response
	if err != nil {
		return err
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}
	defer sess.Save()

	sess.Set("workspace_id", c.Query("workspace_id"))
	sess.Set("user_id", c.Query("user_id"))
	sess.Set("company", c.Query("company"))

	return c.Status(status).JSON(returnBody)
}

func OAuthCallBackHandler(c *fiber.Ctx, crmService crm_exporter.Crm) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("session error")
	}
	defer sess.Save()

	workspaceID := sess.Get("workspace_id")
	userID := sess.Get("user_id")
	company := sess.Get("company")

	defer sess.Delete("workspace_id")
	defer sess.Delete("user_id")
	defer sess.Delete("company")

	if workspaceID == nil || userID == nil || company == nil {
		return c.Status(fiber.StatusBadRequest).SendString("necessary session data not found")
	}

	_ , err = crmService.OAuthCallback(c, workspaceID, userID, company)

	status := fiber.StatusNoContent
	if err != nil {
		status = fiber.StatusInternalServerError
		return c.Status(status).JSON(err);
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
