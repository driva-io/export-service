package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/v2"
)

func getContext(c *fiber.Ctx) context.Context {
	tx := apm.TransactionFromContext(c.Context())

	return apm.ContextWithTransaction(context.Background(), tx)
}
