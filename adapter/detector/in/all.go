package in

import "github.com/taoq-ai/wuming/domain/port"

// All returns all India locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewAadhaarDetector(),
		NewPANDetector(),
		NewPhoneDetector(),
		NewPINCodeDetector(),
		NewPassportDetector(),
		NewGSTINDetector(),
	}
}
