package torrents

import (
	"context"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-torrent/v4/pkg/engine"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type Service struct {
	l   *log.Entry
	db  Database
	rep RepresentationService
	e   engine.TorrentEngine
}

func New(cfg config.Torrent, db Database, rep RepresentationService) (*Service, error) {
	e, err := createTorrentEngine(cfg)
	if err != nil {
		return nil, err
	}

	// TODO: mount exists torrents to the representation

	return &Service{
		db:  db,
		l:   log.WithField("from", "torrent-service"),
		rep: rep,
		e:   e,
	}, nil
}

func (s *Service) Add(ctx context.Context, record *model.Torrent, content []byte) error {
	torrentInfo, err := s.e.Add(context.Background(), determineCategory(record), record.Title, nil, content)
	if err != nil {
		return err
	}

	record.ID = torrentInfo.ID
	record.Title = torrentInfo.Title

	s.rep.Register(record, torrentInfo.Location)

	if err = s.db.PutTorrent(record); err != nil {
		s.l.Errorf("Store info about torrent failed: %s", err)
	}

	return nil
}

func (s *Service) GetTorrentsList(mediaType media.ContentType) ([]*model.Torrent, error) {
	return s.db.LoadTorrents(mediaType)
}

func (s *Service) Remove(ctx context.Context, id string) error {
	// TODO: make consistent
	t, err := s.db.GetTorrent(id)
	if err != nil {
		return err
	}

	if err := s.e.Remove(ctx, id); err != nil {
		return err
	}

	s.rep.Unregister(t)
	return s.db.RemoveTorrent(id)
}

func (s *Service) GetContentDirectory() string {
	return ""
	// TODO: find the right way for it
	// return filepath.Join(s.layout.contentDir, mainRoute)
}

func (s *Service) Stop() {
	_ = s.e.Stop()
}
