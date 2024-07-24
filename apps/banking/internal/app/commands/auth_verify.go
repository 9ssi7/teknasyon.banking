package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"go.opentelemetry.io/otel/trace"
)

type AuthVerify struct {
	Token string `json:"token" validate:"required,uuid"`
}

type AuthVerifyHandler cqrs.HandlerFunc[AuthVerify, *cqrs.Empty]

func NewAuthVerifyHandler(tracer trace.Tracer, userRepo abstracts.UserRepo) AuthVerifyHandler {
	return func(ctx context.Context, cmd AuthVerify) (*cqrs.Empty, error) {
		ctx, span := tracer.Start(ctx, "AuthVerifyHandler")
		defer span.End()
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
