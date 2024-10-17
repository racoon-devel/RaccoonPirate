package selector

import (
	"sort"
	"strings"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/antzucaro/matchr"
	"github.com/apex/log"
)

func (s MediaSelector) getMovieRankFunction(criteria Criteria) rankFunc {
	switch criteria {
	case CriteriaQuality:
		return makeRankFunc(s.limitBySize, s.rankByQuality, s.rankWeight(2, s.getRankByVoiceFunc()))
	case CriteriaFastest:
		return makeRankFunc(s.rankBySize, s.rankBySeeders, s.rankWeight(0.5, s.getRankByVoiceFunc()))
	case CriteriaCompact:
		return makeRankFunc(s.limitBySize, s.rankBySeeders, s.rankByQuality, s.rankWeight(2, s.rankBySeasons), s.rankWeight(2, s.getRankByVoiceFunc()))
	}
	panic("unknown criteria")
}

func (s MediaSelector) SelectMovie(l *log.Entry, criteria Criteria, list []*models.SearchTorrentsResult) *models.SearchTorrentsResult {
	rank := s.getMovieRankFunction(criteria)
	ranks := rank(l, list)
	_, _, best := findMax(ranks, func(elem float32) float32 {
		return elem
	})
	for i := range ranks {
		l.Debugf("%d rank: %.4f", i, ranks[i])
	}
	sel := list[best]
	l.Infof("Selected { Title: %s, Voice: %s, Size: %d, Seeders: %d, Quality: %s }", getString(sel.Title), sel.Voice, getValue(sel.Size), getValue(sel.Seeders), sel.Quality)
	return sel
}

func (s MediaSelector) SortMovies(l *log.Entry, criteria Criteria, list []*models.SearchTorrentsResult) {
	rank := s.getMovieRankFunction(criteria)
	ranks := rank(l, list)
	sort.SliceStable(list, func(i, j int) bool {
		return ranks[j] < ranks[i]
	})
}

func (s MediaSelector) rankBySeasons(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	_, max, _ := findMax(list, func(t *models.SearchTorrentsResult) int {
		return len(t.Seasons)
	})

	for i, t := range list {
		ranks[i] = float32(len(t.Seasons)) / float32(max)
	}

	return ranks
}

func (s MediaSelector) getRankByVoiceFunc() rankFunc {
	return func(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
		if s.Voice == "" {
			return s.rankByVoiceList(l, list)
		}
		return s.rankByVoice(l, list)
	}
}

func (s MediaSelector) rankByVoiceList(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	perItemWeight := 1 / float32(len(s.VoiceList))
	for i, t := range list {
		tVoice := strings.ToLower(t.Voice)
	ScanVoice:
		for j, voice := range s.VoiceList {
			for _, w := range voice {
				if strings.Index(tVoice, w) >= 0 {
					ranks[i] = float32(len(s.VoiceList)-j) * perItemWeight
					l.Debugf("%d rank by voice list: %.4f", i, ranks[i])
					break ScanVoice
				}
			}
		}
	}
	return ranks
}

func (s MediaSelector) rankByVoice(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
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
		l.Debugf("%d rank by voice: %.4f", j, ranks[j])
	}
	return ranks
}
