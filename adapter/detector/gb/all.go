package gb

import "github.com/taoq-ai/wuming/domain/port"

// All returns all British (GB) locale PII detectors.
func All() []port.Detector {
	return []port.Detector{
		NewNINDetector(),
		NewNHSDetector(),
		NewUTRDetector(),
		NewPhoneDetector(),
		NewPostcodeDetector(),
	}
}
