package entities

import (
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/currency"
	"github.com/9ssi7/banking/pkg/iban"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	Base
	UserId   uuid.UUID            `json:"user_id" gorm:"type:uuid;not null"`
	Name     string               `json:"name" gorm:"type:varchar(255);not null"`
	Owner    string               `json:"owner" gorm:"type:varchar(255);not null"`
	Iban     string               `json:"iban" gorm:"type:varchar(255);not null"`
	Currency currency.Currency    `json:"currency" gorm:"type:varchar(3);not null"`
	Balance  decimal.Decimal      `json:"balance" gorm:"type:decimal;not null"`
	Status   valobj.AccountStatus `json:"status" gorm:"type:varchar(255);not null"`
}

func (a *Account) Credit(amount decimal.Decimal) {
	a.Balance = a.Balance.Add(amount)
}

func (a *Account) Debit(amount decimal.Decimal) {
	a.Balance = a.Balance.Sub(amount)
}

func (a *Account) Lock() {
	a.Status = valobj.AccountStatusLocked
}

func (a *Account) Activate() {
	a.Status = valobj.AccountStatusActive
}

func (a *Account) Freeze() {
	a.Status = valobj.AccountStatusFrozen
}

func (a *Account) Suspend() {
	a.Status = valobj.AccountStatusSuspended
}

func (a *Account) IsAvailable() bool {
	return a.Status == valobj.AccountStatusActive
}

func (a *Account) CanCredit(amount decimal.Decimal) bool {
	return a.IsAvailable() && amount.GreaterThan(decimal.Zero) && a.Balance.GreaterThanOrEqual(amount)
}

func NewAccount(userId uuid.UUID, name string, owner string, currency currency.Currency) *Account {
	return &Account{
		UserId:   userId,
		Name:     name,
		Owner:    owner,
		Iban:     iban.New(),
		Currency: currency,
		Balance:  decimal.Zero,
		Status:   valobj.AccountStatusActive,
	}
}
