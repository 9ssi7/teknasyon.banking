package repos

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/9ssi7/banking/pkg/rescode"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountRepo struct {
	syncRepo
	txnGormRepo
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) abstracts.AccountRepo {
	return &accountRepo{
		db:          db,
		txnGormRepo: newTxnGormRepo(db),
	}
}

func (r *accountRepo) Save(ctx context.Context, account *entities.Account) error {
	r.syncRepo.Lock()
	defer r.syncRepo.Unlock()
	if err := r.adapter.GetCurrent(ctx).Save(account).Error; err != nil {
		return rescode.Failed
	}
	return nil
}

func (r *accountRepo) ListByUserId(ctx context.Context, userId uuid.UUID, pagi *list.PagiRequest) (*list.PagiResponse[*entities.Account], error) {
	var accounts []*entities.Account
	query := r.adapter.GetCurrent(ctx).Model(&entities.Account{}).Where("user_id = ?", userId)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, rescode.Failed
	}
	var filteredTotal int64
	if err := query.Count(&filteredTotal).Error; err != nil {
		return nil, rescode.Failed
	}
	if err := query.Limit(*pagi.Limit).Offset(pagi.Offset()).Find(&accounts).Error; err != nil {
		return nil, rescode.Failed
	}
	return &list.PagiResponse[*entities.Account]{
		List:          accounts,
		Total:         total,
		Limit:         *pagi.Limit,
		TotalPage:     pagi.TotalPage(filteredTotal),
		FilteredTotal: filteredTotal,
		Page:          *pagi.Page,
	}, nil
}

func (r *accountRepo) FindByIbanAndOwner(ctx context.Context, iban string, owner string) (*entities.Account, error) {
	var account entities.Account
	if err := r.adapter.GetCurrent(ctx).Model(&entities.Account{}).Where("iban = ? AND owner = ?", iban, owner).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rescode.NotFound
		}
		return nil, rescode.Failed
	}
	return &account, nil
}

func (r *accountRepo) FindByUserIdAndId(ctx context.Context, userId uuid.UUID, id uuid.UUID) (*entities.Account, error) {
	var account entities.Account
	if err := r.adapter.GetCurrent(ctx).Model(&entities.Account{}).Where("user_id = ? AND id = ?", userId, id).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rescode.NotFound
		}
		return nil, rescode.Failed
	}
	return &account, nil
}

func (r *accountRepo) FindById(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	var account entities.Account
	if err := r.adapter.GetCurrent(ctx).Model(&entities.Account{}).Where("id = ?", id).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rescode.NotFound
		}
		return nil, rescode.Failed
	}
	return &account, nil
}
