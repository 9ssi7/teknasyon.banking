package abstracts

import (
	"context"

	"github.com/9ssi7/banking/internal/domain/aggregates"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/list"
	"github.com/google/uuid"
)

type UserRepo interface {
	Save(ctx context.Context, user *entities.User) error
	IsExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByToken(ctx context.Context, token string) (*entities.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByPhone(ctx context.Context, phone string) (*entities.User, error)
	Filter(ctx context.Context, req *list.PagiRequest, search string, isActive string) (*list.PagiResponse[*entities.User], error)
}

type AccountRepo interface {
	Save(ctx context.Context, account *entities.Account) error
	ListByUserId(ctx context.Context, userId uuid.UUID, pagi *list.PagiRequest) (*list.PagiResponse[*entities.Account], error)
	FindByIbanAndOwner(ctx context.Context, iban string, owner string) (*entities.Account, error)
	FindByUserIdAndId(ctx context.Context, userId uuid.UUID, id uuid.UUID) (*entities.Account, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.Account, error)
}

type TransactionRepo interface {
	Save(ctx context.Context, transaction *entities.Transaction) error
	Filter(ctx context.Context, accountId uuid.UUID, pagi *list.PagiRequest, filters *valobj.TransactionFilters) (*list.PagiResponse[*entities.Transaction], error)
}

type SessionRepo interface {
	Save(ctx context.Context, userId uuid.UUID, session *aggregates.Session) error
	FindByIds(ctx context.Context, userId uuid.UUID, deviceId string) (*aggregates.Session, error)
	FindAllByUserId(ctx context.Context, userId uuid.UUID) ([]*aggregates.Session, error)
	Destroy(ctx context.Context, userId uuid.UUID, deviceId string) error
}

type VerifyRepo interface {
	Save(ctx context.Context, token string, verify *aggregates.Verify) error
	IsExists(ctx context.Context, token string, deviceId string) (bool, error)
	Find(ctx context.Context, token string, deviceId string) (*aggregates.Verify, error)
	Delete(ctx context.Context, token string, deviceId string) error
}

type Repositories struct {
	VerifyRepo      VerifyRepo
	SessionRepo     SessionRepo
	UserRepo        UserRepo
	AccountRepo     AccountRepo
	TransactionRepo TransactionRepo
}
