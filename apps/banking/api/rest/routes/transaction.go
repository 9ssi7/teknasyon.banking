package routes

import (
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/gofiber/fiber/v2"
)

func Transaction(router fiber.Router, srv restsrv.Srv, app app.App) {
	group := router.Group("/transactions")
	group.Get("/", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(transactionList(app)))
}

func transactionList(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var query queries.TransactionList
		if err := c.QueryParser(&query); err != nil {
			return err
		}
		query.Pagi.Default()
		res, err := app.Queries.TransactionList(c.UserContext(), query)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
