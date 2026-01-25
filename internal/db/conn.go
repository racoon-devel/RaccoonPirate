package db

import (
	"fmt"
	"path/filepath"

	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

const currentDatabaseVersion = 2
const byteStorageDir = "torrents"

func Open(cfg config.Database) (Database, error) {
	storageDir := filepath.Join(filepath.Dir(cfg.Path), byteStorageDir)

	var dbase databaseInternal
	var err error
	switch cfg.Driver {
	case "cloverdb":
		dbase, err = newBoltDB(cfg)
	case "json":
		bs, err := newByteStorage(storageDir)
		if err != nil {
			return nil, fmt.Errorf("create byte storage directory failed: %w", err)
		}
		dbase, err = newJsonDB(cfg, bs)
	default:
		return nil, fmt.Errorf("database driver '%s' is not implemented", cfg.Driver)
	}

	if err != nil {
		return nil, err
	}

	version, err := dbase.GetDatabaseVersion()
	if err != nil {
		return nil, fmt.Errorf("get database version failed: %w", err)
	}

	if version == 0 {
		version = 1
	}

	if version != currentDatabaseVersion {
		log.Warnf("database version %d != %d, trying to migrate...", version, currentDatabaseVersion)
		if err = migrateDatabase(dbase, version); err != nil {
			return nil, fmt.Errorf("migration failed: %w", err)
		}
	}

	return dbase, nil
}
