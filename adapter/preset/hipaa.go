package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "hipaa",
		Description: "US Health Insurance Portability and Accountability Act — protected health information",
		Locales:     []string{"common", "us"},
		PIITypes: []model.PIIType{
			model.NationalID,
			model.HealthID,
			model.Phone,
			model.Email,
			model.PostalCode,
			model.DateOfBirth,
		},
		MinSeverity: model.Medium,
	})
}
