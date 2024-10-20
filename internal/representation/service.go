package representation

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"unicode"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

const mediaPerms = 0755

type Service struct {
	l          *log.Entry
	cfg        config.Representation
	flat       bool
	contentDir string
}

func New(cfg config.Representation, contentDirectory string) *Service {
	return &Service{
		l:          log.WithField("from", "representation"),
		cfg:        cfg,
		flat:       !cfg.Categories.Alphabet && !cfg.Categories.Genres && !cfg.Categories.Type && !cfg.Categories.Year,
		contentDir: contentDirectory,
	}
}

func (s *Service) Initialize(db Storage) error {
	if !s.cfg.Enabled {
		return nil
	}

	_ = os.RemoveAll(s.cfg.Directory)
	if err := os.MkdirAll(s.cfg.Directory, mediaPerms); err != nil {
		return fmt.Errorf("create representation directory failed: %s", err)
	}

	torrents, err := db.LoadAllTorrents()
	if err != nil {
		return fmt.Errorf("load torrents failed: %s", err)
	}

	for _, t := range torrents {
		s.createTorrentLayout(t)
	}
	return nil
}

func (s *Service) mapTorrent(t *model.Torrent) []string {
	basePath := filepath.Join(s.cfg.Directory, mapMediaTypeToDir(t.Type))
	if t.BelongsTo == "" || s.flat || t.Type == media.Other {
		return []string{basePath}
	}
	result := []string{}
	if t.Type == media.Movies {
		if s.cfg.Categories.Type {
			result = append(result, filepath.Join(basePath, mapMovieTypeToDir(t.MovieType), escape(t.BelongsTo)))
		}
		if s.cfg.Categories.Year && t.Year != 0 {
			result = append(result, filepath.Join(basePath, byYearDirectory, strconv.Itoa(int(t.Year)), escape(t.BelongsTo)))
		}
		if s.cfg.Categories.Alphabet {
			letter := unicode.ToUpper([]rune(escape(t.BelongsTo))[0])
			result = append(result, filepath.Join(basePath, byAlphabetDirectory, string(letter), escape(t.BelongsTo)))
		}
		if s.cfg.Categories.Genres {
			genres := t.GetGenres()
			for _, g := range genres {
				result = append(result, filepath.Join(basePath, byGenreDirectory, escape(g), escape(t.BelongsTo)))
			}
		}
	} else if t.Type == media.Music {
		result = append(result, filepath.Join(basePath, escape(t.BelongsTo)))
	}

	return result
}

func (s *Service) createTorrentLayout(t *model.Torrent) {
	var err error
	layout := s.mapTorrent(t)
	symlink := escape(t.Title)
	for _, dir := range layout {
		if err = os.MkdirAll(dir, mediaPerms); err != nil {
			s.l.Warnf("Create directory '%s' failed: %s", dir, err)
			continue
		}
		content := filepath.Join(s.contentDir, t.Title)
		if err = os.Symlink(content, filepath.Join(dir, symlink)); err != nil {
			s.l.Warnf("Create symlink to '%s' failed: %s", content, err)
		}
	}
}
