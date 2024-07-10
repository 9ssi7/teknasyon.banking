package entities

import (
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Base
	SenderId    uuid.UUID              `json:"sender_id" gorm:"type:uuid;not null"`
	ReceiverId  uuid.UUID              `json:"receiver_id" gorm:"type:uuid;not null"`
	Amount      decimal.Decimal        `json:"amount" gorm:"type:decimal;not null"`
	Description string                 `json:"description" gorm:"type:varchar(255);not null"`
	Kind        valobj.TransactionKind `json:"kind" gorm:"type:varchar(255);not null"`
}

func (t *Transaction) IsItself() bool {
	return t.SenderId == t.ReceiverId
}

func (t *Transaction) IsUserSender(userId uuid.UUID) bool {
	return t.SenderId == userId
}

func (t *Transaction) IsUserReceiver(userId uuid.UUID) bool {
	return t.ReceiverId == userId
}

func NewTransaction(senderId, receiverId uuid.UUID, amount decimal.Decimal, description string, kind valobj.TransactionKind) *Transaction {
	return &Transaction{
		SenderId:    senderId,
		ReceiverId:  receiverId,
		Amount:      amount,
		Description: description,
		Kind:        kind,
	}
}
