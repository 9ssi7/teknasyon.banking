package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/currency"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
)

type AccountCreate struct {
	UserId   uuid.UUID `json:"user_id" validate:"-"`
	Name     string    `json:"name" validate:"required,min=3,max=255"`
	Owner    string    `json:"owner" validate:"required,min=3,max=255"`
	Currency string    `json:"currency" validate:"required,currency"`
}

type AccountCreateHandler cqrs.HandlerFunc[AccountCreate, *uuid.UUID]

func NewAccountCreateHandler(v validation.Service, accountRepo abstracts.AccountRepo) AccountCreateHandler {
	return func(ctx context.Context, cmd AccountCreate) (*uuid.UUID, error) {
		if err := v.ValidateStruct(ctx, cmd); err != nil {
			return nil, err
		}
		account := entities.NewAccount(cmd.UserId, cmd.Name, cmd.Owner, currency.Currency(cmd.Currency))
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &account.Id, nil
	}
}
