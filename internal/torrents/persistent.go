package torrents

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RacoonMediaServer/rms-torrent/v4/pkg/engine"
)

type persistentStorage struct {
	torrentsDir string
}

func newPersistentStorage(torrentsDir string) (engine.TorrentDatabase, error) {
	return &persistentStorage{torrentsDir: torrentsDir}, os.MkdirAll(torrentsDir, 0744)
}

// Add implements engine.TorrentDatabase.
func (p persistentStorage) Add(t engine.TorrentRecord) error {
	filePath := filepath.Join(p.torrentsDir, fmt.Sprintf("%s.torrent", t.ID))
	return os.WriteFile(filePath, t.Content, 0744)
}

// Complete implements engine.TorrentDatabase.
func (p persistentStorage) Complete(id string) error {
	return errors.ErrUnsupported
}

// Del implements engine.TorrentDatabase.
func (p persistentStorage) Del(id string) error {
	filePath := filepath.Join(p.torrentsDir, fmt.Sprintf("%s.torrent", id))
	return os.Remove(filePath)
}

// Load implements engine.TorrentDatabase.
func (p persistentStorage) Load() ([]engine.TorrentRecord, error) {
	result := []engine.TorrentRecord{}
	files, err := os.ReadDir(p.torrentsDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			data, err := os.ReadFile(filepath.Join(p.torrentsDir, f.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, engine.TorrentRecord{Content: data})
		}
	}
	return result, nil
}
