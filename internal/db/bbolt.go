package db

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
)

type bboltDb struct {
	conn *bbolt.DB
}

var torrentsBucket = []byte("torrents")
var metainfoBucket = []byte("metainfo")
var filesBucket = []byte("files")

var dbVersionKey = []byte("dbVersion")
var versionKey = []byte("version")

// Close implements Database.
func (b *bboltDb) Close() error {
	return b.Close()
}

// GetDatabaseVersion implements databaseInternal.
func (b *bboltDb) GetDatabaseVersion() (version uint, err error) {
	err = b.conn.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(metainfoBucket)
		value := bucket.Get(dbVersionKey)
		if value == nil {
			return errors.New("dbVersion not found")
		}
		parsed, err := strconv.ParseUint(string(value), 10, 32)
		if err != nil {
			return err
		}
		version = uint(parsed)
		return nil
	})
	return
}

// GetTorrent implements Database.
func (b *bboltDb) GetTorrent(id string) (*model.Torrent, error) {
	panic("unimplemented")
}

// GetVersion implements Database.
func (b *bboltDb) GetVersion() (version string, err error) {
	err = b.conn.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(metainfoBucket)
		value := bucket.Get(versionKey)
		if value == nil {
			return errors.New("version not found")
		}
		version = string(value)
		return nil
	})
	return
}

// LoadAllTorrents implements Database.
func (b *bboltDb) LoadAllTorrents() ([]*model.Torrent, error) {
	panic("unimplemented")
}

// LoadTorrents implements Database.
func (b *bboltDb) LoadTorrents(mediaType media.ContentType) ([]*model.Torrent, error) {
	panic("unimplemented")
}

// PutTorrent implements Database.
func (b *bboltDb) PutTorrent(t *model.Torrent) error {
	data, err := msgpack.Marshal(t)
	if err != nil {
		return fmt.Errorf("serialize torrent data failed: %w", err)
	}

	return b.conn.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(torrentsBucket)
		return bucket.Put([]byte(t.ID), data)
	})
}

// RemoveTorrent implements Database.
func (b *bboltDb) RemoveTorrent(id string) error {
	return b.conn.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(torrentsBucket)
		return bucket.Delete([]byte(id))
	})
}

// SetDatabaseVersion implements databaseInternal.
func (b *bboltDb) SetDatabaseVersion(version uint) error {
	return b.conn.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(metainfoBucket)
		return bucket.Put(dbVersionKey, []byte(strconv.Itoa(int(version))))
	})
}

// SetVersion implements Database.
func (b *bboltDb) SetVersion(version string) error {
	return b.conn.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(metainfoBucket)
		return bucket.Put(versionKey, []byte(version))
	})
}

func newBoltDB(cfg config.Database) (databaseInternal, error) {
	db, err := bbolt.Open(cfg.Path, 0644, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(torrentsBucket); err != nil {
			return err
		}
		mi, err := tx.CreateBucketIfNotExists(metainfoBucket)
		if err != nil {
			return err
		}

		if cv := mi.Get(dbVersionKey); cv == nil {
			return mi.Put(dbVersionKey, []byte(strconv.Itoa(currentDatabaseVersion)))
		}
		return nil
	})

	return &bboltDb{conn: db}, err
}
