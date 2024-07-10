package events

import (
	"fmt"

	"github.com/9ssi7/banking/assets"
	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/internal/app/messages"
	"github.com/9ssi7/banking/internal/infra/mail"
)

type AuthRegistered struct {
	Name             string
	Email            string
	VerificationCode string
}

func OnAuthRegistered(e AuthRegistered) {
	go func() {
		mail.GetClient().SendWithTemplate(mail.SendWithTemplateConfig{
			SendConfig: mail.SendConfig{
				To:      []string{e.Email},
				Subject: messages.AuthEmailSubjectRegistered,
			},
			Template: assets.Templates.AuthRegistered,
			Data: map[string]interface{}{
				"Name":            e.Name,
				"VerificationUrl": fmt.Sprintf("%s/auth/verify/%s", config.ReadValue().PublicHost, e.VerificationCode),
			},
		})
	}()
}
