package seeds

import (
	"time"

	"github.com/9ssi7/banking/config/roles"
	"github.com/9ssi7/banking/internal/domain/entities"
	"github.com/9ssi7/banking/pkg/ptr"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func runUserSeeds(db *gorm.DB) {
	var user entities.User
	if err := db.Model(&entities.User{}).Where("email = ?", "test@test.com").First(&user).Error; err != nil {
		db.Create(&entities.User{
			Email: "test@test.com",
			Name:  "Test",
			Roles: pq.StringArray{
				roles.Admin,
				roles.AdminSuper,
			},
			IsActive:   true,
			VerifiedAt: ptr.Time(time.Now()),
		})
	}
}
