package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) getTorrentsHandler(ctx *gin.Context) {
	list, err := s.TorrentService.List()
	if err != nil {
		s.l.Errorf("Load existing torrents list failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Load torrents list failed")
		return
	}

	page := struct {
		uiPage
		Torrents []string
	}{
		Torrents: list,
	}
	ctx.HTML(http.StatusOK, "multimedia.downloads.tmpl", &page)
}

func (s *Server) deleteTorrentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	l := s.l.WithField("id", id)
	if err := s.TorrentService.Remove(id); err != nil {
		s.l.Errorf("Remove failed: %s", err)
		displayError(ctx, http.StatusNotFound, "Remove torrent failed")
		return
	}
	l.Info("Removed")
	displayOK(ctx, "Torrent removed", "/torrents")
}
