package webui

import (
	"embed"
	"io/fs"
)

//go:embed dist/**
var distFiles embed.FS

func DistFS() (fs.FS, bool) {
	subFS, err := fs.Sub(distFiles, "dist")
	if err != nil {
		return nil, false
	}

	if _, err := fs.Stat(subFS, "index.html"); err != nil {
		return nil, false
	}

	return subFS, true
}
