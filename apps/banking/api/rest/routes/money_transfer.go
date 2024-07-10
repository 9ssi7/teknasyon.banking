package routes

import (
	"github.com/9ssi7/banking/api/rest/middlewares"
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/commands"
	"github.com/gofiber/fiber/v2"
)

func MoneyTransfer(router fiber.Router, srv restsrv.Srv, app app.App) {
	group := router.Group("/money-transfer")
	group.Post("/", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(moneyTransfer(app)))
}

func moneyTransfer(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.MoneyTranfer
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		currentUser := middlewares.AccessMustParse(c)
		cmd.UserId = currentUser.Id
		cmd.UserEmail = currentUser.Email
		cmd.UserName = currentUser.Name
		res, err := app.Commands.MoneyTransfer(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
