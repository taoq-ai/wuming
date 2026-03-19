package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "pci-dss",
		Description: "Payment Card Industry Data Security Standard — credit card data protection",
		Locales:     []string{"common"},
		PIITypes: []model.PIIType{
			model.CreditCard,
		},
		MinSeverity: model.Critical,
	})
}
