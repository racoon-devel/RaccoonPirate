package smartsearch

import (
	"context"

	"github.com/RacoonMediaServer/rms-library/pkg/movsearch"
	"github.com/RacoonMediaServer/rms-library/pkg/selector"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/go-openapi/runtime"

	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

type Service struct {
	auth runtime.ClientAuthInfoWriter
	cli  *client.Client
}

func NewService(clientFactory ClientFactory, conf config.Discovery) *Service {
	s := Service{}
	s.auth, s.cli = clientFactory.NewDiscoveryClient(conf.ApiPath)
	return &s
}

func convertMovieInfo(mov *model.Movie) *rms_library.MovieInfo {
	info := &rms_library.MovieInfo{
		Title:         mov.Title,
		Description:   mov.Description,
		Year:          uint32(mov.Year),
		Poster:        mov.Poster,
		Genres:        mov.Genres,
		Rating:        mov.Rating,
		Type:          rms_library.MovieType_TvSeries,
		OriginalTitle: mov.OriginalTitle,
	}

	if mov.Type == model.MovieType_Movie {
		info.Type = rms_library.MovieType_Film
	}

	if mov.Seasons != 0 {
		seasons := uint32(mov.Seasons)
		info.Seasons = &seasons
	}

	return info
}

func (s Service) SmartSearchMovieTorrents(ctx context.Context, mov *model.Movie, sel selector.MediaSelector, selopts selector.Options, season *int64) (torrents [][]byte, err error) {
	searchEngine := movsearch.NewRemoteSearchEngine(s.cli.Torrents, s.auth)

	var strategy movsearch.Strategy
	if mov.Type == model.MovieType_Movie {
		strategy = &movsearch.SimpleStrategy{Engine: searchEngine, Selector: sel}
	} else if season != nil {
		strategy = &movsearch.SeasonStrategy{Engine: searchEngine, Selector: sel, SeasonNo: uint(*season)}
	} else {
		strategy = &movsearch.FullStrategy{Engine: searchEngine, Selector: sel}
	}

	var results []movsearch.Result
	results, err = strategy.Search(ctx, mov.ID, convertMovieInfo(mov), selopts)
	if err != nil {
		return
	}

	for _, r := range results {
		torrents = append(torrents, r.Torrent)
	}

	return
}

func (s Service) SearchMovieTorrents(ctx context.Context, mov *model.Movie, season *int64) ([]*models.SearchTorrentsResult, error) {
	searchEngine := movsearch.NewRemoteSearchEngine(s.cli.Torrents, s.auth)
	var seasonNo *uint
	if season != nil {
		seasonNo = new(uint)
		*seasonNo = uint(*season)
	}
	return searchEngine.SearchTorrents(ctx, mov.ID, convertMovieInfo(mov), seasonNo)
}
