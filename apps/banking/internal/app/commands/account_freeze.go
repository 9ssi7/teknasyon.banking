package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/google/uuid"
)

type AccountFreeze struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id" validate:"-"`
}

type AccountFreezeHandler cqrs.HandlerFunc[AccountFreeze, *cqrs.Empty]

func NewAccountFreezeHandler(accountRepo abstracts.AccountRepo) AccountFreezeHandler {
	return func(ctx context.Context, query AccountFreeze) (*cqrs.Empty, error) {
		account, err := accountRepo.FindByUserIdAndId(ctx, query.UserId, query.AccountId)
		if err != nil {
			return nil, err
		}
		account.Freeze()
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
