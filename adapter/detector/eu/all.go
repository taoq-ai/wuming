package eu

import "github.com/taoq-ai/wuming/domain/port"

// All returns all EU locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewPassportMRZDetector(),
		NewVATDetector(),
	}
}
