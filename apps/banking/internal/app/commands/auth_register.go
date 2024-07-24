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
)

type AuthRegister struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type AuthRegisterHandler cqrs.HandlerFunc[AuthRegister, *cqrs.Empty]

func NewAuthRegisterHandler(v validation.Service, userRepo abstracts.UserRepo) AuthRegisterHandler {
	return func(ctx context.Context, cmd AuthRegister) (*cqrs.Empty, error) {
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
