package web

import (
	"io"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/gin-gonic/gin"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type uploadPage struct {
	uiPage
	MediaType string
}

func (s *Server) getUploadHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "multimedia.upload.tmpl", &uploadPage{MediaType: "movies"})
}

func (s *Server) postUploadHandler(ctx *gin.Context) {
	mediaType := ctx.PostForm("media-type")
	contentType, ok := frontend.DetermineContentType(mediaType)
	if !ok {
		contentType = media.Movies
	}

	torrentRecord := model.Torrent{Type: contentType}
	file, err := ctx.FormFile("file")
	if err != nil {
		s.l.Errorf("Upload torrent file failed: %s", err)
		displayError(ctx, http.StatusBadRequest, "Upload file failed")
		return
	}
	f, err := file.Open()
	if err != nil {
		s.l.Errorf("Upload torrent file failed: %s", err)
		displayError(ctx, http.StatusBadRequest, "Upload file failed")
		return
	}

	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		s.l.Errorf("Upload torrent file failed: %s", err)
		displayError(ctx, http.StatusBadRequest, "Upload file failed")
		return
	}

	torrentRecord.Content = buf
	if err = s.TorrentService.Add(ctx, &torrentRecord); err != nil {
		s.l.Errorf("Add torrent failed %s", err)
		displayError(ctx, http.StatusInternalServerError, "Add torrent failed")
		return
	}
	displayOK(ctx, "Файл загружен", "/torrents?media-type="+mediaType)
}
