package torrents

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tConfig "github.com/RacoonMediaServer/distribyted/config"
	"github.com/RacoonMediaServer/distribyted/fuse"
	"github.com/RacoonMediaServer/distribyted/torrent"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/anacrolix/missinggo/v2/filecache"
	aTorrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
	uuid "github.com/satori/go.uuid"
)

type Service struct {
	l      *log.Entry
	layout layout
	db     Database
	rep    RepresentationService

	fuse      *fuse.Handler
	fileStore *torrent.FileItemStore
	cli       *aTorrent.Client
	service   *torrent.Service
}

func New(cfg config.Storage, db Database, rep RepresentationService) (*Service, error) {
	absPath, err := filepath.Abs(cfg.Directory)
	if err != nil {
		return nil, err
	}

	torrents, err := db.LoadAllTorrents()
	if err != nil {
		return nil, err
	}

	s := Service{
		layout: newLayout(absPath),
		db:     db,
		l:      log.WithField("from", "torrent-service"),
		rep:    rep,
	}
	if err := s.layout.makeLayout(); err != nil {
		return nil, err
	}

	fileCache, err := filecache.NewCache(s.layout.cacheDir)
	if err != nil {
		return nil, fmt.Errorf("create cache failed: %w", err)
	}
	fileCache.SetCapacity(int64(cfg.Limit) * 1024 * 1024 * 1024)

	torrentStorage := storage.NewResourcePieces(fileCache.AsResourceProvider())

	fileStorage, err := torrent.NewFileItemStore(s.layout.itemsDir, time.Duration(cfg.TTL)*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("create file store failed: %w", err)
	}

	id, err := torrent.GetOrCreatePeerID(filepath.Join(s.layout.baseDir, "ID"))
	if err != nil {
		return nil, fmt.Errorf("create ID failed: %w", err)
	}

	conf := tConfig.TorrentGlobal{
		ReadTimeout:     int(cfg.ReadTimeout),
		AddTimeout:      int(cfg.AddTimeout),
		GlobalCacheSize: -1,
		MetadataFolder:  s.layout.baseDir,
	}

	cli, err := torrent.NewClient(torrentStorage, fileStorage, &conf, id)
	if err != nil {
		return nil, fmt.Errorf("start torrent client failed: %w", err)
	}

	stats := torrent.NewStats()

	loaders := []torrent.DatabaseLoader{&s.layout}
	service := torrent.NewService(loaders, stats, cli, conf.AddTimeout, conf.ReadTimeout, true)

	fss, err := service.Load()
	if err != nil {
		return nil, fmt.Errorf("load torrents failed: %w", err)
	}

	mh := fuse.NewHandler(true, s.layout.contentDir)
	if err = mh.Mount(fss); err != nil {
		return nil, fmt.Errorf("mount fuse directory: %w", err)
	}

	for _, t := range torrents {
		rep.Register(t, filepath.Join(s.layout.contentDir, mainRoute, t.Title))
	}

	s.fuse = mh
	s.fileStore = fileStorage
	s.cli = cli
	s.service = service

	return &s, nil
}

func (s *Service) Add(record *model.Torrent, content []byte) error {
	title, _, err := s.service.Add(mainRoute, content)
	if err != nil {
		return err
	}

	record.ID = uuid.NewV4().String()
	record.Title = title

	s.rep.Register(record, filepath.Join(s.layout.contentDir, mainRoute, title))

	if err = os.WriteFile(filepath.Join(s.layout.torrentsDir, fmt.Sprintf("%s.torrent", record.ID)), content, 0744); err != nil {
		s.l.Errorf("Create torrent faile failed: %s", err)
		return nil
	}

	if err = s.db.PutTorrent(record); err != nil {
		s.l.Errorf("Store info about torrent failed: %s", err)
	}

	return nil
}

func (s *Service) GetTorrentsList(mediaType media.ContentType) ([]*model.Torrent, error) {
	return s.db.LoadTorrents(mediaType)
}

func (s *Service) Remove(id string) error {
	if err := s.layout.Remove(id); err != nil {
		return err
	}

	t, err := s.db.GetTorrent(id)
	if err != nil {
		return err
	}

	s.rep.Unregister(t)
	return s.db.RemoveTorrent(id)
}

func (s *Service) GetContentDirectory() string {
	return filepath.Join(s.layout.contentDir, mainRoute)
}

func (s *Service) Stop() {
	_ = s.fileStore.Close()
	s.cli.Close()
	s.fuse.Unmount()
}
