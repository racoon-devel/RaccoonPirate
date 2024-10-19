package web

import (
	"context"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	internalModel "github.com/racoon-devel/raccoon-pirate/internal/model"
)

type DiscoveryService interface {
	SearchMovies(ctx context.Context, q string) ([]*model.Movie, error)
	SearchMusic(ctx context.Context, q string) ([]model.Music, error)
	SearchMovieTorrents(ctx context.Context, mov *model.Movie, season *int64) ([]*models.SearchTorrentsResult, error)
	SearchMusicTorrents(ctx context.Context, m model.Music) ([]*models.SearchTorrentsResult, error)
	SearchOtherTorrents(ctx context.Context, q string) ([]*models.SearchTorrentsResult, error)
	GetTorrent(ctx context.Context, link string) ([]byte, error)
}

type TorrentService interface {
	Add(record *internalModel.Torrent, data []byte) error
	GetTorrentsList(mediaType media.ContentType) ([]*internalModel.Torrent, error)
	Remove(torrent string) error
}
