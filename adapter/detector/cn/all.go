package cn

import "github.com/taoq-ai/wuming/domain/port"

// All returns all China locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewResidentIDDetector(),
		NewPhoneDetector(),
		NewPostalDetector(),
		NewPassportDetector(),
		NewUSCCDetector(),
	}
}
