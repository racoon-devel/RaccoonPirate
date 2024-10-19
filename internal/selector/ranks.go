package selector

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
)

func (s selection) rankWeight(weight float32, f rankFunc) rankFunc {
	return func(list []*models.SearchTorrentsResult) []float32 {
		ranks := f(list)
		for i := range ranks {
			ranks[i] = weight * ranks[i]
		}
		return ranks
	}
}

func (s selection) rankBySize(list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	_, max, _ := findMax(list, func(t *models.SearchTorrentsResult) int64 {
		return getValue(t.Size)
	})

	for i, t := range list {
		ranks[i] = 1 - (float32(getValue(t.Size)) / float32(max))
		s.log().Debugf("%d rank by size: %.4f", i, ranks[i])
	}
	return ranks
}

func (s selection) limitBySize(list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))

	for i, t := range list {
		size := getValue(t.Size)
		seasons := int64(len(t.Seasons))
		if seasons == 0 {
			seasons = 1
		}

		if size < s.MinSeasonSizeMB*seasons || size >= s.MaxSeasonSizeMB*seasons {
			ranks[i] = -1
			s.log().Debugf("%d limit by size: %.4f", i, ranks[i])
		}

	}
	return ranks
}

func (s selection) rankBySeeders(list []*models.SearchTorrentsResult) []float32 {
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
		s.log().Debugf("%d rank by seeders: %.4f", i, ranks[i])
	}
	return ranks
}

func (s selection) rankByQuality(list []*models.SearchTorrentsResult) []float32 {
	ranks := make([]float32, len(list))
	perQualityWeight := 1 / float32(len(s.QualityPrior))
	for i, t := range list {
		for j, q := range s.QualityPrior {
			if t.Quality == q {
				ranks[i] = float32(len(s.QualityPrior)-j) * perQualityWeight
				break
			}
		}
		s.log().Debugf("%d rank by quality: %.4f", i, ranks[i])
	}

	return ranks
}
