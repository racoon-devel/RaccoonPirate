package db

import (
	"path/filepath"

	"github.com/dgraph-io/badger/v3"
	"github.com/ostafen/clover/v2"
	badgerstore "github.com/ostafen/clover/v2/store/badger"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

type Database struct {
	conn *clover.DB
}

func Open(cfg config.Storage) (*Database, error) {
	dbPath := filepath.Join(cfg.Directory, "database")
	store, err := badgerstore.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return &Database{}, err
	}

	conn, err := clover.OpenWithStore(store)
	if err != nil {
		return &Database{}, err
	}

	exists, err := conn.HasCollection(torrentsCollection)
	if err != nil {
		_ = conn.Close()
		return &Database{}, err
	}

	if !exists {
		if err = conn.CreateCollection(torrentsCollection); err != nil {
			_ = conn.Close()
			return &Database{}, err
		}
	}

	return &Database{conn: conn}, nil
}

func (d *Database) Close() error {
	return d.conn.Close()
}
