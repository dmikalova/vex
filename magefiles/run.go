//go:build mage

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/sh"

	"github.com/dmikalova/vactrol/internal/hotreload"
)

// pkgsiteVersion pins the doc server (the renderer behind pkg.go.dev). Unlike the
// formatter and linter it is a read-only local viewer that no gate depends on, so
// bumping it is low-risk; it is pinned for reproducible fetches all the same.
const pkgsiteVersion = "v0.4.0"

// Docs serves this module's Go documentation at http://localhost:6060 with the
// same rendering as pkg.go.dev. It runs pkgsite over the local module via
// `go run`, so the tool never enters the module's own dependency graph, and
// blocks until Ctrl-C. The first run fetches pkgsite and may take a minute.
func Docs() error {
	const addr = "localhost:6060"
	return sh.RunV("go", "run",
		"golang.org/x/pkgsite/cmd/pkgsite@"+pkgsiteVersion, "-http", addr, ".")
}

// WebWasm builds the web client to WebAssembly (web/app.wasm). -trimpath makes
// the build reproducible; -ldflags="-s -w" drops debug info to shrink the bundle.
func WebWasm() error {
	return sh.RunWithV(
		map[string]string{"GOARCH": "wasm", "GOOS": "js"},
		"go", "build", "-trimpath", "-ldflags=-s -w", "-o", "web/app.wasm", "./cmd/web",
	)
}

// Web serves the wasm client with live rebuilds. It listens on
// http://localhost:8000 and rebuilds and restarts on any Go or CSS change so edits
// show up live. It also serves the Style gallery at /style, which no other
// deployment does. Each restart
// bumps go-app's version; the browser polls for it (see cmd/web devReload),
// reloads, and OnMount resumes the in-progress match. No external watcher needed;
// press Ctrl-C to stop.
//
// Set WEB_SETTLE to a Go duration (e.g. WEB_SETTLE=5s) to debounce rapid bursts
// of file changes: the watcher waits for that quiet period before rebuilding.
// Defaults to 5s when unset.
func Web() error {
	bin := filepath.Join(os.TempDir(), "vactrol-web-dev")
	// The dev server is the one place the Style gallery (/style) is meant to
	// exist, so this is where it is switched on; the served binary inherits it.
	if err := os.Setenv("VACTROL_STYLE", "1"); err != nil {
		return err
	}
	return hotreload.Serve(hotreload.Config{
		Build: func() error {
			if err := WebWasm(); err != nil {
				return err
			}
			// No WebAssets here: precompressing (max-level brotli) the multi-megabyte
			// wasm on every rebuild is what made the dev loop take ~50s. The server
			// falls back to the raw files when the .br/.gz siblings are absent, so the
			// dev server serves them straight from disk.
			return sh.Run("go", "build", "-o", bin, "./cmd/web")
		},
		Command:    bin,
		Extensions: []string{".go", ".css"},
		Settle:     webSettle(),
	})
}

// webSettle reads the WEB_SETTLE debounce duration from the environment; an unset
// or unparsable value defaults to 5s.
func webSettle() time.Duration {
	d, err := time.ParseDuration(os.Getenv("WEB_SETTLE"))
	if err != nil {
		return 5 * time.Second
	}
	return d
}
