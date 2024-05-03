package web

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchPage struct {
	uiPage
	Query  string
	Movies []*model.Movie
}

func (s *Server) searchHandler(ctx *gin.Context) {
	q := ctx.Query("q")
	page := searchPage{
		Query: q,
	}
	if q != "" {
		var err error
		l := s.l.WithField("query", q)
		l.Debugf("Search")
		page.Movies, err = s.DiscoveryService.SearchMovies(ctx, q)
		if err != nil {
			l.Errorf("Search movies failed: %s", err)
			displayError(ctx, http.StatusInternalServerError, "Something went wrong...")
			return
		}
		if len(page.Movies) == 0 {
			l.Warnf("Nothing found")
			displayError(ctx, http.StatusBadRequest, "No results found")
			return
		}
	}
	ctx.HTML(http.StatusOK, "multimedia.search.tmpl", &page)
}
