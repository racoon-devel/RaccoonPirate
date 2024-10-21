package representation

import (
	"io"
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

type Service interface {
	Register(t *model.Torrent, pathToContent string)
	Unregister(t *model.Torrent)
	Clean()
}

type serviceImpl struct {
	l    *log.Entry
	cfg  config.Representation
	flat bool
}

func New(cfg config.Representation) Service {
	if !cfg.Enabled {
		return &disabledImpl{}
	}

	_ = os.RemoveAll(cfg.Directory)

	return &serviceImpl{
		l:    log.WithField("from", "representation"),
		cfg:  cfg,
		flat: !cfg.Categories.Alphabet && !cfg.Categories.Genres && !cfg.Categories.Type && !cfg.Categories.Year,
	}
}

func (s *serviceImpl) mapTorrent(t *model.Torrent) []string {
	basePath := mapMediaTypeToDir(t.Type)
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

func (s *serviceImpl) Register(t *model.Torrent, pathToContent string) {
	layout := s.mapTorrent(t)
	symlink := escape(t.Title)
	for _, dir := range layout {
		fullPath := filepath.Join(s.cfg.Directory, dir)
		if err := os.MkdirAll(fullPath, mediaPerms); err != nil {
			s.l.Warnf("Create directory '%s' failed: %s", fullPath, err)
			continue
		}
		if err := os.Symlink(pathToContent, filepath.Join(fullPath, symlink)); err != nil {
			s.l.Warnf("Create symlink to '%s' failed: %s", pathToContent, err)
		}
	}
}

func isEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func (s *serviceImpl) removeIfEmpty(dir string) {
	if dir == "" {
		return
	}
	for {
		path, _ := filepath.Split(dir)
		fullPath := filepath.Join(s.cfg.Directory, dir)
		empty, err := isEmpty(fullPath)
		if err != nil || !empty {
			return
		}
		_ = os.Remove(fullPath)
		dir = path
	}
}

func (s *serviceImpl) Unregister(t *model.Torrent) {
	layout := s.mapTorrent(t)
	symlink := escape(t.Title)
	for _, dir := range layout {
		_ = os.Remove(filepath.Join(s.cfg.Directory, dir, symlink))
		s.removeIfEmpty(dir)
	}
}

// Clean implements Service.
func (s *serviceImpl) Clean() {
	_ = os.RemoveAll(s.cfg.Directory)
}
