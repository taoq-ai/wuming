package ca

import "github.com/taoq-ai/wuming/domain/port"

// All returns all Canada locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewSINDetector(),
		NewPhoneDetector(),
		NewPostalCodeDetector(),
		NewPassportDetector(),
	}
}
