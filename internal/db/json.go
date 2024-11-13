package db

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type jsonDb struct {
	path string
}

type fileSchema struct {
	Torrents map[string]*model.Torrent
	Version  string
}

func newJsonDB(cfg config.Storage) (Database, error) {
	dbPath := filepath.Join(cfg.Directory, "database.db")
	db := jsonDb{path: dbPath}
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(cfg.Directory, 0755); err != nil {
			return nil, err
		}
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
		log.Errorf("Fail to write to '%s': %s", d.path, err)
		return err
	}
	return nil
}

// LoadAllTorrents implements Database.
func (d *jsonDb) LoadAllTorrents() ([]*model.Torrent, error) {
	content, err := d.load()
	result := make([]*model.Torrent, 0, len(content.Torrents))
	for _, t := range content.Torrents {
		result = append(result, t)
	}
	return result, err
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
	if content.Torrents == nil {
		content.Torrents = map[string]*model.Torrent{}
	}
	content.Torrents[t.ID] = t
	return d.save(content)
}

// GetTorrent implements Database.
func (d *jsonDb) GetTorrent(id string) (*model.Torrent, error) {
	content, err := d.load()
	if err != nil {
		return &model.Torrent{}, err
	}
	result, ok := content.Torrents[id]
	if !ok {
		return &model.Torrent{}, errors.New("not found")
	}
	return result, nil
}

// RemoveTorrent implements Database.
func (d *jsonDb) RemoveTorrent(id string) error {
	content, err := d.load()
	if err != nil {
		return err
	}
	delete(content.Torrents, id)
	return d.save(content)
}

// GetVersion implements Database.
func (d *jsonDb) GetVersion() (string, error) {
	content, err := d.load()
	if err != nil {
		return "", err
	}
	return content.Version, nil
}

// SetVersion implements Database.
func (d *jsonDb) SetVersion(version string) error {
	content, err := d.load()
	if err != nil {
		return err
	}
	content.Version = version
	return d.save(content)
}
