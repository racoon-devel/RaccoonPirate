package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type jsonDb struct {
	path string
	bs   *byteStorage
}

type fileSchema struct {
	Torrents        map[string]*model.Torrent
	Version         string
	DatabaseVersion uint
}

func newJsonDB(cfg config.Database, bs *byteStorage) (databaseInternal, error) {
	db := jsonDb{path: cfg.Path, bs: bs}
	_, err := os.Stat(cfg.Path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
			return nil, err
		}
		if err = db.save(&fileSchema{DatabaseVersion: currentDatabaseVersion}); err != nil {
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
		bytes, err := d.bs.Load(t.ID, "torrent")
		if err == nil {
			t.Content = bytes
			result = append(result, t)
		} else {
			log.Warnf("Load torrent content for %s failed: %s, skip", err)
		}
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

	if err = d.bs.Add(t.ID, "torrent", t.Content); err != nil {
		return fmt.Errorf("save torrent file failed: %w", err)
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
	if result.Content, err = d.bs.Load(id, "torrent"); err != nil {
		return &model.Torrent{}, fmt.Errorf("load torrent file failed: %w", err)
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
	_ = d.bs.Del(id, "torrent")
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

// GetDatabaseVersion implements databaseInternal.
func (d *jsonDb) GetDatabaseVersion() (uint, error) {
	content, err := d.load()
	if err != nil {
		return 0, err
	}
	return content.DatabaseVersion, nil
}

// SetDatabaseVersion implements databaseInternal.
func (d *jsonDb) SetDatabaseVersion(version uint) error {
	content, err := d.load()
	if err != nil {
		return err
	}
	content.DatabaseVersion = version
	return d.save(content)
}
