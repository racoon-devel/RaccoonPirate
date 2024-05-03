package web

import (
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type DiscoveryService interface {
	SearchMovies(ctx context.Context, q string) ([]*model.Movie, error)
	SearchTorrents(ctx context.Context, id string, season *int64) ([]*models.SearchTorrentsResult, error)
	GetSeasonsCount(id string) uint
	GetTorrent(ctx context.Context, link string) ([]byte, error)
}

type TorrentService interface {
	Add(data []byte) error
}
