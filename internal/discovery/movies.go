package discovery

import (
	"context"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client/movies"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func convertMovieInfo(in *models.SearchMoviesResult) *model.Movie {
	out := &model.Movie{
		ID:            *in.ID,
		Title:         *in.Title,
		OriginalTitle: in.OriginalTitle,
		Description:   in.Description,
		Year:          uint(in.Year),
		Poster:        in.Poster,
		Genres:        in.Genres,
		Rating:        float32(in.Rating),
		Seasons:       uint(in.Seasons),
	}

	if in.Type == "tv-series" {
		out.Type = model.MovieType_TvSeries
	} else {
		out.Type = model.MovieType_Movie
	}

	return out
}

func (s *Service) SearchMovies(ctx context.Context, q string) ([]*model.Movie, error) {
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
		mov := convertMovieInfo(r)
		result = append(result, mov)
	}

	return result, nil
}
