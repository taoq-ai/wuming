package common

import "github.com/taoq-ai/wuming/domain/port"

// All returns all common (locale-independent) PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewEmailDetector(),
		NewCreditCardDetector(),
		NewIPDetector(),
		NewURLDetector(),
		NewIBANDetector(),
		NewMACDetector(),
	}
}
