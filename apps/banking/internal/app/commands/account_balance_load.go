package commands

import (
	"context"
	"errors"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/events"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/trace"
)

type AccountBalanceLoad struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	UserEmail string    `json:"user_email" validate:"-"`
	UserName  string    `json:"user_name" validate:"-"`
	AccountId uuid.UUID `json:"account_id"  params:"account_id" validate:"required,uuid"`
	Amount    string    `json:"amount" validate:"required,amount"`
}

type AccountBalanceLoadHandler cqrs.HandlerFunc[AccountBalanceLoad, *cqrs.Empty]

func NewAccountBalanceLoadHandler(tracer trace.Tracer, v validation.Service, accountRepo abstracts.AccountRepo, transactionRepo abstracts.TransactionRepo) AccountBalanceLoadHandler {
	return func(ctx context.Context, cmd AccountBalanceLoad) (*cqrs.Empty, error) {
		ctx, span := tracer.Start(ctx, "AccountBalanceLoadHandler")
		defer span.End()
		if err := v.ValidateStruct(ctx, cmd); err != nil {
			return nil, err
		}
		account, err := accountRepo.FindByUserIdAndId(ctx, cmd.UserId, cmd.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.IsAvailable() {
			return nil, rescode.AccountNotAvailable(errors.New("sender account not available"))
		}
		amount, err := decimal.NewFromString(cmd.Amount)
		if err != nil {
			return nil, rescode.Failed(err)
		}
		account.Credit(amount)
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		t := entities.NewTransaction(account.Id, account.Id, amount, "Load balance", valobj.TransactionKindDeposit)
		if err := transactionRepo.Save(ctx, t); err != nil {
			return nil, err
		}
		events.OnTransferIncoming(events.TranfserIncoming{
			Name:        cmd.UserName,
			Amount:      amount.String(),
			Currency:    account.Currency.String(),
			Email:       cmd.UserEmail,
			Account:     account.Name,
			Description: "Load balance",
		})
		return &cqrs.Empty{}, nil
	}
}
