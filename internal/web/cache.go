package web

import "github.com/RacoonMediaServer/rms-media-discovery/pkg/model"

func (s *Server) movieFromCache(id string) (*model.Movie, bool) {
	mov := &model.Movie{}
	val, ok := s.cache.Load(id)
	if ok {
		mov, ok = val.(*model.Movie)
	}

	return mov, ok
}
