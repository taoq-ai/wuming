package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "pipeda",
		Description: "Canada Personal Information Protection and Electronic Documents Act",
		Locales:     []string{"common", "ca"},
		PIITypes: []model.PIIType{
			model.NationalID,
			model.HealthID,
			model.Phone,
			model.Email,
		},
		MinSeverity: model.Medium,
	})
}
