package commands

import (
	"context"
	"errors"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/aggregates"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/events"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/state"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AuthStart struct {
	Phone  string         `json:"phone" validate:"required_without=Email,omitempty,phone"`
	Email  string         `json:"email" validate:"required_without=Phone,omitempty,email"`
	Device *valobj.Device `json:"-"`
}

type AuthStartRes struct {
	VerifyToken string `json:"-"`
}

type AuthStartHandler cqrs.HandlerFunc[AuthStart, *AuthStartRes]

func NewAuthStartHandler(tracer trace.Tracer, v validation.Service, verifyRepo abstracts.VerifyRepo, userRepo abstracts.UserRepo) AuthStartHandler {
	return func(ctx context.Context, cmd AuthStart) (*AuthStartRes, error) {
		ctx, span := tracer.Start(ctx, "AuthStartHandler")
		defer span.End()
		err := v.ValidateStruct(ctx, cmd)
		if err != nil {
			return nil, err
		}
		var user *entities.User
		if cmd.Phone != "" {
			user, err = userRepo.FindByPhone(ctx, cmd.Phone)
			if err != nil {
				return nil, err
			}
		} else {
			user, err = userRepo.FindByEmail(ctx, cmd.Email)
			if err != nil {
				return nil, err
			}
		}
		if user == nil {
			return nil, rescode.NotFound(errors.New("user not found"))
		}
		if !user.IsActive {
			return nil, rescode.UserDisabled(errors.New("user disabled"))
		}
		if user.TempToken != nil && *user.TempToken != "" {
			return nil, rescode.UserVerifyRequired(errors.New("user verify required"))
		}
		verifyToken := uuid.New().String()
		verify := aggregates.NewVerify(user.Id, state.GetDeviceId(ctx), state.GetLocale(ctx))
		err = verifyRepo.Save(ctx, verifyToken, verify)
		if err != nil {
			return nil, err
		}
		events.OnAuthStarted(events.AuthStarted{
			Email:  user.Email,
			Code:   verify.Code,
			Device: *cmd.Device,
		})
		return &AuthStartRes{
			VerifyToken: verifyToken,
		}, nil
	}
}
