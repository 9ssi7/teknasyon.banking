package commands

import (
	"context"

	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/events"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MoneyTranfer struct {
	UserId      uuid.UUID `validate:"-"`
	UserEmail   string    `validate:"-"`
	UserName    string    `json:"user_name" validate:"-"`
	AccountId   uuid.UUID `json:"account_id" validate:"required,uuid"`
	Amount      string    `json:"amount" validate:"required,amount"`
	ToIban      string    `json:"to_iban" validate:"required,iban"`
	ToOwner     string    `json:"to_owner" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"required,min=3,max=255"`
}

type MoneyTransferHandler cqrs.HandlerFunc[MoneyTranfer, *cqrs.Empty]

func NewMoneyTransferHandler(v validation.Service, userRepo abstracts.UserRepo, accountRepo abstracts.AccountRepo, transactionRepo abstracts.TransactionRepo) MoneyTransferHandler {
	return func(ctx context.Context, cmd MoneyTranfer) (*cqrs.Empty, error) {
		if err := v.ValidateStruct(ctx, cmd); err != nil {
			return nil, err
		}
		toAccount, err := accountRepo.FindByIbanAndOwner(ctx, cmd.ToIban, cmd.ToOwner)
		if err != nil {
			return nil, rescode.AccountNotFound
		}
		account, err := accountRepo.FindByUserIdAndId(ctx, cmd.UserId, cmd.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.IsAvailable() {
			return nil, rescode.AccountNotAvailable
		}
		if !toAccount.IsAvailable() {
			return nil, rescode.ToAccountNotAvailable
		}
		if account.Id == toAccount.Id {
			return nil, rescode.AccountTransferToSameAccount
		}
		if account.Currency != toAccount.Currency {
			return nil, rescode.AccountCurrencyMismatch
		}
		amountToTransfer, err := decimal.NewFromString(cmd.Amount)
		if err != nil {
			return nil, err
		}
		amountToPay := amountToTransfer
		if account.UserId != toAccount.UserId && config.ReadValue().ProcessFee != 0 {
			amountToPay = amountToTransfer.Add(decimal.NewFromInt(int64(config.ReadValue().ProcessFee)))
		}

		if !account.CanCredit(amountToPay) {
			return nil, rescode.AccountBalanceInsufficient
		}

		transaction := entities.NewTransaction(account.Id, toAccount.Id, amountToTransfer, cmd.Description, valobj.TransactionKindTransfer)
		if err := transactionRepo.Save(ctx, transaction); err != nil {
			return nil, err
		}
		if !amountToPay.Equal(amountToTransfer) {
			transactionFee := entities.NewTransaction(account.Id, account.Id, decimal.NewFromInt(int64(config.ReadValue().ProcessFee)), "Process Fee", valobj.TransactionKindFee)
			if err := transactionRepo.Save(ctx, transactionFee); err != nil {
				return nil, err
			}
		}

		account.Debit(amountToPay)
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		toAccount.Credit(amountToTransfer)
		if err := accountRepo.Save(ctx, toAccount); err != nil {
			return nil, err
		}

		if toAccount.UserId != account.UserId {
			toUser, err := userRepo.FindById(ctx, toAccount.UserId)
			if err != nil {
				return nil, rescode.Failed
			}
			events.OnTransferIncoming(events.TranfserIncoming{
				Email:       toUser.Email,
				Name:        toUser.Name,
				Amount:      amountToTransfer.String(),
				Currency:    toAccount.Currency.String(),
				Account:     toAccount.Name,
				Description: cmd.Description,
			})
			events.OnTransferOutgoing(events.TranfserOutgoing{
				Amount:      amountToPay.String(),
				Email:       cmd.UserEmail,
				Name:        cmd.UserName,
				Currency:    account.Currency.String(),
				Account:     account.Name,
				Description: cmd.Description,
			})
		}
		return &cqrs.Empty{}, nil
	}
}
