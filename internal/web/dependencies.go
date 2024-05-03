package web

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type DiscoveryService interface {
	SearchMovies(ctx context.Context, q string) ([]*model.Movie, error)
}
