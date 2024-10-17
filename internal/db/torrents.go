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

func (d *Database) LoadTorrents() ([]*model.Torrent, error) {
	docs, err := d.conn.FindAll(query.NewQuery(torrentsCollection))
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
