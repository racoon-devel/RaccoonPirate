package discovery

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client/torrents"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func asPtr[T any](val T) *T {
	return &val
}

func wait(ctx context.Context, interval time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(interval):
		return nil
	}
}

func (s *Service) searchTorrents(ctx context.Context, req *torrents.SearchTorrentsAsyncParams) ([]*models.SearchTorrentsResult, error) {
	task, err := s.cli.Torrents.SearchTorrentsAsync(req, s.auth)
	if err != nil {
		return nil, err
	}
	defer s.cli.Torrents.SearchTorrentsAsyncCancel(&torrents.SearchTorrentsAsyncCancelParams{ID: task.Payload.ID}, s.auth)

	for {
		if err = wait(ctx, time.Duration(task.Payload.PollIntervalMs)*time.Millisecond); err != nil {
			return nil, err
		}
		resp, err := s.cli.Torrents.SearchTorrentsAsyncStatus(&torrents.SearchTorrentsAsyncStatusParams{ID: task.Payload.ID, Context: ctx}, s.auth)
		if err != nil {
			return nil, err
		}
		switch *resp.Payload.Status {
		case "error":
			return nil, errors.New(resp.Payload.Error)
		case "ready":
			return resp.Payload.Results, nil
		}
	}
}

func (s *Service) SearchMovieTorrents(ctx context.Context, mov *model.Movie, season *int64) ([]*models.SearchTorrentsResult, error) {
	year := int64(mov.Year)

	req := torrents.SearchTorrentsAsyncParams{
		SearchParameters: torrents.SearchTorrentsAsyncBody{
			Limit:  int64(searchResultsLimit),
			Q:      &mov.Title,
			Strong: asPtr(true),
			Type:   "movies",
			Year:   year,
		},
		Context: ctx,
	}
	if season != nil {
		req.SearchParameters.Season = *season
	}

	return s.searchTorrents(ctx, &req)
}

func (s *Service) SearchMusicTorrents(ctx context.Context, m model.Music) ([]*models.SearchTorrentsResult, error) {
	q := m.Title()
	discrography := m.IsArtist()

	if m.IsAlbum() {
		q = m.AsAlbum().Artist + " " + q
	}
	req := torrents.SearchTorrentsAsyncParams{
		SearchParameters: torrents.SearchTorrentsAsyncBody{
			Limit:       int64(searchResultsLimit),
			Q:           &q,
			Strong:      asPtr(false),
			Type:        "music",
			Discography: &discrography,
		},
		Context: ctx,
	}

	return s.searchTorrents(ctx, &req)
}

func (s *Service) SearchOtherTorrents(ctx context.Context, q string) ([]*models.SearchTorrentsResult, error) {
	req := torrents.SearchTorrentsAsyncParams{
		SearchParameters: torrents.SearchTorrentsAsyncBody{
			Limit:  int64(searchResultsLimit),
			Q:      &q,
			Strong: asPtr(false),
			Type:   "others",
		},
		Context: ctx,
	}

	return s.searchTorrents(ctx, &req)
}

func (s *Service) GetTorrent(ctx context.Context, link string) ([]byte, error) {
	req := torrents.DownloadTorrentParams{
		Link:    link,
		Context: ctx,
	}
	buf := bytes.NewBuffer([]byte{})

	_, err := s.cli.Torrents.DownloadTorrent(&req, s.auth, buf)
	return buf.Bytes(), err
}
