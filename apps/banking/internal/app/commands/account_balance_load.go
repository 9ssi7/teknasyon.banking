package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountBalanceLoad struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id" validate:"required,uuid"`
	Amount    string    `json:"amount" validate:"required,amount"`
}

type AccountBalanceLoadHandler cqrs.HandlerFunc[AccountBalanceLoad, *cqrs.Empty]

func NewAccountBalanceLoadHandler(accountRepo abstracts.AccountRepo) AccountBalanceLoadHandler {
	return func(ctx context.Context, cmd AccountBalanceLoad) (*cqrs.Empty, error) {
		account, err := accountRepo.FindByUserIdAndId(ctx, cmd.UserId, cmd.AccountId)
		if err != nil {
			return nil, err
		}
		amount, err := decimal.NewFromString(cmd.Amount)
		if err != nil {
			return nil, rescode.Failed
		}
		account.AddBalance(amount)
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
