package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "pipl",
		Description: "China Personal Information Protection Law — covers all personal data",
		Locales:     []string{"common", "cn"},
		PIITypes:    allPIITypes(),
		MinSeverity: model.Low,
	})
}
