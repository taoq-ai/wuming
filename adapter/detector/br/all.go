package br

import "github.com/taoq-ai/wuming/domain/port"

// All returns all BR locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewCPFDetector(),
		NewCNPJDetector(),
		NewPhoneDetector(),
		NewCEPDetector(),
		NewPISDetector(),
		NewCNHDetector(),
	}
}
