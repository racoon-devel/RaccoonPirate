package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchPage struct {
	uiPage
	Query string
}

func (s *Server) searchHandler(ctx *gin.Context) {
	page := searchPage{
		Query: ctx.Query("q"),
	}
	ctx.HTML(http.StatusOK, "multimedia.search.tmpl", &page)
}
