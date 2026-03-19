package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "dpdp",
		Description: "India Digital Personal Data Protection Act — covers all personal data",
		Locales:     []string{"common", "in"},
		PIITypes:    allPIITypes(),
		MinSeverity: model.Low,
	})
}
