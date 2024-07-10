package events

import (
	"github.com/9ssi7/banking/assets"
	"github.com/9ssi7/banking/internal/app/messages"
	"github.com/9ssi7/banking/internal/domain/valobj"
	"github.com/9ssi7/banking/internal/infra/mail"
)

type AuthStarted struct {
	Email  string
	Code   string
	Device valobj.Device
}

func OnAuthStarted(e AuthStarted) {
	go func() {
		mail.GetClient().SendWithTemplate(mail.SendWithTemplateConfig{
			SendConfig: mail.SendConfig{
				To:      []string{e.Email},
				Subject: messages.AuthVerifyEmailSubject,
				Message: e.Code,
			},
			Template: assets.Templates.AuthVerify,
			Data: map[string]interface{}{
				"Code":    e.Code,
				"IP":      mail.GetField(e.Device.IP),
				"Browser": mail.GetField(e.Device.Name),
				"OS":      mail.GetField(e.Device.OS),
			},
		})

	}()
}
