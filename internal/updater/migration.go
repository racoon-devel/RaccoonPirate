package updater

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/db"
)

type migrationRequest struct {
	major uint64
	minor uint64
	dbase db.Database
	cfg   config.Config
}

type migrateFn func(m *migrationRequest) error

var migrations = map[uint64]migrateFn{
	4: migrateV4ToV5,
}

func (u Updater) migrate(m *migrationRequest) error {
	currentVersion, err := ParseVersion(u.CurrentVersion)
	if err != nil {
		return fmt.Errorf("parse built-in version failed: %w", err)
	}

	if m.major != currentVersion.Major {
		return errors.New("migration between major versions is unsupported")
	}

	if m.minor >= currentVersion.Minor {
		return errors.New("migration from future versions is unsupported")
	}

	for version := m.minor; version < currentVersion.Minor; version++ {
		migrationFn, ok := migrations[version]
		if !ok {
			continue
		}
		if err = migrationFn(m); err != nil {
			return fmt.Errorf("migration from %d.%d to %d.%d failed: %w", currentVersion.Major, version, currentVersion.Major, version+1, err)
		}
	}

	log.Info("Migration was successful!")
	return u.Storage.SetVersion(u.CurrentVersion)
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func migrateV4ToV5(m *migrationRequest) error {
	torrentsDir := filepath.Join(filepath.Dir(m.cfg.Database.Path), "torrents")
	files, err := os.ReadDir(torrentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn("No data to migrate")
			return nil
		}
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(torrentsDir, f.Name()))
		if err != nil {
			log.Warnf("Cannot read '%s', skip: %s", f.Name(), err)
			continue
		}

		t, err := m.dbase.GetTorrent(fileNameWithoutExtension(f.Name()))
		if err != nil {
			log.Warnf("Torrent '%s' hasn't registered in database: %s", err)
			continue
		}

		mi, err := metainfo.Load(bytes.NewReader(data))
		if err != nil {
			log.Warnf("Parse torrent %s failed: %s", f.Name(), err)
			continue
		}

		t.Content = data
		oldId := t.ID
		t.ID = mi.HashInfoBytes().HexString()

		if err = m.dbase.PutTorrent(t); err != nil {
			log.Warnf("Re-add torrent %s failed: %s", f.Name(), err)
			continue
		}

		if oldId != t.ID {
			_ = m.dbase.RemoveTorrent(oldId)
		}
	}

	return nil
}
