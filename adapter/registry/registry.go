// Package registry aggregates all PII detectors across locales,
// providing convenient access to the full set of available detectors.
package registry

import (
	"sort"

	"github.com/taoq-ai/wuming/adapter/detector/au"
	"github.com/taoq-ai/wuming/adapter/detector/cn"
	"github.com/taoq-ai/wuming/adapter/detector/common"
	"github.com/taoq-ai/wuming/adapter/detector/de"
	"github.com/taoq-ai/wuming/adapter/detector/eu"
	"github.com/taoq-ai/wuming/adapter/detector/fr"
	"github.com/taoq-ai/wuming/adapter/detector/gb"
	"github.com/taoq-ai/wuming/adapter/detector/jp"
	"github.com/taoq-ai/wuming/adapter/detector/kr"
	"github.com/taoq-ai/wuming/adapter/detector/nl"
	"github.com/taoq-ai/wuming/adapter/detector/us"
	"github.com/taoq-ai/wuming/domain/port"
)

// localeProviders maps locale names to their All() functions.
var localeProviders = map[string]func() []port.Detector{
	"au":     au.All,
	"cn":     cn.All,
	"common": common.All,
	"de":     de.All,
	"eu":     eu.All,
	"fr":     fr.All,
	"gb":     gb.All,
	"jp":     jp.All,
	"kr":     kr.All,
	"nl":     nl.All,
	"us":     us.All,
}

// AllDetectors returns every registered PII detector across all locales.
func AllDetectors() []port.Detector {
	var all []port.Detector
	for _, provider := range localeProviders {
		all = append(all, provider()...)
	}
	return all
}

// DetectorsForLocale returns all detectors for the given locale.
// Common/global detectors are always included.
func DetectorsForLocale(locale string) []port.Detector {
	detectors := common.All()
	if locale == "common" {
		return detectors
	}
	if provider, ok := localeProviders[locale]; ok {
		detectors = append(detectors, provider()...)
	}
	return detectors
}

// Locales returns a sorted list of all supported locale names.
func Locales() []string {
	locales := make([]string, 0, len(localeProviders))
	for k := range localeProviders {
		locales = append(locales, k)
	}
	sort.Strings(locales)
	return locales
}
