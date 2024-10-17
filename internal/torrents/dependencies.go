package torrents

import "github.com/racoon-devel/raccoon-pirate/internal/model"

type Database interface {
	LoadTorrents() ([]*model.Torrent, error)
	PutTorrent(t *model.Torrent) error
	RemoveTorrent(id string) error
}
