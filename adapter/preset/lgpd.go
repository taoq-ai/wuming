package preset

import "github.com/taoq-ai/wuming/domain/model"

func init() {
	register(Preset{
		Name:        "lgpd",
		Description: "Brazil Lei Geral de Proteção de Dados — covers all personal data",
		Locales:     []string{"common", "br"},
		PIITypes:    allPIITypes(),
		MinSeverity: model.Low,
	})
}
