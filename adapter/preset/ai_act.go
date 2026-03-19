package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "ai-act",
		Description: "EU AI Act — scrub training/validation data for high-risk AI systems (Articles 10, 15)",
		Locales: []string{
			"au", "br", "ca", "cn", "common", "de", "eu",
			"fr", "gb", "in", "jp", "kr", "nl", "us",
		},
		PIITypes:    allPIITypes(),
		MinSeverity: model.Low,
	})
}
