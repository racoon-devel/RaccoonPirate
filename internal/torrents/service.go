package torrents

import (
	"context"
	"fmt"
	"sync"

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

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func New(cfg config.Torrent, db Database, rep RepresentationService) (*Service, error) {
	e, err := createTorrentEngine(cfg, db)
	if err != nil {
		return nil, err
	}

	torrents, err := db.LoadTorrents(true)
	if err != nil {
		return nil, fmt.Errorf("load all torrents failed: %w", err)
	}

	srv := &Service{
		db:  db,
		l:   log.WithField("from", "torrent-service"),
		rep: rep,
		e:   e,
	}

	srv.ctx, srv.cancel = context.WithCancel(context.Background())
	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		srv.trySyncTorrents(torrents)
	}()

	return srv, nil
}

func (s *Service) Add(ctx context.Context, record *model.Torrent) error {
	torrentInfo, err := s.e.Add(context.Background(), determineCategory(record), record.Title, nil, record.Content)
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
	torrents, err := s.db.LoadTorrents(false)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Torrent, 0, len(torrents))
	for _, t := range torrents {
		if t.Type == mediaType {
			result = append(result, t)
		}
	}

	return result, nil
}

func (s *Service) Remove(ctx context.Context, id string) error {
	t, err := s.db.GetTorrent(id)
	if err != nil {
		return err
	}

	if err = s.db.RemoveTorrent(id); err != nil {
		return fmt.Errorf("remove torrent from db failed: %w", err)
	}

	if err = s.e.Remove(ctx, id); err != nil {
		s.l.Warnf("Remove torrent from engine failed: %s", err)
	}

	s.rep.Unregister(t)
	return nil
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	_ = s.e.Stop()
}
