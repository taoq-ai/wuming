package nl

import "github.com/taoq-ai/wuming/domain/port"

// All returns all Dutch (NL) locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewBSNDetector(),
		NewPhoneDetector(),
		NewPostalDetector(),
		NewKvKDetector(),
		NewIDDocumentDetector(),
	}
}
