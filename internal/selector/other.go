package selector

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"

func (s selection) getOtherRankFunction() rankFunc {
	return makeRankFunc(s.rankBySeeders, s.getRankByTextOtherFunc())
}

func (s selection) getRankByTextOtherFunc() rankFunc {
	return func(list []*models.SearchTorrentsResult) []float32 {
		return s.rankByText(s.Query, list)
	}
}
