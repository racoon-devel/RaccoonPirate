package selector

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/heuristic"
	"github.com/antzucaro/matchr"
)

func (s selection) getMusicRankFunction() rankFunc {
	funcs := []rankFunc{s.rankBySeeders}
	funcs = append(funcs, s.getRankByTextFunc()...)
	return makeRankFunc(funcs...)
}

func (s selection) getRankByTextFunc() []rankFunc {
	var result []rankFunc
	if s.Discography {
		result = append(result, func(list []*models.SearchTorrentsResult) []float32 {
			return s.rankByText(s.Query+" дискография", list)
		})
		result = append(result, func(list []*models.SearchTorrentsResult) []float32 {
			return s.rankByText(s.Query+" discography", list)
		})
	} else {
		result = append(result, func(list []*models.SearchTorrentsResult) []float32 {
			return s.rankByText(s.Query, list)
		})
	}
	return result
}

func (s selection) rankByText(query string, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	distance := make([]int, len(list))

	target := heuristic.NormalizeWithoutBraces(query)
	for i, t := range list {
		title := heuristic.NormalizeWithoutBraces(*t.Title)
		distance[i] = matchr.Levenshtein(title, target)
	}

	_, max, _ := findMax(distance, func(elem int) int {
		return elem
	})
	for j, d := range distance {
		ranks[j] = 1 - (float32(d) / float32(max))
	}
	return ranks
}
