package queries

import (
	"context"
	"errors"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/state"
)

type AuthCheck struct {
	VerifyToken string `json:"-"`
}

type AuthCheckHandler cqrs.HandlerFunc[AuthCheck, *cqrs.Empty]

func NewAuthCheckHandler(verifyRepo abstracts.VerifyRepo) AuthCheckHandler {
	return func(ctx context.Context, query AuthCheck) (*cqrs.Empty, error) {
		exists, err := verifyRepo.IsExists(ctx, query.VerifyToken, state.GetDeviceId(ctx))
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, rescode.NotFound(errors.New("verify token not exists"))
		}
		return &cqrs.Empty{}, nil
	}
}
