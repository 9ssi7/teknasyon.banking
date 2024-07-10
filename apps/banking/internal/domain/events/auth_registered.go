package events

import (
	"github.com/9ssi7/banking/assets"
	"github.com/9ssi7/banking/internal/app/messages"
	"github.com/9ssi7/banking/internal/infra/mail"
)

type AuthRegistered struct {
	Name  string
	Email string
}

func OnAuthRegistered(e AuthRegistered) {
	go func() {
		mail.GetClient().SendWithTemplate(mail.SendWithTemplateConfig{
			SendConfig: mail.SendConfig{
				To:      []string{e.Email},
				Subject: messages.AuthRegisteredEmailSubject,
			},
			Template: assets.Templates.AuthRegistered,
			Data: map[string]interface{}{
				"Name": e.Name,
			},
		})
	}()
}
