package handlers

import (
	"errors"
	"export-service/internal/core/ports"
	"export-service/internal/repositories"
	"export-service/internal/repositories/presentation_spec_repo"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetPresentationSpecHandler(c *fiber.Ctx, p *presentation_spec_repo.PgPresentationSpecRepository) error {
	companyName := c.Query("user_company")
	userEmail := c.Query("user_email")
	service := c.Query("service")
	base := c.Query("base")

	pSpec, err := p.Get(c.Context(), ports.PresentationSpecQueryParams{
		UserEmail:   userEmail,
		UserCompany: companyName,
		Service:     service,
		DataSource:  base,
	})

	status := fiber.StatusOK
	var returnBody any
	returnBody = pSpec
	if err != nil {
		returnBody = err
		var notFoundErr repositories.PresentationSpecNotFoundError
		var notUniqueErr repositories.PresentationSpecNotUniqueError
		var invalidQueryParams ports.InvalidQueryParamsError

		switch {
		case errors.As(err, &notFoundErr):
			status = fiber.StatusNotFound
		case errors.As(err, &notUniqueErr):
			status = fiber.StatusConflict
		case errors.As(err, &invalidQueryParams):
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}
	}

	return c.Status(status).JSON(returnBody)
}

func AddPresentationSpecHandler(c *fiber.Ctx, p *presentation_spec_repo.PgPresentationSpecRepository) error {
	companyName := c.Query("user_company")
	userEmail := c.Query("user_email")
	service := c.Query("service")
	base := c.Query("base")

	var body ports.PresentationSpecAddBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())

	}

	validate := validator.New()
	if invalidStruct := validate.Struct(body); invalidStruct != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
	}

	addedSpec, err := p.Add(c.Context(), ports.PresentationSpecQueryParams{
		UserEmail:   userEmail,
		UserCompany: companyName,
		Service:     service,
		DataSource:  base,
	}, body)

	status := fiber.StatusOK
	var returnBody any
	returnBody = addedSpec
	if err != nil {
		returnBody = err
		var invalidQueryParams ports.InvalidQueryParamsError
		var invalidBody ports.InvalidBodyError

		switch {
		case errors.As(err, &invalidBody):
			status = fiber.StatusBadRequest
		case errors.As(err, &invalidQueryParams):
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}
	}

	return c.Status(status).JSON(returnBody)
}

func PatchPresentationSpecHandler(c *fiber.Ctx, p *presentation_spec_repo.PgPresentationSpecRepository) error {
	id := c.Params("id")

	var body ports.PresentationSpecAddBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
	}

	validate := validator.New()
	if invalidStruct := validate.Struct(body); invalidStruct != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
	}

	updatedSpec, err := p.Patch(c.Context(), id, body)
	status := fiber.StatusOK
	var returnBody any
	returnBody = updatedSpec
	if err != nil {
		returnBody = err
		var invalidQueryParams ports.InvalidQueryParamsError
		var invalidBody ports.InvalidBodyError
		switch {
		case errors.As(err, &invalidBody):
			status = fiber.StatusBadRequest
		case errors.As(err, &invalidQueryParams):
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}
	}

	return c.Status(status).JSON(returnBody)
}

func PatchPresentationSpecKeyHandler(c *fiber.Ctx, p *presentation_spec_repo.PgPresentationSpecRepository) error {
	id := c.Params("id")
	key := c.Params("key")

	var body ports.PresentationSpecPatchKey
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
	}

	validate := validator.New()
	if invalidStruct := validate.Struct(body); invalidStruct != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ports.NewInvalidBodyError())
	}

	updatedSpec, err := p.PatchKey(c.Context(), id, key, body)
	status := fiber.StatusOK
	var returnBody any
	returnBody = updatedSpec
	if err != nil {
		returnBody = err
		var invalidQueryParams ports.InvalidQueryParamsError
		var invalidBody ports.InvalidBodyError
		switch {
		case errors.As(err, &invalidBody):
			status = fiber.StatusBadRequest
		case errors.As(err, &invalidQueryParams):
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}
	}

	return c.Status(status).JSON(returnBody)
}

func DeletePresentationSpecHandler(c *fiber.Ctx, p *presentation_spec_repo.PgPresentationSpecRepository) error {
	id := c.Params("id")

	err := p.Delete(c.Context(), id)

	status := fiber.StatusNoContent
	if err != nil {
		var invalidQueryParams ports.InvalidQueryParamsError

		switch {
		case errors.As(err, &invalidQueryParams):
			status = fiber.StatusBadRequest
		default:
			status = fiber.StatusInternalServerError
		}

		return c.Status(status).JSON(err)
	}

	return c.SendStatus(status)
}
