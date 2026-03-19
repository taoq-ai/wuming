package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "pipa",
		Description: "South Korea Personal Information Protection Act — identity and contact data",
		Locales:     []string{"common", "kr"},
		PIITypes: []model.PIIType{
			model.NationalID,
			model.Phone,
		},
		MinSeverity: model.Medium,
	})
}
