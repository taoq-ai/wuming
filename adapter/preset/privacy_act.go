package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "privacy-act",
		Description: "Australia Privacy Act — tax, health, and contact data",
		Locales:     []string{"common", "au"},
		PIITypes: []model.PIIType{
			model.TaxID,
			model.HealthID,
			model.Phone,
			model.Email,
		},
		MinSeverity: model.Medium,
	})
}
