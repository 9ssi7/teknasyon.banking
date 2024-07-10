package routes

import (
	"github.com/9ssi7/banking/api/rest/middlewares"
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/commands"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/gofiber/fiber/v2"
)

func Account(router fiber.Router, srv restsrv.Srv, app app.App) {
	group := router.Group("/account")
	group.Post("/", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountCreate(app)))
	group.Patch("/:account_id/activate", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountActivate(app)))
	group.Patch("/:account_id/freeze", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountFreeze(app)))
	group.Patch("/:account_id/lock", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountLock(app)))
	group.Patch("/:account_id/suspend", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountSuspend(app)))
	group.Patch("/:account_id/balance/load", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountBalanceLoad(app)))
	group.Patch("/:account_id/balance/withdraw", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountBalanceWithdraw(app)))
	group.Get("/", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(accountList(app)))
}

func accountCreate(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountCreate
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountCreate(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountActivate(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountActivate
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountActivate(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountFreeze(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountFreeze
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountFreeze(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountLock(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountLock
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountLock(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountSuspend(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountSuspend
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountSuspend(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountBalanceLoad(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountBalanceLoad
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountBalanceLoad(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountBalanceWithdraw(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AccountBalanceWithdraw
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		cmd.UserId = middlewares.AccessMustParse(c).Id
		res, err := app.Commands.AccountBalanceWithdraw(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func accountList(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagi list.PagiRequest
		if err := c.QueryParser(&pagi); err != nil {
			return err
		}
		pagi.Default()
		query := queries.AccountList{
			UserId: middlewares.AccessMustParse(c).Id,
			Pagi:   pagi,
		}
		res, err := app.Queries.AccountList(c.UserContext(), query)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
