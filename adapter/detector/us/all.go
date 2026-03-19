package us

import "github.com/taoq-ai/wuming/domain/port"

// All returns all US locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewSSNDetector(),
		NewITINDetector(),
		NewEINDetector(),
		NewPhoneDetector(),
		NewZIPDetector(),
		NewPassportDetector(),
		NewMedicareDetector(),
	}
}
