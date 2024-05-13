package selector

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	list     []*models.SearchTorrentsResult
	result   int
	criteria Criteria
	voice    string
}

func makeValuePtr[T any](val T) *T {
	p := new(T)
	*p = val
	return p
}

var testCases = []testCase{
	{
		list: []*models.SearchTorrentsResult{
			{
				Seeders: makeValuePtr[int64](2),
				Size:    makeValuePtr[int64](3225),
				Voice:   "Сыендук",
			},
			{
				Seeders: makeValuePtr[int64](16),
				Size:    makeValuePtr[int64](42434),
			},
			{
				Seeders: makeValuePtr[int64](950),
				Size:    makeValuePtr[int64](14919),
			},
			{
				Seeders: makeValuePtr[int64](387),
				Size:    makeValuePtr[int64](8724),
			},
		},
		criteria: CriteriaFastest,
		result:   3,
	},
	{
		list: []*models.SearchTorrentsResult{
			{
				Seeders: makeValuePtr[int64](2),
				Size:    makeValuePtr[int64](3225),
				Voice:   "Сыендук",
			},
			{
				Seeders: makeValuePtr[int64](16),
				Size:    makeValuePtr[int64](42434),
				Quality: "1080p",
			},
			{
				Seeders: makeValuePtr[int64](950),
				Size:    makeValuePtr[int64](14919),
				Quality: "720p",
				Voice:   "Сыендук",
			},
			{
				Seeders: makeValuePtr[int64](387),
				Size:    makeValuePtr[int64](8724),
				Quality: "480p",
				Voice:   "Vo сыендук",
			},
		},
		criteria: CriteriaQuality,
		result:   2,
	},
	{
		list: []*models.SearchTorrentsResult{
			{
				Seeders: makeValuePtr[int64](2),
				Size:    makeValuePtr[int64](3225),
				Voice:   "Сыендук",
			},
			{
				Seeders: makeValuePtr[int64](16),
				Size:    makeValuePtr[int64](42434),
				Quality: "1080p",
				Seasons: []int64{1, 2, 3, 4},
			},
			{
				Seeders: makeValuePtr[int64](950),
				Size:    makeValuePtr[int64](14919),
				Quality: "720p",
				Voice:   "Сыендук",
				Seasons: []int64{1},
			},
			{
				Seeders: makeValuePtr[int64](387),
				Size:    makeValuePtr[int64](8724),
				Quality: "480p",
				Voice:   "Vo сыендук",
			},
		},
		criteria: CriteriaCompact,
		result:   2,
	},
	{
		list: []*models.SearchTorrentsResult{
			{
				Seeders: makeValuePtr[int64](2),
				Size:    makeValuePtr[int64](3225),
				Voice:   "vo Сыендук Original",
			},
			{
				Seeders: makeValuePtr[int64](16),
				Size:    makeValuePtr[int64](42434),
			},
			{
				Seeders: makeValuePtr[int64](950),
				Size:    makeValuePtr[int64](14919),
				Voice:   "Mvo Tvshows Original",
			},
			{
				Seeders: makeValuePtr[int64](387),
				Size:    makeValuePtr[int64](8724),
			},
		},
		voice:    "Vo Сыендук | TVShows Studio",
		criteria: CriteriaQuality,
		result:   0,
	},
}

func TestMovieSelector_Select(t *testing.T) {
	s := MovieSelector{
		MinSeasonSizeMB:     1024,
		MaxSeasonSizeMB:     1024 * 50,
		MinSeedersThreshold: 50,
		Voice:               "",
		QualityPrior: []string{
			"1080p",
			"720p",
			"480p",
		},
	}

	s.VoiceList.Append("сыендук")

	log.SetLevel(log.DebugLevel)

	for i, test := range testCases {
		log.Debugf("\n\n############### Test %d ###############", i)
		s.Voice = test.voice
		result := s.Select(log.WithField("from", "selector"), test.criteria, test.list)
		assert.Equal(t, test.list[test.result], result, "test case %d failed", i)
	}
}
