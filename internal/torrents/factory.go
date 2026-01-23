package torrents

import (
	"fmt"
	"path/filepath"

	"github.com/RacoonMediaServer/rms-torrent/v4/pkg/engine"
	"github.com/RacoonMediaServer/rms-torrent/v4/pkg/engine/online/builtin"
	"github.com/RacoonMediaServer/rms-torrent/v4/pkg/engine/online/torrserver"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

func createTorrentEngine(cfg config.Torrent) (engine.TorrentEngine, error) {
	switch cfg.Driver {
	case "builtin":
		engineCfg := builtin.Config{
			Directory:   cfg.Builtin.Directory,
			Limit:       cfg.Builtin.Limit,
			AddTimeout:  cfg.Builtin.AddTimeout,
			ReadTimeout: cfg.Builtin.ReadTimeout,
			TTL:         cfg.Builtin.TTL,
		}
		torrentsDir := filepath.Join(engineCfg.Directory, "torrents")
		storage, err := newPersistentStorage(torrentsDir)
		if err != nil {
			return nil, fmt.Errorf("create directory for torrents failed: %w", err)
		}
		return builtin.NewEngine(engineCfg, storage)
	case "torr-server":
		engineCfg := torrserver.Config{
			URL:      cfg.TorrServer.URL,
			Location: cfg.TorrServer.Fusepath,
		}
		return torrserver.NewEngine(engineCfg)
	default:
		return nil, fmt.Errorf("unknown torrent engine id: %s", cfg.Driver)
	}
}
