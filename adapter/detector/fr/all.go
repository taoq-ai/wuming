package fr

import "github.com/taoq-ai/wuming/domain/port"

// All returns all French (FR) locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewNIRDetector(),
		NewNIFDetector(),
		NewPhoneDetector(),
		NewPostalDetector(),
		NewIDCardDetector(),
	}
}
