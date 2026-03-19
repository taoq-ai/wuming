package au

import "github.com/taoq-ai/wuming/domain/port"

// All returns all Australian locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewTFNDetector(),
		NewMedicareDetector(),
		NewABNDetector(),
		NewPhoneDetector(),
		NewPostcodeDetector(),
	}
}
