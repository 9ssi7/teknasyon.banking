package dtos

import (
	"time"

	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/google/uuid"
)

type TransactionListDto struct {
	Id          uuid.UUID  `json:"id"`
	AccountId   *uuid.UUID `json:"account_id,omitempty"`
	AccountName *string    `json:"account_name,omitempty"`
	Amount      string     `json:"amount"`
	Description string     `json:"description"`
	Kind        string     `json:"kind"`
	Direction   string     `json:"direction"`
	CreatedAt   string     `json:"created_at"`
}

func MapTransactionListItem(e *entities.Transaction, userId uuid.UUID) *TransactionListDto {
	d := &TransactionListDto{
		Id:          e.Id,
		Amount:      e.Amount.String(),
		Description: e.Description,
		Kind:        e.Kind.String(),
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
	}
	if e.IsItself() {
		d.Direction = "self"
	} else if e.IsUserSender(userId) {
		d.Direction = "outgoing"
		d.AccountId = &e.ReceiverId
	} else {
		d.Direction = "incoming"
		d.AccountId = &e.SenderId
	}
	return d
}
