package commands

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
)

type AccountActivate struct {
	UserId    uuid.UUID `json:"user_id" validate:"-"`
	AccountId uuid.UUID `json:"account_id"  params:"account_id" validate:"required,uuid"`
}

type AccountActivateHandler cqrs.HandlerFunc[AccountActivate, *cqrs.Empty]

func NewAccountActivateHandler(v validation.Service, accountRepo abstracts.AccountRepo) AccountActivateHandler {
	return func(ctx context.Context, cmd AccountActivate) (*cqrs.Empty, error) {
		if err := v.ValidateStruct(ctx, cmd); err != nil {
			return nil, err
		}
		account, err := accountRepo.FindByUserIdAndId(ctx, cmd.UserId, cmd.AccountId)
		if err != nil {
			return nil, err
		}
		account.Activate()
		if err := accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
		return &cqrs.Empty{}, nil
	}
}
