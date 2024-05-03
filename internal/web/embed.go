package web

import (
	"embed"
	"io/fs"
)

//go:embed content
var webFS embed.FS

//go:embed templates
var templatesFS embed.FS

func wrapFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
