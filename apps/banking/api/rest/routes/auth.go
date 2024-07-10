package routes

import (
	"github.com/9ssi7/banking/api/rest/middlewares"
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/commands"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/gofiber/fiber/v2"
)

func Auth(router fiber.Router, srv restsrv.Srv, app app.App) {
	group := router.Group("/auth")
	group.Post("/start", srv.VerifyTokenExcluded(), srv.Timeout(authStart(srv, app)))
	group.Post("/login", srv.AccessInit(), srv.AccessExcluded(), srv.VerifyTokenRequired(), srv.Timeout(authLogin(srv, app)))
	group.Post("/register", srv.AccessInit(), srv.AccessExcluded(), srv.Turnstile(), srv.Timeout(authRegister(app)))
	group.Post("/logout", srv.AccessInit(true), srv.AccessRequired(true), srv.Timeout(authLogout(app)))
	group.Put("/refresh", srv.RefreshInit(), srv.RefreshRequired(), srv.Timeout(authRefresh(app)))
	group.Get("/verify/:token", srv.AccessInit(), srv.AccessExcluded(), srv.Turnstile(), srv.Timeout(authVerify(app)))
	group.Get("/check", srv.AccessInit(), srv.AccessExcluded(), srv.Timeout(authCheck(app)))
	group.Get("/", srv.AccessInit(), srv.AccessRequired(), srv.Timeout(authCurrent(app)))
}

func authLogin(srv restsrv.Srv, app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AuthLogin
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		cmd.Device = srv.MakeDevice(c)
		cmd.VerifyToken = middlewares.VerifyTokenParse(c)
		res, err := app.Commands.AuthLogin(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		middlewares.VerifyTokenRemove(c)
		middlewares.AccessTokenSetCookie(c, res.AccessToken)
		middlewares.RefreshTokenSetCookie(c, res.RefreshToken)
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authRegister(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AuthRegister
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		res, err := app.Commands.AuthRegister(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authVerify(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AuthVerify
		if err := c.ParamsParser(&cmd); err != nil {
			return err
		}
		res, err := app.Commands.AuthVerify(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authLogout(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cmd := commands.AuthLogout{
			UserId: middlewares.AccessMustParse(c).Id,
		}
		res, err := app.Commands.AuthLogout(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		middlewares.AccessTokenRemoveCookie(c)
		middlewares.RefreshTokenRemoveCookie(c)
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authRefresh(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cmd := commands.AuthRefresh{
			UserId:       middlewares.RefreshMustParse(c).Id,
			AccessToken:  middlewares.AccessGetToken(c),
			RefreshToken: middlewares.RefreshParseToken(c),
			IpAddress:    middlewares.IpMustParse(c),
		}
		res, err := app.Commands.AuthRefresh(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		middlewares.AccessTokenSetCookie(c, res.AccessToken)
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authStart(srv restsrv.Srv, app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cmd commands.AuthStart
		if err := c.BodyParser(&cmd); err != nil {
			return err
		}
		cmd.Device = srv.MakeDevice(c)
		res, err := app.Commands.AuthStart(c.UserContext(), cmd)
		if err != nil {
			return err
		}
		middlewares.VerifyTokenSet(c, res.VerifyToken)
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authCheck(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := queries.AuthCheck{
			VerifyToken: middlewares.VerifyTokenParse(c),
		}
		res, err := app.Queries.AuthCheck(c.UserContext(), query)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func authCurrent(app app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(middlewares.AccessMustParse(c).User)
	}
}
