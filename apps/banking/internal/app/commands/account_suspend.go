package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/google/uuid"
)

type AccountSuspend struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id" validate:"-"`
}

type AccountSuspendHandler cqrs.HandlerFunc[AccountSuspend, *cqrs.Empty]

func NewAccountSuspendHandler(accountRepo abstracts.AccountRepo) AccountSuspendHandler {
	return func(ctx context.Context, query AccountSuspend) (*cqrs.Empty, error) {
		account, err := accountRepo.FindByUserIdAndId(ctx, query.UserId, query.AccountId)
		if err != nil {
			return nil, err
		}
		account.Suspend()
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
