package db

import (
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type Database interface {
	PutTorrent(t *model.Torrent) error
	RemoveTorrent(id string) error
	LoadTorrents(includeContent bool) ([]*model.Torrent, error)
	GetTorrent(id string) (*model.Torrent, error)
	GetVersion() (string, error)
	SetVersion(version string) error
	Close() error
}

type databaseInternal interface {
	Database
	GetDatabaseVersion() (uint, error)
	SetDatabaseVersion(version uint) error
}
