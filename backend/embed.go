package main

import (
	"embed"
	"io/fs"
	"os"

	"github.com/iwanesko/chess-web-site/backend/internal/server"
)

// distFS embeds the built frontend. A committed placeholder keeps this
// compilable before the first frontend build; the Docker image and the local
// build script overwrite backend/dist with the real Vite output.
//
//go:embed all:dist
var distFS embed.FS

// resolveStatic picks the static file source: the on-disk dev directory when
// DEV_STATIC_DIR is set, otherwise the binary-embedded build.
func resolveStatic(cfg server.Config) (fs.FS, error) {
	if cfg.DevStaticDir != "" {
		return os.DirFS(cfg.DevStaticDir), nil
	}
	return fs.Sub(distFS, "dist")
}
