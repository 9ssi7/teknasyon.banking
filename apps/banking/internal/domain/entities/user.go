package entities

import (
	"time"

	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/pkg/ptr"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	Base
	valobj.Audit
	Name       string         `json:"name" gorm:"type:varchar(255);not null"`
	Email      string         `json:"email" gorm:"type:varchar(255);not null;unique"`
	IsActive   bool           `json:"is_active" gorm:"type:boolean;not null;default:true"`
	Roles      pq.StringArray `json:"roles" gorm:"type:text[]"`
	TempToken  *string        `json:"temp_token" gorm:"type:varchar(255);default:null;index:idx_verifier"`
	VerifiedAt *time.Time     `json:"verified_at" gorm:"type:timestamp;default:null;index:idx_verifier"`
}

func (u *User) AddRole(role string) {
	u.Roles = append(u.Roles, role)
}

func (u *User) RemoveRole(role string) {
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			break
		}
	}
}

func (u *User) Verify() {
	u.VerifiedAt = ptr.Time(time.Now())
	u.TempToken = nil
}

func (u *User) Enable() {
	u.IsActive = true
}

func (u *User) Disable() {
	u.IsActive = false
}

func NewUser(name string, email string) *User {
	return &User{
		Name:      name,
		Email:     email,
		Roles:     pq.StringArray{},
		TempToken: ptr.String(uuid.New().String()),
	}
}

func NewUserFromAdmin(name string, email string, adminId uuid.UUID) *User {
	return &User{
		Audit: valobj.Audit{
			MakedBy: &adminId,
		},
		Name:      name,
		Email:     email,
		Roles:     pq.StringArray{},
		TempToken: ptr.String(uuid.New().String()),
	}
}
