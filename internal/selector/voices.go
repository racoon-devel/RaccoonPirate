package selector

import "strings"

type Voices [][]string

func (v *Voices) Append(names ...string) {
	var normalized []string
	for _, n := range names {
		normalized = append(normalized, strings.ToLower(n))
	}
	*v = append(*v, normalized)
}
