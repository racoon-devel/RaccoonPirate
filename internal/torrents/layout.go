package torrents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type layout struct {
	baseDir            string
	pieceCompletionDir string
	cacheDir           string
	contentDir         string
	itemsDir           string
	torrentsDir        string
}

func newLayout(baseDir string) layout {
	return layout{
		baseDir:            baseDir,
		pieceCompletionDir: filepath.Join(baseDir, "piece-completion"),
		cacheDir:           filepath.Join(baseDir, "cache"),
		contentDir:         filepath.Join(baseDir, "content"),
		itemsDir:           filepath.Join(baseDir, "items"),
		torrentsDir:        filepath.Join(baseDir, "torrents"),
	}
}

func (l layout) makeLayout() error {
	if err := os.MkdirAll(l.pieceCompletionDir, 0744); err != nil {
		return fmt.Errorf("create piece completion directory failed: %w", err)
	}
	if err := os.MkdirAll(l.cacheDir, 0744); err != nil {
		return fmt.Errorf("create cache directory failed: %w", err)
	}
	if err := os.MkdirAll(l.contentDir, 0744); err != nil {
		return fmt.Errorf("create content directory failed: %w", err)
	}
	if err := os.MkdirAll(l.torrentsDir, 0744); err != nil {
		return fmt.Errorf("create torrents directory failed: %w", err)
	}
	return nil
}

func (l layout) ListTorrents() (map[string][][]byte, error) {
	var result [][]byte
	files, err := os.ReadDir(l.torrentsDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			data, err := os.ReadFile(filepath.Join(l.torrentsDir, f.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		}
	}
	return map[string][][]byte{mainRoute: result}, nil
}

func escape(s string) string {
	return strings.Replace(s, "/", "", -1)
}
