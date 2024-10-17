package model

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
