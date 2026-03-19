package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "gdpr",
		Description: "EU General Data Protection Regulation — covers all personal data across EU/EEA locales",
		Locales:     []string{"common", "eu", "nl", "de", "fr", "gb"},
		PIITypes:    allPIITypes(),
		MinSeverity: model.Low,
	})
}
