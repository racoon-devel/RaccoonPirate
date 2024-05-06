package discovery

import (
	"bytes"
	"context"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client/torrents"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func asPtr[T any](val T) *T {
	return &val
}

func (s *Service) SearchTorrents(ctx context.Context, mov *model.Movie, season *int64) ([]*models.SearchTorrentsResult, error) {
	year := int64(mov.Year)

	req := torrents.SearchTorrentsParams{
		Limit:   asPtr[int64](searchResultsLimit),
		Q:       mov.Title,
		Season:  season,
		Strong:  asPtr(true),
		Type:    asPtr("movies"),
		Year:    &year,
		Context: ctx,
	}

	resp, err := s.cli.Torrents.SearchTorrents(&req, s.auth)
	if err != nil {
		return nil, err
	}

	return resp.Payload.Results, nil
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
