package db

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type jsonDb struct {
	path string
}

type fileSchema struct {
	Torrents []*model.Torrent
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func newJsonDB(cfg config.Storage) (Database, error) {
	dbPath := filepath.Join(cfg.Directory, "database.db")
	db := jsonDb{path: dbPath}
	_, err := os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		if err = db.save(&fileSchema{}); err != nil {
			return nil, err
		}
	}
	return &db, nil
}

// Close implements Database.
func (d *jsonDb) Close() error {
	return nil
}

func (d *jsonDb) load() (*fileSchema, error) {
	content, err := os.ReadFile(d.path)
	if err != nil {
		return &fileSchema{}, err
	}

	result := fileSchema{}
	if err = json.Unmarshal(content, &result); err != nil {
		return &fileSchema{}, err
	}

	return &result, nil
}

func (d *jsonDb) save(content *fileSchema) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if err = os.WriteFile(d.path, data, 0755); err != nil {
		return err
	}
	return nil
}

// LoadAllTorrents implements Database.
func (d *jsonDb) LoadAllTorrents() ([]*model.Torrent, error) {
	result, err := d.load()
	return result.Torrents, err
}

// LoadTorrents implements Database.
func (d *jsonDb) LoadTorrents(mediaType media.ContentType) ([]*model.Torrent, error) {
	list, err := d.LoadAllTorrents()
	if err != nil {
		return list, err
	}

	result := make([]*model.Torrent, 0, len(list))
	for _, t := range list {
		if t.Type == mediaType {
			result = append(result, t)
		}
	}
	return result, nil
}

// PutTorrent implements Database.
func (d *jsonDb) PutTorrent(t *model.Torrent) error {
	content, err := d.load()
	if err != nil {
		return err
	}
	content.Torrents = append(content.Torrents, t)
	return d.save(content)
}

// RemoveTorrent implements Database.
func (d *jsonDb) RemoveTorrent(id string) error {
	content, err := d.load()
	if err != nil {
		return err
	}
	found := false
	for i, t := range content.Torrents {
		if t.ID == id {
			content.Torrents = remove(content.Torrents, i)
			found = true
			break
		}
	}
	if !found {
		return errors.New("not found")
	}
	return d.save(content)
}
