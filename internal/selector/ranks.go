package selector

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/apex/log"
)

func (s MediaSelector) rankWeight(weight float32, f rankFunc) rankFunc {
	return func(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
		ranks := f(l, list)
		for i := range ranks {
			ranks[i] = weight * ranks[i]
		}
		return ranks
	}
}

func (s MediaSelector) rankBySize(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	_, max, _ := findMax(list, func(t *models.SearchTorrentsResult) int64 {
		return getValue(t.Size)
	})

	for i, t := range list {
		ranks[i] = 1 - (float32(getValue(t.Size)) / float32(max))
		l.Debugf("%d rank by size: %.4f", i, ranks[i])
	}
	return ranks
}

func (s MediaSelector) limitBySize(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))

	for i, t := range list {
		size := getValue(t.Size)
		seasons := int64(len(t.Seasons))
		if seasons == 0 {
			seasons = 1
		}

		if size < s.MinSeasonSizeMB*seasons || size >= s.MaxSeasonSizeMB*seasons {
			ranks[i] = -1
			l.Debugf("%d limit by size: %.4f", i, ranks[i])
		}

	}
	return ranks
}

func (s MediaSelector) rankBySeeders(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	_, max, _ := findMax(list, func(t *models.SearchTorrentsResult) int64 {
		return getValue(t.Seeders)
	})

	for i, t := range list {
		seeders := getValue(t.Seeders)
		if seeders < s.MinSeedersThreshold {
			ranks[i] = float32(seeders) / float32(max)
		} else {
			ranks[i] = 1
		}
		l.Debugf("%d rank by seeders: %.4f", i, ranks[i])
	}
	return ranks
}

func (s MediaSelector) rankByQuality(l *log.Entry, list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	perQualityWeight := 1 / float32(len(s.QualityPrior))
	for i, t := range list {
		for j, q := range s.QualityPrior {
			if t.Quality == q {
				ranks[i] = float32(len(s.QualityPrior)-j) * perQualityWeight
				break
			}
		}
		l.Debugf("%d rank by quality: %.4f", i, ranks[i])
	}

	return ranks
}
