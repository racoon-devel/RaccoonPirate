package web

import (
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/gin-gonic/gin"
)

type searchPage struct {
	uiPage
	Query     string
	MediaType string
	Movies    []*model.Movie
}

func (s *Server) searchHandler(ctx *gin.Context) {
	q := ctx.Query("q")
	mediaType := ctx.Query("media-type")
	if mediaType == "" {
		mediaType = "movies"
	}
	page := searchPage{
		Query:     q,
		MediaType: mediaType,
	}

	if q != "" {
		var err error
		l := s.l.WithField("query", q).WithField("media-type", mediaType)
		l.Debugf("Search")
		page.Movies, err = s.DiscoveryService.SearchMovies(ctx, q)
		if err != nil {
			l.Errorf("Search failed: %s", err)
			displayError(ctx, http.StatusInternalServerError, "Something went wrong...")
			return
		}
		if len(page.Movies) == 0 {
			l.Warnf("Nothing found")
			displayError(ctx, http.StatusBadRequest, "No results found")
			return
		}
		for _, mov := range page.Movies {
			s.cache.Store(mov.ID, mov)
		}
	}
	ctx.HTML(http.StatusOK, "multimedia.search.tmpl", &page)
}
