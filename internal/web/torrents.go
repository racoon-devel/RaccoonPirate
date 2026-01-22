package web

import (
	"net/http"
	"sort"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/gin-gonic/gin"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type torrentsPage struct {
	uiPage
	Torrents  []*model.Torrent
	MediaType string
}

func (s *Server) getTorrentsHandler(ctx *gin.Context) {
	mediaType := ctx.Query("media-type")
	page := torrentsPage{}

	contentType, ok := frontend.DetermineContentType(mediaType)
	if !ok {
		contentType = media.Movies
		mediaType = frontend.GetContentTypeID(contentType)
	}

	list, err := s.TorrentService.GetTorrentsList(contentType)
	if err != nil {
		s.l.Errorf("Load existing torrents list failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Load torrents list failed")
		return
	}

	sort.SliceStable(list, func(i, j int) bool {
		if list[i].BelongsTo == list[j].BelongsTo {
			return list[i].Title < list[j].Title
		}
		return list[i].BelongsTo < list[j].BelongsTo
	})

	page.Torrents = list
	page.MediaType = mediaType

	ctx.HTML(http.StatusOK, "multimedia.downloads.tmpl", &page)
}

func (s *Server) deleteTorrentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	mediaType := ctx.Query("media-type")
	l := s.l.WithField("id", id)
	if err := s.TorrentService.Remove(ctx, id); err != nil {
		s.l.Errorf("Remove failed: %s", err)
		displayError(ctx, http.StatusNotFound, "Remove torrent failed")
		return
	}
	l.Info("Removed")
	displayOK(ctx, "Torrent removed", "/torrents?media-type="+mediaType)
}
