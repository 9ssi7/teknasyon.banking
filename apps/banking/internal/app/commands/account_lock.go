package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/google/uuid"
)

type AccountLock struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id" validate:"-"`
}

type AccountLockHandler cqrs.HandlerFunc[AccountLock, *cqrs.Empty]

func NewAccountLockHandler(accountRepo abstracts.AccountRepo) AccountLockHandler {
	return func(ctx context.Context, query AccountLock) (*cqrs.Empty, error) {
		account, err := accountRepo.FindByUserIdAndId(ctx, query.UserId, query.AccountId)
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
