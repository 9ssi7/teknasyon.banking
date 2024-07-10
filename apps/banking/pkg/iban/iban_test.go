package iban

import (
	"testing"
)

func TestIbanGeneration(t *testing.T) {
	type args struct {
		countryCode string
		expectedLen int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test TR IBAN",
			args: args{
				countryCode: "tr",
				expectedLen: 26,
			},
		},
		{
			name: "Test EN IBAN",
			args: args{
				countryCode: "en",
				expectedLen: 27,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iban := generateIBAN(tt.args.countryCode)
			if len(iban) != tt.args.expectedLen {
				t.Errorf("generateIBAN() = %v, want %v", len(iban), tt.args.expectedLen)
			}
		})
	}
}
