package selector

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/apex/log"
)

type selection struct {
	Settings
	Options
}

func (s selection) getRankFunction() rankFunc {
	switch s.MediaType {
	case media.Movies:
		return s.getMovieRankFunction()
	case media.Music:
		return s.getMusicRankFunction()
	default:
		return s.getOtherRankFunction()
	}
}

func (s selection) log() *log.Entry {
	if s.Log != nil {
		return s.Log
	}

	l := log.WithField("from", "selector")
	l.Level = log.InvalidLevel
	return l
}
