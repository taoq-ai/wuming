package cn

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/taoq-ai/wuming/domain/model"
)

// Resident ID (居民身份证): 18 characters.
// 6-digit area code + 8-digit birth date (YYYYMMDD) + 3-digit sequence + 1 check digit (0-9 or X).
var residentIDRe = regexp.MustCompile(`\b[1-9]\d{5}(?:19|20)\d{2}(?:0[1-9]|1[0-2])(?:0[1-9]|[12]\d|3[01])\d{3}[\dX]\b`)

// weights for mod-11 check digit calculation.
var residentIDWeights = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// checkMap maps the mod-11 remainder to the expected check character.
var residentIDCheckMap = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

// ResidentIDDetector detects Chinese Resident Identity Card numbers.
type ResidentIDDetector struct{}

func NewResidentIDDetector() *ResidentIDDetector { return &ResidentIDDetector{} }

func (d *ResidentIDDetector) Name() string              { return "cn/resident_id" }
func (d *ResidentIDDetector) Locales() []string         { return []string{locale} }
func (d *ResidentIDDetector) PIITypes() []model.PIIType { return []model.PIIType{model.NationalID} }

func (d *ResidentIDDetector) Detect(_ context.Context, text string) ([]model.Match, error) {
	results := residentIDRe.FindAllStringIndex(text, -1)
	if len(results) == 0 {
		return nil, nil
	}

	var matches []model.Match
	for _, loc := range results {
		id := text[loc[0]:loc[1]]
		if !isValidResidentID(id) {
			continue
		}
		matches = append(matches, model.Match{
			Type:       model.NationalID,
			Value:      id,
			Start:      loc[0],
			End:        loc[1],
			Confidence: 0.90,
			Locale:     locale,
			Detector:   d.Name(),
		})
	}
	return matches, nil
}

func isValidResidentID(id string) bool {
	upper := strings.ToUpper(id)

	// Verify check digit using mod-11 algorithm.
	sum := 0
	for i := 0; i < 17; i++ {
		digit, err := strconv.Atoi(string(upper[i]))
		if err != nil {
			return false
		}
		sum += digit * residentIDWeights[i]
	}
	remainder := sum % 11
	if upper[17] != residentIDCheckMap[remainder] {
		return false
	}

	// Validate birth date is reasonable: not in the future and not before 1900.
	year, _ := strconv.Atoi(upper[6:10])
	month, _ := strconv.Atoi(upper[10:12])
	day, _ := strconv.Atoi(upper[12:14])
	birthDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if birthDate.Year() != year || int(birthDate.Month()) != month || birthDate.Day() != day {
		return false // invalid date (e.g. Feb 30)
	}
	if birthDate.After(time.Now()) {
		return false
	}

	return true
}
