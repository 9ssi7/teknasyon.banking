package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
)

type AuthVerify struct {
	Token string `json:"token" validate:"required,uuid"`
}

type AuthVerifyHandler cqrs.HandlerFunc[AuthVerify, *cqrs.Empty]

func NewAuthVerifyHandler(userRepo abstracts.UserRepo) AuthVerifyHandler {
	return func(ctx context.Context, cmd AuthVerify) (*cqrs.Empty, error) {
		u, err := userRepo.FindByToken(ctx, cmd.Token)
		if err != nil {
			return nil, err
		}
		u.Verify()
		err = userRepo.Save(ctx, u)
		if err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
