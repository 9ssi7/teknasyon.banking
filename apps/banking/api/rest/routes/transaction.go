package routes

import (
	"github.com/9ssi7/banking/api/rest/middlewares"
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/gofiber/fiber/v2"
)

func Transaction(router fiber.Router, srv restsrv.Srv, app app.App) {
	group := router.Group("/transactions")
	group.Get("/", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(transactionList(app)))
}

func transactionList(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagi list.PagiRequest
		if err := c.QueryParser(&pagi); err != nil {
			return err
		}
		var filters valobj.TransactionFilters
		if err := c.QueryParser(&filters); err != nil {
			return err
		}
		var query queries.TransactionList
		if err := c.QueryParser(&query); err != nil {
			return err
		}
		pagi.Default()
		query.Pagi = pagi
		query.Filters = filters
		query.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Queries.TransactionList(c.UserContext(), query)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
