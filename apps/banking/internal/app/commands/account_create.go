package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/currency"
	"github.com/google/uuid"
)

type AccountCreate struct {
	UserId   uuid.UUID `json:"user_id" validate:"-"`
	Name     string    `json:"name" validate:"required,min=3,max=255"`
	Owner    string    `json:"owner" validate:"required,min=3,max=255"`
	Currency string    `json:"currency" validate:"required,currency"`
}

type AccountCreateHandler cqrs.HandlerFunc[AccountCreate, *uuid.UUID]

func NewAccountCreateHandler(accountRepo abstracts.AccountRepo) AccountCreateHandler {
	return func(ctx context.Context, query AccountCreate) (*uuid.UUID, error) {
		account := entities.NewAccount(query.UserId, query.Name, query.Owner, currency.Currency(query.Currency))
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &account.Id, nil
	}
}
