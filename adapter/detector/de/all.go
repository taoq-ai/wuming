package de

import "github.com/taoq-ai/wuming/domain/port"

// All returns all German (DE) locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewSteuerIDDetector(),
		NewSozialversicherungDetector(),
		NewPhoneDetector(),
		NewPLZDetector(),
		NewIDCardDetector(),
	}
}
