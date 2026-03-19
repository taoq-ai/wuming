package model

import "testing"

func TestPIITypeString(t *testing.T) {
	tests := []struct {
		piiType PIIType
		want    string
	}{
		{Email, "EMAIL"},
		{Phone, "PHONE"},
		{CreditCard, "CREDIT_CARD"},
		{IBAN, "IBAN"},
		{IPAddress, "IP_ADDRESS"},
		{URL, "URL"},
		{MACAddress, "MAC_ADDRESS"},
		{NationalID, "NATIONAL_ID"},
		{TaxID, "TAX_ID"},
		{Passport, "PASSPORT"},
		{DriversLicense, "DRIVERS_LICENSE"},
		{HealthID, "HEALTH_ID"},
		{DateOfBirth, "DATE_OF_BIRTH"},
		{Name, "NAME"},
		{Address, "ADDRESS"},
		{PostalCode, "POSTAL_CODE"},
		{BankAccount, "BANK_ACCOUNT"},
		{SocialMedia, "SOCIAL_MEDIA"},
		{Custom, "CUSTOM"},
		{PIIType(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.piiType.String(); got != tt.want {
			t.Errorf("PIIType(%d).String() = %q, want %q", tt.piiType, got, tt.want)
		}
	}
}

func TestSeverityString(t *testing.T) {
	tests := []struct {
		severity Severity
		want     string
	}{
		{Low, "LOW"},
		{Medium, "MEDIUM"},
		{High, "HIGH"},
		{Critical, "CRITICAL"},
		{Severity(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.severity.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.severity, got, tt.want)
		}
	}
}
