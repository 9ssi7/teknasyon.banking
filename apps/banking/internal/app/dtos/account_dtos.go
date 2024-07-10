package dtos

import (
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/pkg/currency"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/google/uuid"
)

type AccountListDto struct {
	Id       uuid.UUID         `json:"id"`
	Name     string            `json:"name"`
	Owner    string            `json:"owner"`
	Iban     string            `json:"iban"`
	Currency currency.Currency `json:"currency"`
	Balance  string            `json:"balance"`
	Status   string            `json:"status"`
}

func MapAccountList(res *list.PagiResponse[*entities.Account]) *list.PagiResponse[*AccountListDto] {
	var result []*AccountListDto
	for _, item := range res.List {
		result = append(result, &AccountListDto{
			Id:       item.Id,
			Name:     item.Name,
			Owner:    item.Owner,
			Iban:     item.Iban,
			Currency: item.Currency,
			Balance:  item.Balance.String(),
			Status:   item.Status.String(),
		})
	}
	return &list.PagiResponse[*AccountListDto]{
		List:          result,
		Total:         res.Total,
		Limit:         res.Limit,
		TotalPage:     res.TotalPage,
		FilteredTotal: res.FilteredTotal,
		Page:          res.Page,
	}
}
