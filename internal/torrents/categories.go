package torrents

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	mediaModel "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

func determineCategory(t *model.Torrent) string {
	// map to TorrServer categories
	switch t.Type {
	case media.Movies:
		if t.MovieType == mediaModel.MovieType_TvSeries {
			return "tv"
		}
		return "movie"
	case media.Music:
		return "music"
	default:
		return "other"
	}
}
