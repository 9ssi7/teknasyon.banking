package iban

import (
	"fmt"
	"time"

	"math/rand"
)

func New() string {
	return generateIBAN("tr") // add i18n support
}

func Validate(iban string) bool {
	return validateIBAN(iban)
}

func validateIBAN(iban string) bool {
	if len(iban) != 26 && len(iban) != 27 {
		return false
	}
	for _, char := range iban[:2] {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	for _, char := range iban[2:] {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func generateIBAN(countryCode string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var ibanFormat string
	switch countryCode {
	case "en":
		ibanFormat = "US%02d%04d%017d"
	default:
		ibanFormat = "TR%02d%04d%016d"
	}
	iban := fmt.Sprintf(ibanFormat,
		rand.Intn(100),               // Bank code (random)
		rand.Intn(10000),             // Branch code (random)
		rand.Intn(10000000000000000), // Account number (random)
	)
	iban += generateCheckDigits(iban)
	return iban
}

func generateCheckDigits(iban string) string {
	iban = iban[4:] + iban[:4]
	numericIBAN := ""
	for _, char := range iban {
		if char >= 'A' && char <= 'Z' {
			numericIBAN += fmt.Sprintf("%d", int(char-'A')+10)
		} else {
			numericIBAN += string(char)
		}
	}
	remainder := mod97(numericIBAN)
	return fmt.Sprintf("%02d", 98-remainder)
}

func mod97(number string) int {
	var remainder int
	for _, char := range number {
		remainder = (remainder*10 + int(char-'0')) % 97
	}
	return remainder
}
