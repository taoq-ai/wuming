// Package model defines the core domain types for PII detection and redaction.
package model

// PIIType categorizes a detected PII entity.
type PIIType int

const (
	// Global patterns.
	Email PIIType = iota + 1
	Phone
	CreditCard
	IBAN
	IPAddress
	URL
	MACAddress

	// Identity documents.
	NationalID
	TaxID
	Passport
	DriversLicense
	HealthID

	// Personal information.
	DateOfBirth
	Name
	Address
	PostalCode

	// Financial.
	BankAccount

	// Online.
	SocialMedia

	// User-defined.
	Custom
)

var piiTypeNames = map[PIIType]string{
	Email:          "EMAIL",
	Phone:          "PHONE",
	CreditCard:     "CREDIT_CARD",
	IBAN:           "IBAN",
	IPAddress:      "IP_ADDRESS",
	URL:            "URL",
	MACAddress:     "MAC_ADDRESS",
	NationalID:     "NATIONAL_ID",
	TaxID:          "TAX_ID",
	Passport:       "PASSPORT",
	DriversLicense: "DRIVERS_LICENSE",
	HealthID:       "HEALTH_ID",
	DateOfBirth:    "DATE_OF_BIRTH",
	Name:           "NAME",
	Address:        "ADDRESS",
	PostalCode:     "POSTAL_CODE",
	BankAccount:    "BANK_ACCOUNT",
	SocialMedia:    "SOCIAL_MEDIA",
	Custom:         "CUSTOM",
}

// String returns the human-readable name of the PII type.
func (p PIIType) String() string {
	if name, ok := piiTypeNames[p]; ok {
		return name
	}
	return "UNKNOWN"
}

// Severity indicates how sensitive a piece of PII is.
type Severity int

const (
	// Low sensitivity — public-ish data (e.g. postal code).
	Low Severity = iota + 1
	// Medium sensitivity — semi-sensitive (e.g. phone number).
	Medium
	// High sensitivity — highly sensitive (e.g. SSN, health ID).
	High
	// Critical sensitivity — regulated data (e.g. credit card under PCI-DSS).
	Critical
)

var severityNames = map[Severity]string{
	Low:      "LOW",
	Medium:   "MEDIUM",
	High:     "HIGH",
	Critical: "CRITICAL",
}

// String returns the human-readable name of the severity level.
func (s Severity) String() string {
	if name, ok := severityNames[s]; ok {
		return name
	}
	return "UNKNOWN"
}

// Match represents a single PII detection result within a text.
type Match struct {
	// Type is the category of PII detected.
	Type PIIType
	// Value is the matched text.
	Value string
	// Start is the byte offset where the match begins.
	Start int
	// End is the byte offset where the match ends (exclusive).
	End int
	// Confidence is a score from 0.0 to 1.0 indicating detection certainty.
	Confidence float64
	// Locale identifies which locale this match belongs to (e.g. "nl", "us").
	// Empty string means the pattern is locale-independent.
	Locale string
	// Detector is the name of the detector that found this match.
	Detector string
}
