package discovery

import (
	"context"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/client/music"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
)

func convertMusicInfo(in *models.SearchMusicResult) model.Music {
	switch *in.Type {
	case "artist":
		a := model.Artist{
			Name:       in.Artist,
			PictureUrl: in.Picture,
			Albums:     uint(in.AlbumsCount),
		}
		return model.PackMusic(&a)
	case "album":
		a := model.AlbumResult{
			Artist: in.Artist,
			Album: model.Album{
				Title:    in.Album,
				CoverUrl: in.Picture,
				Tracks:   uint(in.TracksCount),
				Genres:   in.Genres,
			},
		}
		return model.PackMusic(&a)
	default:
		return model.Music{}
	}
}

func (s *Service) SearchMusic(ctx context.Context, q string) ([]model.Music, error) {
	req := &music.SearchMusicParams{
		Limit:   asPtr(int64(searchResultsLimit)),
		Q:       q,
		Context: ctx,
	}

	resp, err := s.cli.Music.SearchMusic(req, s.auth)
	if err != nil {
		return nil, err
	}

	total := resp.Payload.Results

	result := make([]model.Music, 0, len(total))
	for _, item := range total {
		if *item.Type != "track" {
			result = append(result, convertMusicInfo(item))
		}
	}

	return result, nil
}
