package kr

import "github.com/taoq-ai/wuming/domain/port"

// All returns all South Korea locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewRRNDetector(),
		NewPhoneDetector(),
		NewPostalDetector(),
		NewPassportDetector(),
	}
}
