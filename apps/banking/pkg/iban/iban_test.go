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

func TestIbanValidate(t *testing.T) {
	type args struct {
		iban     string
		expected bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test TR IBAN",
			args: args{
				iban:     "TR298813938771869288426798",
				expected: true,
			},
		},
		{
			name: "Test EN IBAN",
			args: args{
				iban:     "US021000021",
				expected: false,
			},
		},
		{
			name: "Test Invalid IBAN",
			args: args{
				iban:     "TR330006100519786457841327",
				expected: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateIBAN(tt.args.iban); got != tt.args.expected {
				t.Errorf("validateIBAN() = %v, want %v", got, tt.args.expected)
			}
		})
	}
}
