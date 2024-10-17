package web

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/racoon-devel/raccoon-pirate/internal/selector"
)

const gracefulShutdownTimeout = 10 * time.Second

type Server struct {
	l                *log.Entry
	g                *gin.Engine
	srv              http.Server
	cache            sync.Map
	DiscoveryService DiscoveryService
	TorrentService   TorrentService
	Selector         selector.MovieSelector
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

	s.srv.Addr = fmt.Sprintf("%s:%d", host, port)
	s.srv.Handler = s.g.Handler()

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Serve HTTP failed: %s", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown failed: %s", err)
	}
}
