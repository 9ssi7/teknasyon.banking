package repos

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepo struct {
	syncRepo
	txnGormRepo
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) abstracts.TransactionRepo {
	return &transactionRepo{
		db:          db,
		txnGormRepo: newTxnGormRepo(db),
	}
}

func (r *transactionRepo) Save(ctx context.Context, transaction *entities.Transaction) error {
	r.syncRepo.Lock()
	defer r.syncRepo.Unlock()
	if err := r.adapter.GetCurrent(ctx).Save(transaction).Error; err != nil {
		return rescode.Failed
	}
	return nil
}

func (r *transactionRepo) Filter(ctx context.Context, accountId uuid.UUID, pagi *list.PagiRequest, filters *valobj.TransactionFilters) (*list.PagiResponse[*entities.Transaction], error) {
	var transactions []*entities.Transaction
	query := r.adapter.GetCurrent(ctx).Model(&entities.Transaction{}).Where("sender_id = ? OR receiver_id = ?", accountId, accountId)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, rescode.Failed
	}
	if filters.StartDate != "" {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != "" {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	if filters.Kind != "" {
		query = query.Where("kind = ?", filters.Kind)
	}
	var filteredTotal int64
	if err := query.Count(&filteredTotal).Error; err != nil {
		return nil, rescode.Failed
	}
	if err := query.Limit(*pagi.Limit).Offset(pagi.Offset()).Find(&transactions).Error; err != nil {
		return nil, rescode.Failed
	}
	return &list.PagiResponse[*entities.Transaction]{
		List:          transactions,
		Total:         total,
		Limit:         *pagi.Limit,
		TotalPage:     pagi.TotalPage(filteredTotal),
		FilteredTotal: filteredTotal,
		Page:          *pagi.Page,
	}, nil
}
