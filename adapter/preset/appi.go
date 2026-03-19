package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "appi",
		Description: "Japan Act on the Protection of Personal Information — covers all personal data",
		Locales:     []string{"common", "jp"},
		PIITypes:    allPIITypes(),
		MinSeverity: model.Low,
	})
}
