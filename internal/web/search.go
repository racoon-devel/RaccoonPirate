package web

import (
	"fmt"
	"net/http"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	uuid "github.com/satori/go.uuid"
)

type searchPage struct {
	uiPage
	Query     string
	MediaType string
	Movies    []*model.Movie
	Artists   []*model.Artist
	Albums    []*model.AlbumResult
}

func (s *Server) searchMovies(l *log.Entry, ctx *gin.Context, q string) []*model.Movie {
	movies, err := s.DiscoveryService.SearchMovies(ctx, q)
	if err != nil {
		l.Errorf("Search failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Something went wrong...")
		return nil
	}
	if len(movies) == 0 {
		l.Warnf("Nothing found")
		displayError(ctx, http.StatusBadRequest, "No results found")
		return nil
	}
	for _, mov := range movies {
		s.cache.Store(mov.ID, mov)
	}

	return movies
}

func (s *Server) searchMusic(l *log.Entry, ctx *gin.Context, q string) []model.Music {
	music, err := s.DiscoveryService.SearchMusic(ctx, q)
	if err != nil {
		l.Errorf("Search failed: %s", err)
		displayError(ctx, http.StatusInternalServerError, "Something went wrong...")
		return nil
	}
	if len(music) == 0 {
		l.Warnf("Nothing found")
		displayError(ctx, http.StatusBadRequest, "No results found")
		return nil
	}
	for _, m := range music {
		s.cache.Store(m.Title(), m)
	}

	return music
}

func extractArtistsAlbums(music []model.Music) (artists []*model.Artist, albums []*model.AlbumResult) {
	for _, m := range music {
		if m.IsArtist() {
			artists = append(artists, m.AsArtist())
		} else if m.IsAlbum() {
			albums = append(albums, m.AsAlbum())
		}
	}
	return
}

func (s *Server) searchHandler(ctx *gin.Context) {
	q := ctx.Query("q")

	mediaType := ctx.Query("media-type")
	contentType, ok := frontend.DetermineContentType(mediaType)

	if !ok {
		contentType = media.Movies
		mediaType = frontend.GetContentTypeID(media.Movies)
	}

	page := searchPage{
		Query:     q,
		MediaType: mediaType,
	}

	if q == "" {
		ctx.HTML(http.StatusOK, "multimedia.search.tmpl", &page)
		return
	}

	l := s.l.WithField("query", q).WithField("media-type", mediaType)
	l.Debugf("Search")

	switch contentType {
	case media.Movies:
		page.Movies = s.searchMovies(l, ctx, q)
		if page.Movies == nil {
			return
		}
	case media.Music:
		music := s.searchMusic(l, ctx, q)
		if music == nil {
			return
		}
		page.Artists, page.Albums = extractArtistsAlbums(music)
	case media.Other:
		id := uuid.NewV4().String()
		s.cache.Store(id, q)
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/add/%s?select=true", id))
		return
	}

	ctx.HTML(http.StatusOK, "multimedia.search.tmpl", &page)
}
