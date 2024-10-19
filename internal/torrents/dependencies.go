package torrents

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type Database interface {
	LoadAllTorrents() ([]*model.Torrent, error)
	LoadTorrents(mediaType media.ContentType) ([]*model.Torrent, error)
	PutTorrent(t *model.Torrent) error
	RemoveTorrent(id string) error
}
