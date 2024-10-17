package db

import (
	"errors"

	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

const torrentsCollection = "torrents"

func (d *Database) PutTorrent(t *model.Torrent) error {
	doc := document.NewDocumentOf(t)
	if doc == nil {
		return errors.New("deserialize document failed")
	}

	_, err := d.conn.InsertOne(torrentsCollection, doc)
	return err
}

func (d *Database) RemoveTorrent(id string) error {
	return d.conn.Delete(query.NewQuery(torrentsCollection).Where(query.Field("ID").Eq(id)))
}

func (d *Database) loadTorrents(t *model.MediaType) ([]*model.Torrent, error) {
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

func (d *Database) LoadAllTorrents() ([]*model.Torrent, error) {
	return d.loadTorrents(nil)
}

func (d *Database) LoadTorrents(mediaType model.MediaType) ([]*model.Torrent, error) {
	return d.loadTorrents(&mediaType)
}
