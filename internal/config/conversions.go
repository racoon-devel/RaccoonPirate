package config

import "github.com/racoon-devel/raccoon-pirate/internal/selector"

func (c Selector) GetCriterion() selector.Criteria {
	switch c.Criterion {
	case "quality":
		return selector.CriteriaQuality
	case "fastest":
		return selector.CriteriaFastest
	case "compact":
		return selector.CriteriaCompact
	default:
		return selector.CriteriaQuality
	}
}
