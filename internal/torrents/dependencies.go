package torrents

import (
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type Database interface {
	LoadTorrents(includeContent bool) ([]*model.Torrent, error)
	GetTorrent(id string) (*model.Torrent, error)
	PutTorrent(t *model.Torrent) error
	RemoveTorrent(id string) error
}

type RepresentationService interface {
	Register(t *model.Torrent, location string)
	Unregister(t *model.Torrent)
}
