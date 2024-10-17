package model

import (
	"encoding/json"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

type MediaType uint

const (
	MediaTypeMovie MediaType = iota
	MediaTypeTvSeries
	MediaTypeArtist
	MediaTypeOther
)

type Torrent struct {
	// ID is unique auto-generated torrent ID
	ID string

	// Title of original torrent file
	Title string

	// Type is suggested type of torrent's content
	Type MediaType

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
	if mov.Type == model.MovieType_Movie {
		t.Type = MediaTypeMovie
	} else {
		t.Type = MediaTypeTvSeries
	}

	t.BelongsTo = mov.Title
	t.Year = mov.Year
	t.MediaID = mov.ID
	t.SetGenres(mov.Genres)
}
