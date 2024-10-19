package model

import (
	"encoding/json"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type Torrent struct {
	// ID is unique auto-generated torrent ID
	ID string

	// Title of original torrent file
	Title string

	// Type is suggested type of torrent's content
	Type media.ContentType

	// BelongsTo means relation to discovered film, tv-series or artist
	BelongsTo string

	// Year of production
	Year uint

	// Genres id a list of content genres (stored as json array)
	Genres string

	// MediaID is an ID of discovered item. Not used, but may be useful in the feature
	MediaID string
}

func (t *Torrent) SetGenres(list []string) {
	data, _ := json.Marshal(list)
	t.Genres = string(data)
}

func (t *Torrent) GetGenres() []string {
	result := []string{}
	_ = json.Unmarshal([]byte(t.Genres), &result)
	return result
}

func (t *Torrent) ExpandByMovie(mov *model.Movie) {
	t.Type = media.Movies
	t.BelongsTo = mov.Title
	t.Year = mov.Year
	t.MediaID = mov.ID
	t.SetGenres(mov.Genres)
}

func (t *Torrent) ExpandByMusic(m model.Music) {
	t.Type = media.Music
	if m.IsArtist() {
		t.BelongsTo = m.AsArtist().Name
	} else if m.IsAlbum() {
		t.BelongsTo = m.AsAlbum().Artist
	}
}
