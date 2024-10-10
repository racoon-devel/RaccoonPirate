package web

import (
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/racoon-devel/raccoon-pirate/internal/selector"
	"net/http"
	"strconv"
)

type selectSeasonPage struct {
	uiPage
	ID      string
	Select  bool
	Seasons []uint
}

type selectTorrentPage struct {
	uiPage
	ID       string
	Select   bool
	Torrents []*models.SearchTorrentsResult
}

func getSeasonNo(season string) *int64 {
	if season == "" {
		return nil
	}
	no, err := strconv.ParseUint(season, 10, 32)
	if err != nil {
		return nil
	}

	result := int64(no)
	return &result
}

func (s *Server) addHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	l := s.l.WithField("id", id)
	season := ctx.Query("season")
	selectTorrent := ctx.Query("select") == "true"
	torrent := ctx.Query("torrent")

	mov, ok := s.movieFromCache(id)
	if !ok {
		l.Errorf("Item not found")
		displayError(ctx, http.StatusNotFound, "Info about media not found in a cache")
		return
	}

	if mov.Type == model.MovieType_TvSeries && mov.Seasons != 0 && season == "" && torrent == "" {
		page := selectSeasonPage{
			ID:      id,
			Select:  selectTorrent,
			Seasons: iotaSeasons(mov.Seasons),
		}
		ctx.HTML(http.StatusOK, "multimedia.download.tmpl", &page)
		return
	}

	if torrent == "" {
		list, err := s.DiscoveryService.SearchTorrents(ctx, mov, getSeasonNo(season))
		if err != nil {
			l.Errorf("Search torrents failed: %s", err)
			displayError(ctx, http.StatusInternalServerError, "Search torrents failed")
			return
		}
		if len(list) == 0 {
			l.Warnf("Nothing found")
			displayError(ctx, http.StatusNotFound, "Nothing found")
			return
		}
		if selectTorrent {
			page := selectTorrentPage{
				ID:       id,
				Select:   selectTorrent,
				Torrents: list,
			}
			ctx.HTML(http.StatusOK, "multimedia.download.select.tmpl", &page)
			return
		}

		selected := s.Selector.Select(l, selector.CriteriaQuality, list)
		torrent = *selected.Link
	}

	content, err := s.DiscoveryService.GetTorrent(ctx, torrent)
	if err != nil {
		l.Errorf("Download torrent '%s' failed: %s", torrent, err)
		displayError(ctx, http.StatusInternalServerError, "Download torrent failed")
		return
	}

	if err = s.TorrentService.Add(content); err != nil {
		l.Errorf("Add torrent failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Add torrent failed")
		return
	}

	l.Info("Added")
	displayOK(ctx, "Torrent added", "/torrents")
}
