package db

import (
	"errors"
	"path/filepath"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/dgraph-io/badger/v3"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	badgerstore "github.com/ostafen/clover/v2/store/badger"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

const (
	torrentsCollection = "torrents"
	metaCollection     = "meta"
	versionKey         = "Version"
)

var knownCollections = []string{
	torrentsCollection,
	metaCollection,
}

type metaInfo struct {
	Key   string
	Value string
}

type cloverDb struct {
	conn *clover.DB
}

func newCloverDB(cfg config.Storage) (Database, error) {
	dbPath := filepath.Join(cfg.Directory, "database")
	store, err := badgerstore.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}

	conn, err := clover.OpenWithStore(store)
	if err != nil {
		return nil, err
	}

	for _, collection := range knownCollections {
		exists, err := conn.HasCollection(collection)
		if err != nil {
			_ = conn.Close()
			return nil, err
		}

		if !exists {
			if err = conn.CreateCollection(collection); err != nil {
				_ = conn.Close()
				return nil, err
			}
		}
	}

	return &cloverDb{conn: conn}, nil
}

func (d *cloverDb) PutTorrent(t *model.Torrent) error {
	doc := document.NewDocumentOf(t)
	if doc == nil {
		return errors.New("deserialize document failed")
	}

	_, err := d.conn.InsertOne(torrentsCollection, doc)
	return err
}

func (d *cloverDb) RemoveTorrent(id string) error {
	return d.conn.Delete(query.NewQuery(torrentsCollection).Where(query.Field("ID").Eq(id)))
}

func (d *cloverDb) loadTorrents(t *media.ContentType) ([]*model.Torrent, error) {
	q := query.NewQuery(torrentsCollection)
	if t != nil {
		q = q.Where(query.Field("Type").Eq(*t))
	}
	docs, err := d.conn.FindAll(q)
	if err != nil {
		return []*model.Torrent{}, err
	}

	result := make([]*model.Torrent, len(docs))
	for i, doc := range docs {
		if err = doc.Unmarshal(&result[i]); err != nil {
			return []*model.Torrent{}, err
		}
	}

	return result, nil
}

// GetTorrent implements Database.
func (d *cloverDb) GetTorrent(id string) (*model.Torrent, error) {
	result := model.Torrent{}
	doc, err := d.conn.FindFirst(query.NewQuery(torrentsCollection).Where(query.Field("ID").Eq(id)))
	if err != nil {
		return &result, err
	}
	err = doc.Unmarshal(&result)
	return &result, err
}

func (d *cloverDb) LoadAllTorrents() ([]*model.Torrent, error) {
	return d.loadTorrents(nil)
}

func (d *cloverDb) LoadTorrents(mediaType media.ContentType) ([]*model.Torrent, error) {
	return d.loadTorrents(&mediaType)
}

func (d *cloverDb) Close() error {
	return d.conn.Close()
}

// GetVersion implements Database.
func (d *cloverDb) GetVersion() (string, error) {
	result := metaInfo{}
	doc, err := d.conn.FindFirst(query.NewQuery(metaCollection).Where(query.Field("Key").Eq(versionKey)))
	if err != nil {
		return "", err
	}
	err = doc.Unmarshal(&result)
	return result.Value, err
}

// SetVersion implements Database.
func (d *cloverDb) SetVersion(version string) error {
	doc := document.NewDocumentOf(&metaInfo{Key: versionKey, Value: version})
	if doc == nil {
		return errors.New("deserialize document failed")
	}
	d.conn.Delete(query.NewQuery(metaCollection).Where(query.Field("Key").Eq(versionKey)))
	_, err := d.conn.InsertOne(metaCollection, doc)
	return err
}
