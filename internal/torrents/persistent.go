package torrents

import (
	"github.com/RacoonMediaServer/rms-torrent/v4/pkg/engine"
)

type persistentStorage struct {
	dbase Database
}

// Add implements engine.TorrentDatabase.
func (p persistentStorage) Add(t engine.TorrentRecord) error {
	// content has already stored in the database
	return nil
}

// Complete implements engine.TorrentDatabase.
func (p persistentStorage) Complete(id string) error {
	// isn't actual for online torrent engines
	return nil
}

// Del implements engine.TorrentDatabase.
func (p persistentStorage) Del(id string) error {
	// content managed by Database
	return nil
}

// Load implements engine.TorrentDatabase.
func (p persistentStorage) Load() ([]engine.TorrentRecord, error) {
	torrents, err := p.dbase.LoadAllTorrents()
	if err != nil {
		return nil, err
	}

	result := make([]engine.TorrentRecord, 0, len(torrents))
	for _, t := range torrents {
		record := engine.TorrentRecord{
			TorrentDescription: engine.TorrentDescription{
				ID:    t.ID,
				Title: t.Title,
			},
			Content: t.Content,
		}
		result = append(result, record)
	}

	return result, nil
}
