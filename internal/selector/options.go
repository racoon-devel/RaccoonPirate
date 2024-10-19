package selector

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/apex/log"
)

type Options struct {
	Log         *log.Entry
	Criteria    Criteria
	MediaType   media.ContentType
	Query       string
	Discography bool
}

func (o *Options) WithLogger(l *log.Entry) *Options {
	o.Log = l
	return o
}

func (o *Options) WithCriteria(criteria Criteria) *Options {
	o.Criteria = criteria
	return o
}

func (o *Options) WithMediaType(mediaType media.ContentType) *Options {
	o.MediaType = mediaType
	return o
}

func (o *Options) WithQuery(q string) *Options {
	o.Query = q
	return o
}

func (o *Options) WithDiscrography(d bool) *Options {
	o.Discography = d
	return o
}
