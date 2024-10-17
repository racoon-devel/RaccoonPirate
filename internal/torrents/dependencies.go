package torrents

import "github.com/racoon-devel/raccoon-pirate/internal/model"

type Database interface {
	LoadAllTorrents() ([]*model.Torrent, error)
	LoadTorrents(mediaType model.MediaType) ([]*model.Torrent, error)
	PutTorrent(t *model.Torrent) error
	RemoveTorrent(id string) error
}
