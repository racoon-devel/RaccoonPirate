package db

import (
	"errors"
	"fmt"

	"github.com/apex/log"
)

type migrateFn func(databaseInternal) error

var migrations = map[uint]migrateFn{
	1: migrateV1toV2,
}

func migrateDatabase(dbase databaseInternal, version uint) error {
	if version > currentDatabaseVersion {
		return fmt.Errorf("database created by the newest version: %d", version)
	}
	if version < 1 {
		return fmt.Errorf("database created by the unknown version: %d", version)
	}

	return errors.ErrUnsupported

	for nextVersion := version; version < currentDatabaseVersion; version++ {
		if err := migrations[nextVersion](dbase); err != nil {
			return fmt.Errorf("migrate from %d to %d failed: %w", nextVersion, nextVersion+1, err)
		}
	}

	return nil
}

func migrateV1toV2(dbase databaseInternal) error {
	torrents, err := dbase.LoadAllTorrents()
	if err != nil {
		return fmt.Errorf("load all torrents failed: %w", err)
	}

	for _, t := range torrents {
		// TODO
		if err = dbase.RemoveTorrent(t.ID); err != nil {
			log.Warnf("Remove torrent %s failed: %s", t.ID, err)
		}
		if err = dbase.PutTorrent(t); err != nil {
			log.Warnf("Re-add torrent %s failed: %s", t.ID, err)
		}
	}

	return dbase.SetDatabaseVersion(2)
}
