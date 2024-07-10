package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountBalanceWithdraw struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id"  params:"account_id" validate:"required,uuid"`
	Amount    string    `json:"amount" validate:"required,amount"`
}

type AccountBalanceWithdrawHandler cqrs.HandlerFunc[AccountBalanceWithdraw, *cqrs.Empty]

func NewAccountBalanceWithdrawHandler(v validation.Service, accountRepo abstracts.AccountRepo, transactionRepo abstracts.TransactionRepo) AccountBalanceWithdrawHandler {
	return func(ctx context.Context, cmd AccountBalanceWithdraw) (*cqrs.Empty, error) {
		if err := v.ValidateStruct(ctx, cmd); err != nil {
			return nil, err
		}
		account, err := accountRepo.FindByUserIdAndId(ctx, cmd.UserId, cmd.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.IsAvailable() {
			return nil, rescode.AccountNotAvailable
		}
		amount, err := decimal.NewFromString(cmd.Amount)
		if err != nil {
			return nil, rescode.Failed
		}
		if !account.CanCredit(amount) {
			return nil, rescode.AccountBalanceInsufficient
		}
		account.SubBalance(amount)
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		t := entities.NewTransaction(account.Id, account.Id, amount, "Withdraw balance", valobj.TransactionKindWithdrawal)
		if err := transactionRepo.Save(ctx, t); err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
