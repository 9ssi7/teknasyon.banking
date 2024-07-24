package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AccountLock struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id"  params:"account_id" validate:"required,uuid"`
}

type AccountLockHandler cqrs.HandlerFunc[AccountLock, *cqrs.Empty]

func NewAccountLockHandler(tracer trace.Tracer, v validation.Service, accountRepo abstracts.AccountRepo) AccountLockHandler {
	return func(ctx context.Context, cmd AccountLock) (*cqrs.Empty, error) {
		ctx, span := tracer.Start(ctx, "AccountLockHandler")
		defer span.End()
		if err := v.ValidateStruct(ctx, cmd); err != nil {
			return nil, err
		}
		account, err := accountRepo.FindByUserIdAndId(ctx, cmd.UserId, cmd.AccountId)
		if err != nil {
			return nil, err
		}
		account.Lock()
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
