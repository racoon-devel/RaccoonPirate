package web

import (
	"net/http"
	"strconv"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"

	internalModel "github.com/racoon-devel/raccoon-pirate/internal/model"
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

type addQuery struct {
	id            string
	season        string
	selectTorrent bool
	torrent       string
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
	torrentRecord := &internalModel.Torrent{}
	q := addQuery{
		id:            ctx.Param("id"),
		season:        ctx.Query("season"),
		selectTorrent: ctx.Query("select") == "true",
		torrent:       ctx.Query("torrent"),
	}
	l := s.l.WithField("id", q.id)

	value, ok := s.cache.Load(q.id)
	if !ok {
		l.Errorf("Item not found")
		displayError(ctx, http.StatusNotFound, "Info about media not found in a cache")
		return
	}

	switch item := value.(type) {
	case *model.Movie:
		if !s.selectMovieTorrent(ctx, l, &q, item) {
			return
		}
		torrentRecord.ExpandByMovie(item)
	default:
		l.Errorf("Unknown type of media: %T", item)
		displayError(ctx, http.StatusInternalServerError, "Type of media is unsupported")
		return
	}

	content, err := s.DiscoveryService.GetTorrent(ctx, q.torrent)
	if err != nil {
		l.Errorf("Download torrent '%s' failed: %s", q.torrent, err)
		displayError(ctx, http.StatusInternalServerError, "Download torrent failed")
		return
	}

	if err = s.TorrentService.Add(torrentRecord, content); err != nil {
		l.Errorf("Add torrent failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Add torrent failed")
		return
	}

	l.Info("Added")
	displayOK(ctx, "Torrent added", "/torrents")
}

func (s *Server) selectMovieTorrent(ctx *gin.Context, l *log.Entry, q *addQuery, mov *model.Movie) bool {
	l = l.WithField("media-type", "movie").WithField("title", mov.Title)

	// Select season in tv-series case
	if mov.Type == model.MovieType_TvSeries && mov.Seasons != 0 && q.season == "" && q.torrent == "" {
		page := selectSeasonPage{
			ID:      q.id,
			Select:  q.selectTorrent,
			Seasons: iotaSeasons(mov.Seasons),
		}
		ctx.HTML(http.StatusOK, "multimedia.download.tmpl", &page)
		return false
	}

	// If torrent has been selected - just return
	if q.torrent != "" {
		return true
	}

	// Search torrents
	list, err := s.DiscoveryService.SearchTorrents(ctx, mov, getSeasonNo(q.season))
	if err != nil {
		l.Errorf("Search torrents failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Search torrents failed")
		return false
	}

	if len(list) == 0 {
		l.Warnf("Nothing found")
		displayError(ctx, http.StatusNotFound, "Nothing found")
		return false
	}

	// Select concrete torrent manually by user
	if q.selectTorrent {
		s.Selector.Sort(l, s.SelectCriterion, list)
		page := selectTorrentPage{
			ID:       q.id,
			Select:   q.selectTorrent,
			Torrents: list,
		}
		ctx.HTML(http.StatusOK, "multimedia.download.select.tmpl", &page)
		return false
	}

	// Auto selection of torrent
	selected := s.Selector.Select(l, s.SelectCriterion, list)
	q.torrent = *selected.Link
	return true
}
