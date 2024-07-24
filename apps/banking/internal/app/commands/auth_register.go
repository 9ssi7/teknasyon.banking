package commands

import (
	"context"
	"errors"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/events"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/validation"
	"go.opentelemetry.io/otel/trace"
)

type AuthRegister struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type AuthRegisterHandler cqrs.HandlerFunc[AuthRegister, *cqrs.Empty]

func NewAuthRegisterHandler(tracer trace.Tracer, v validation.Service, userRepo abstracts.UserRepo) AuthRegisterHandler {
	return func(ctx context.Context, cmd AuthRegister) (*cqrs.Empty, error) {
		ctx, span := tracer.Start(ctx, "AuthRegisterHandler")
		defer span.End()
		err := v.ValidateStruct(ctx, cmd)
		if err != nil {
			return nil, err
		}
		exists, err := userRepo.IsExistsByEmail(ctx, cmd.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, rescode.EmailAlreadyExists(errors.New("email already exists"))
		}
		u := entities.NewUser(cmd.Name, cmd.Email)
		err = userRepo.Save(ctx, u)
		if err != nil {
			return nil, err
		}
		events.OnAuthRegistered(events.AuthRegistered{
			Name:             cmd.Name,
			Email:            cmd.Email,
			VerificationCode: *u.TempToken,
		})
		return &cqrs.Empty{}, nil
	}
}
