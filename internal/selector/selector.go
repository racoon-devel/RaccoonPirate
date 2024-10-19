package selector

import (
	"sort"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
)

type MediaSelector struct {
	settings Settings
}

func New(settings Settings) MediaSelector {
	return MediaSelector{settings: settings}
}

func (s MediaSelector) Select(list []*models.SearchTorrentsResult, opts Options) *models.SearchTorrentsResult {
	selCtx := selection{
		Settings: s.settings,
		Options:  opts,
	}
	rank := selCtx.getRankFunction()
	ranks := rank(list)
	_, _, best := findMax(ranks, func(elem float32) float32 {
		return elem
	})
	if opts.Log != nil {
		for i := range ranks {
			opts.Log.Debugf("%d rank: %.4f", i, ranks[i])
		}
	}
	sel := list[best]
	if opts.Log != nil {
		opts.Log.Infof("Selected { Title: %s, Voice: %s, Size: %d, Seeders: %d, Quality: %s }", getString(sel.Title), sel.Voice, getValue(sel.Size), getValue(sel.Seeders), sel.Quality)
	}
	return sel
}

func (s MediaSelector) Sort(list []*models.SearchTorrentsResult, opts Options) {
	selCtx := selection{
		Settings: s.settings,
		Options:  opts,
	}
	rank := selCtx.getRankFunction()
	ranks := rank(list)
	sort.SliceStable(list, func(i, j int) bool {
		return ranks[j] < ranks[i]
	})
}
