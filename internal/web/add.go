package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RacoonMediaServer/rms-library/pkg/movsearch"
	"github.com/RacoonMediaServer/rms-library/pkg/selector"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"

	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
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
	if season == "" || season == "all" {
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

	mediaType := ""
	switch item := value.(type) {
	case *model.Movie:
		if !s.selectMovieTorrent(ctx, l, &q, item) {
			return
		}
		torrentRecord.ExpandByMovie(item)
		mediaType = frontend.GetContentTypeID(media.Movies)
	case model.Music:
		if !s.selectMusicTorrent(ctx, l, &q, item) {
			return
		}
		torrentRecord.ExpandByMusic(item)
		mediaType = frontend.GetContentTypeID(media.Music)
	case string:
		if !s.selectOtherTorrent(ctx, l, &q, item) {
			return
		}
		torrentRecord.Type = media.Other
		mediaType = frontend.GetContentTypeID(media.Other)
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

	if err = s.TorrentService.Add(ctx, torrentRecord, content); err != nil {
		l.Errorf("Add torrent failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Add torrent failed")
		return
	}

	l.Info("Added")
	displayOK(ctx, "Torrent added", "/torrents?media-type="+mediaType)
}

func (s *Server) selectMovieTorrent(ctx *gin.Context, l *log.Entry, q *addQuery, mov *model.Movie) bool {
	l = l.WithField("media-type", frontend.GetContentTypeID(media.Movies)).WithField("title", mov.Title)

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

	criteria := s.SelectCriterion
	if q.season == "all" {
		criteria = selector.CriteriaCompact
	}

	opts := selector.Options{
		Log:       l,
		Criteria:  criteria,
		MediaType: media.Movies,
	}

	// Imporved torrent movie selection
	if !q.selectTorrent {
		s.smartSelectMovieTorrent(ctx, l, mov, getSeasonNo(q.season), opts)
		return false
	}

	// Search torrents
	list, err := s.SmartSearchService.SearchMovieTorrents(ctx, mov, getSeasonNo(q.season))
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
		s.Selector.Sort(list, opts)
		page := selectTorrentPage{
			ID:       q.id,
			Select:   q.selectTorrent,
			Torrents: list,
		}
		ctx.HTML(http.StatusOK, "multimedia.download.select.tmpl", &page)
		return false
	}

	// Auto selection of torrent
	selected := s.Selector.Select(list, opts)
	q.torrent = *selected.Link
	return true
}

func (s *Server) smartSelectMovieTorrent(ctx *gin.Context, l *log.Entry, mov *model.Movie, season *int64, opts selector.Options) {
	torrents, err := s.SmartSearchService.SmartSearchMovieTorrents(ctx, mov, s.Selector, opts, season)
	if err != nil {
		l.Errorf("Find torrents failed: %s", err)
		if errors.Is(err, movsearch.ErrAnyTorrentsNotFound) {
			displayError(ctx, http.StatusNotFound, "Nothing found")
			return
		}
		displayError(ctx, http.StatusInternalServerError, "Search torrents failed")
		return
	}

	l.Infof("FOUND %d torrent", len(torrents))

	torrentRecord := &internalModel.Torrent{}
	torrentRecord.ExpandByMovie(mov)

	somethingAdded := false
	for _, torrent := range torrents {
		if err = s.TorrentService.Add(ctx, torrentRecord, torrent); err != nil {
			l.Errorf("Add torrent failed: %s", err)
		} else {
			somethingAdded = true
		}
	}

	if !somethingAdded {
		displayError(ctx, http.StatusInternalServerError, "Add torrents failed")
	} else {
		displayOK(ctx, "Added", "/torrents?media-type="+frontend.GetContentTypeID(media.Movies))
	}
}

func (s *Server) selectMusicTorrent(ctx *gin.Context, l *log.Entry, q *addQuery, m model.Music) bool {
	l = l.WithField("media-type", frontend.GetContentTypeID(media.Music)).WithField("title", m.Title())

	// If torrent has been selected - just return
	if q.torrent != "" {
		return true
	}

	// Search torrents
	list, err := s.DiscoveryService.SearchMusicTorrents(ctx, m)
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

	opts := selector.Options{
		Log:         l,
		Criteria:    s.SelectCriterion,
		MediaType:   media.Music,
		Query:       m.Title(),
		Discography: m.IsArtist(),
	}

	// Select concrete torrent manually by user
	if q.selectTorrent {
		s.Selector.Sort(list, opts)
		page := selectTorrentPage{
			ID:       q.id,
			Select:   q.selectTorrent,
			Torrents: list,
		}
		ctx.HTML(http.StatusOK, "multimedia.download.select.tmpl", &page)
		return false
	}

	// Auto selection of torrent
	selected := s.Selector.Select(list, opts)
	q.torrent = *selected.Link
	return true
}

func (s *Server) selectOtherTorrent(ctx *gin.Context, l *log.Entry, q *addQuery, tq string) bool {
	l = l.WithField("media-type", frontend.GetContentTypeID(media.Other)).WithField("title", tq)

	// If torrent has been selected - just return
	if q.torrent != "" {
		return true
	}

	// Search torrents
	list, err := s.DiscoveryService.SearchOtherTorrents(ctx, tq)
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

	opts := selector.Options{
		Log:       l,
		Criteria:  s.SelectCriterion,
		MediaType: media.Other,
		Query:     tq,
	}

	// Select concrete torrent manually by user
	if q.selectTorrent {
		s.Selector.Sort(list, opts)
		page := selectTorrentPage{
			ID:       q.id,
			Select:   q.selectTorrent,
			Torrents: list,
		}
		ctx.HTML(http.StatusOK, "multimedia.download.select.tmpl", &page)
		return false
	}

	// Auto selection of torrent
	selected := s.Selector.Select(list, opts)
	q.torrent = *selected.Link
	return true
}
