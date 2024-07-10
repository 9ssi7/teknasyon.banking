package queries

import (
	"context"

	"github.com/9ssi7/banking/internal/app/dtos"
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/cqrs"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/google/uuid"
)

type TransactionList struct {
	AccountId uuid.UUID `json:"account_id" params:"account_id" validate:"required,uuid"`
	Pagi      list.PagiRequest
	Filters   valobj.TransactionFilters
}

type TransactionListHandler cqrs.HandlerFunc[TransactionList, *list.PagiResponse[*dtos.TransactionListDto]]

func NewTransactionListHandler(v validation.Service, transactionRepo abstracts.TransactionRepo, accountRepo abstracts.AccountRepo) TransactionListHandler {
	return func(ctx context.Context, query TransactionList) (*list.PagiResponse[*dtos.TransactionListDto], error) {
		if err := v.ValidateStruct(ctx, query); err != nil {
			return nil, err
		}
		_, err := accountRepo.FindById(ctx, query.AccountId)
		if err != nil {
			return nil, err
		}
		filters, err := transactionRepo.Filter(ctx, query.AccountId, &query.Pagi, &query.Filters)
		if err != nil {
			return nil, err
		}
		result := &list.PagiResponse[*dtos.TransactionListDto]{
			List:          make([]*dtos.TransactionListDto, 0, len(filters.List)),
			Total:         filters.Total,
			FilteredTotal: filters.FilteredTotal,
			Limit:         filters.Limit,
			TotalPage:     filters.TotalPage,
			Page:          filters.Page,
		}
		for _, e := range filters.List {
			dto := dtos.MapTransactionListItem(e, query.AccountId)
			if dto.AccountId != nil {
				a, err := accountRepo.FindById(ctx, *dto.AccountId)
				if err != nil {
					return nil, err
				}
				dto.AccountName = &a.Name
			}
		}
		return result, nil
	}
}
