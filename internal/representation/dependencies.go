package representation

import "github.com/racoon-devel/raccoon-pirate/internal/model"

type Storage interface {
	LoadAllTorrents() ([]*model.Torrent, error)
}
