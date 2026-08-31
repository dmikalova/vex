//go:build mage

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/sh"

	"github.com/dmikalova/vactrol/internal/hotreload"
)

// Tui launches the vactrol TUI (card explorer / play a hotseat game).
func Tui() error {
	return sh.RunV("go", "run", "./cmd/tui")
}

// WebWasm builds the web client to WebAssembly (web/app.wasm). -trimpath makes
// the build reproducible; -ldflags="-s -w" drops debug info to shrink the bundle.
func WebWasm() error {
	return sh.RunWithV(
		map[string]string{"GOARCH": "wasm", "GOOS": "js"},
		"go", "build", "-trimpath", "-ldflags=-s -w", "-o", "web/app.wasm", "./cmd/web",
	)
}

// Web builds the wasm client and serves it at http://localhost:8000, rebuilding
// and restarting on any Go or CSS change so edits show up live. Each restart
// bumps go-app's version; the browser polls for it (see cmd/web devReload),
// reloads, and OnMount resumes the in-progress match. No external watcher needed;
// press Ctrl-C to stop.
//
// Set WEB_SETTLE to a Go duration (e.g. WEB_SETTLE=5s) to debounce rapid bursts
// of file changes: the watcher waits for that quiet period before rebuilding.
// Defaults to 5s when unset.
func Web() error {
	bin := filepath.Join(os.TempDir(), "vactrol-web-dev")
	return hotreload.Serve(hotreload.Config{
		Build: func() error {
			if err := WebWasm(); err != nil {
				return err
			}
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
