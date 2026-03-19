package replacer

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/taoq-ai/wuming/domain/model"
)

// Hash replaces each match with a deterministic SHA-256 hash (truncated).
// The same input value always produces the same hash output.
type Hash struct {
	// Length is the number of hex characters to keep. Defaults to 16.
	Length int
	// Salt is an optional salt prepended to the value before hashing.
	Salt string
}

// NewHash creates a Hash replacer with default settings (16 hex chars, no salt).
func NewHash() *Hash {
	return &Hash{Length: 16}
}

func (h *Hash) Name() string { return "hash" }

func (h *Hash) Replace(text string, matches []model.Match) (string, error) {
	if len(matches) == 0 {
		return text, nil
	}

	sorted := make([]model.Match, len(matches))
	copy(sorted, matches)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start > sorted[j].Start
	})

	result := []byte(text)
	for _, m := range sorted {
		if m.Start < 0 || m.End > len(result) || m.Start >= m.End {
			continue
		}
		sum := sha256.Sum256([]byte(h.Salt + m.Value))
		hashed := hex.EncodeToString(sum[:])
		length := h.Length
		if length > len(hashed) {
			length = len(hashed)
		}
		result = append(result[:m.Start], append([]byte(hashed[:length]), result[m.End:]...)...)
	}
	return string(result), nil
}
