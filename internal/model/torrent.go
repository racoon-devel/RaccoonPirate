package model

type MediaType uint

const (
	MediaTypeMovie MediaType = iota
	MediaTypeTvSeries
	MediaTypeArtist
	MediaTypeOther
)

type Torrent struct {
	ID        string
	Title     string
	Type      MediaType
	BelongsTo string
	Year      uint
	Genres    []string
	Content   []byte
}
