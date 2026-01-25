package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/apex/log"
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
	return b.conn.Close()
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
func (b *bboltDb) GetTorrent(id string) (t *model.Torrent, err error) {
	t = &model.Torrent{}

	err = b.conn.View(func(tx *bbolt.Tx) error {
		torrents := tx.Bucket(torrentsBucket)
		files := tx.Bucket(filesBucket)

		rawData := torrents.Get([]byte(id))
		if rawData == nil {
			return ErrNotFound
		}
		if err := msgpack.Unmarshal(rawData, t); err != nil {
			return fmt.Errorf("deserialize torrent info failed: %+w", err)
		}

		rawFile := files.Get([]byte(id))
		if rawFile == nil {
			return fmt.Errorf("cannot load torrent file: %w", ErrNotFound)
		}
		t.Content = rawFile

		return nil
	})
	return
}

// GetVersion implements Database.
func (b *bboltDb) GetVersion() (version string, err error) {
	err = b.conn.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(metainfoBucket)
		value := bucket.Get(versionKey)
		if value == nil {
			return nil
		}
		version = string(value)
		return nil
	})
	return
}

// LoadAllTorrents implements Database.
func (b *bboltDb) LoadTorrents(includeContent bool) (result []*model.Torrent, err error) {
	err = b.conn.View(func(tx *bbolt.Tx) error {
		torrents := tx.Bucket(torrentsBucket)
		files := tx.Bucket(filesBucket)
		c := torrents.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := model.Torrent{}
			id := string(k)
			if err := msgpack.Unmarshal(v, &t); err != nil {
				log.Warnf("Deserialize %s torrent info failed: %s", id, err)
				continue
			}

			if includeContent {
				t.Content = files.Get(k)
			}

			result = append(result, &t)
		}

		return nil
	})

	return
}

// PutTorrent implements Database.
func (b *bboltDb) PutTorrent(t *model.Torrent) error {
	data, err := msgpack.Marshal(t)
	if err != nil {
		return fmt.Errorf("serialize torrent data failed: %w", err)
	}

	return b.conn.Update(func(tx *bbolt.Tx) error {
		files := tx.Bucket(filesBucket)
		if err := files.Put([]byte(t.ID), t.Content); err != nil {
			return err
		}

		torrents := tx.Bucket(torrentsBucket)
		return torrents.Put([]byte(t.ID), data)
	})
}

// RemoveTorrent implements Database.
func (b *bboltDb) RemoveTorrent(id string) error {
	return b.conn.Update(func(tx *bbolt.Tx) error {
		torrents := tx.Bucket(torrentsBucket)
		if err := torrents.Delete([]byte(id)); err != nil {
			return err
		}
		files := tx.Bucket(filesBucket)
		return files.Delete([]byte(id))
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
	_, err := os.Stat(cfg.Path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
			return nil, err
		}
	}

	db, err := bbolt.Open(cfg.Path, 0644, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(torrentsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(filesBucket); err != nil {
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
