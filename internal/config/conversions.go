package config

import "github.com/RacoonMediaServer/rms-library/pkg/selector"

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
