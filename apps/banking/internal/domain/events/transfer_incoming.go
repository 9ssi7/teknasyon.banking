package events

import (
	"fmt"

	"github.com/9ssi7/banking/assets"
	"github.com/9ssi7/banking/internal/app/messages"
	"github.com/9ssi7/banking/internal/infra/mail"
)

type TranfserIncoming struct {
	Email       string
	Name        string
	Amount      string
	Currency    string
	Account     string
	Description string
}

func OnTransferIncoming(e TranfserIncoming) {
	go func() {
		mail.GetClient().SendWithTemplate(mail.SendWithTemplateConfig{
			SendConfig: mail.SendConfig{
				To:      []string{e.Email},
				Subject: messages.TransactionEmailSubjectIncoming,
			},
			Template: assets.Templates.TransferIncoming,
			Data: map[string]interface{}{
				"Name":        e.Name,
				"Amount":      fmt.Sprintf("%s %s", mail.GetField(e.Amount), e.Currency),
				"Account":     mail.GetField(e.Account),
				"Description": mail.GetField(e.Description),
			},
		})

	}()
}
