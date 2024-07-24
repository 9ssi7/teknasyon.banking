package queries

import (
	"context"

	"github.com/9ssi7/banking/internal/app/dtos"
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AccountList struct {
	UserId uuid.UUID
	Pagi   list.PagiRequest
}

type AccountListHandler cqrs.HandlerFunc[AccountList, *list.PagiResponse[*dtos.AccountListDto]]

func NewAccountListHandler(tracer trace.Tracer, accountRepo abstracts.AccountRepo) AccountListHandler {
	return func(ctx context.Context, query AccountList) (*list.PagiResponse[*dtos.AccountListDto], error) {
		ctx, span := tracer.Start(ctx, "AccountListHandler")
		defer span.End()
		res, err := accountRepo.ListByUserId(ctx, query.UserId, &query.Pagi)
		if err != nil {
			return nil, err
		}
		return dtos.MapAccountList(res), nil
	}
}
