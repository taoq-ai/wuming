package jp

import "github.com/taoq-ai/wuming/domain/port"

// All returns all Japan locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewMyNumberDetector(),
		NewCorporateNumberDetector(),
		NewPhoneDetector(),
		NewPostalDetector(),
		NewPassportDetector(),
	}
}
