package discovery

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/racoon-devel/media-station/internal/config"
)

const searchResultsLimit = 8

type Service struct {
	auth runtime.ClientAuthInfoWriter
	cli  *client.Client
}

func convertMovieInfo(in *models.SearchMoviesResult) *model.Movie {
	out := &model.Movie{
		ID:          *in.ID,
		Title:       *in.Title,
		Description: in.Description,
		Year:        uint(in.Year),
		Poster:      in.Poster,
		Genres:      in.Genres,
		Rating:      float32(in.Rating),
		Seasons:     uint(in.Seasons),
	}

	if in.Type == "tv-series" {
		out.Type = model.MovieType_TvSeries
	} else {
		out.Type = model.MovieType_Movie
	}

	return out
}

func NewService(conf config.Discovery) *Service {
	tr := httptransport.New(conf.Host, conf.Path, []string{conf.Scheme})
	auth := httptransport.APIKeyAuth("X-Token", "header", conf.Identity)
	cli := client.New(tr, strfmt.Default)

	return &Service{
		auth: auth,
		cli:  cli,
	}
}

func (s Service) SearchMovies(ctx context.Context, q string) ([]*model.Movie, error) {
	limit := int64(searchResultsLimit)
	req := &movies.SearchMoviesParams{
		Limit:   &limit,
		Q:       q,
		Context: ctx,
	}
	resp, err := s.cli.Movies.SearchMovies(req, s.auth)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Movie, 0, len(resp.Payload.Results))
	for _, r := range resp.Payload.Results {
		result = append(result, convertMovieInfo(r))
	}

	return result, nil
}
