package middlewares

import (
	"errors"
	"time"

	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/gofiber/fiber/v2"
)

func VerifyRequired(ctx *fiber.Ctx) error {
	token := ctx.Cookies("verify_token")
	if token == "" {
		return rescode.RequiredVerifyToken(errors.New("verify required"))
	}
	return ctx.Next()
}

func VerifyExcluded(ctx *fiber.Ctx) error {
	token := ctx.Cookies("verify_token")
	if token != "" {
		return rescode.ExcludedVerifyToken(errors.New("verify excluded"))
	}
	return ctx.Next()
}

func VerifyTokenParse(ctx *fiber.Ctx) string {
	return ctx.Cookies("verify_token")
}

func VerifyTokenSet(ctx *fiber.Ctx, token string) {
	ctx.Cookie(config.ApplyCookie(&fiber.Cookie{
		Name:    "verify_token",
		Value:   token,
		Domain:  config.ReadValue().HttpHeaders.Domain,
		Expires: time.Now().Add(time.Minute * 5),
	}))
}

func VerifyTokenRemove(ctx *fiber.Ctx) {
	ctx.Cookie(config.ApplyCookie(&fiber.Cookie{
		Name:    "verify_token",
		Value:   "",
		Domain:  config.ReadValue().HttpHeaders.Domain,
		Expires: time.Now().Add(time.Hour * -1),
	}))
}
