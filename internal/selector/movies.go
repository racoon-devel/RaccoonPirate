package selector

import (
	"strings"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/antzucaro/matchr"
)

func (s selection) getMovieRankFunction() rankFunc {
	switch s.Criteria {
	case CriteriaQuality:
		return makeRankFunc(s.limitBySize, s.rankByQuality, s.rankWeight(2, s.getRankByVoiceFunc()))
	case CriteriaFastest:
		return makeRankFunc(s.rankBySize, s.rankBySeeders, s.rankWeight(0.5, s.getRankByVoiceFunc()))
	case CriteriaCompact:
		return makeRankFunc(s.limitBySize, s.rankBySeeders, s.rankByQuality, s.rankWeight(2, s.rankBySeasons), s.rankWeight(2, s.getRankByVoiceFunc()))
	}
	panic("unknown criteria")
}

func (s selection) rankBySeasons(list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	_, max, _ := findMax(list, func(t *models.SearchTorrentsResult) int {
		return len(t.Seasons)
	})

	for i, t := range list {
		ranks[i] = float32(len(t.Seasons)) / float32(max)
	}

	return ranks
}

func (s selection) getRankByVoiceFunc() rankFunc {
	return func(list []*models.SearchTorrentsResult) []float32 {
		if s.Voice == "" {
			return s.rankByVoiceList(list)
		}
		return s.rankByVoice(list)
	}
}

func (s selection) rankByVoiceList(list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	perItemWeight := 1 / float32(len(s.VoiceList))
	for i, t := range list {
		tVoice := strings.ToLower(t.Voice)
	ScanVoice:
		for j, voice := range s.VoiceList {
			for _, w := range voice {
				if strings.Index(tVoice, w) >= 0 {
					ranks[i] = float32(len(s.VoiceList)-j) * perItemWeight
					s.log().Debugf("%d rank by voice list: %.4f", i, ranks[i])
					break ScanVoice
				}
			}
		}
	}
	return ranks
}

func (s selection) rankByVoice(list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	distance := make([]int, len(list))

	target := strings.ToLower(s.Voice)
	for i, t := range list {
		voice := strings.ToLower(t.Voice)
		distance[i] = matchr.Levenshtein(voice, target)
	}

	_, max, _ := findMax(distance, func(elem int) int {
		return elem
	})
	for j, d := range distance {
		ranks[j] = 1 - (float32(d) / float32(max))
		s.log().Debugf("%d rank by voice: %.4f", j, ranks[j])
	}
	return ranks
}
