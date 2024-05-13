package web

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"sync"
)

type Server struct {
	l                *log.Entry
	g                *gin.Engine
	cache            sync.Map
	DiscoveryService DiscoveryService
	TorrentService   TorrentService
}

func (s *Server) Run(host string, port uint16) error {
	s.l = log.WithField("from", "web")
	s.g = gin.Default()

	root := template.New("root")
	templates := template.Must(root.ParseFS(templatesFS, "templates/*.tmpl"))
	s.g.SetHTMLTemplate(templates)
	s.g.StaticFS("/css", http.FS(wrapFS(webFS, "content/css")))
	s.g.StaticFS("/img", http.FS(wrapFS(webFS, "content/img")))
	s.g.StaticFS("/js", http.FS(wrapFS(webFS, "content/js")))

	s.g.NoRoute(func(ctx *gin.Context) {
		displayError(ctx, http.StatusNotFound, "Page not found")
	})

	s.g.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/search")
	})
	s.g.GET("/search", s.searchHandler)
	s.g.GET("/add/:id", s.addHandler)

	s.g.GET("/upload", s.getUploadHandler)
	s.g.POST("/upload", s.postUploadHandler)

	s.g.GET("/torrents", s.getTorrentsHandler)
	s.g.GET("/torrents/delete/:id", s.deleteTorrentHandler)

	return s.g.Run(fmt.Sprintf("%s:%d", host, port))
}
