//go:build mage

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/sh"

	"github.com/dmikalova/vex/internal/hotreload"
)

// pkgsiteVersion pins the doc server (the renderer behind pkg.go.dev). Unlike the
// formatter and linter it is a read-only local viewer that no gate depends on, so
// bumping it is low-risk; it is pinned for reproducible fetches all the same.
const pkgsiteVersion = "v0.4.0"

// Docs serves the local API docs. It runs pkgsite over the local module via
// `go run`, so the tool never enters the module's own dependency graph, and
// blocks until Ctrl-C. The first run fetches pkgsite and may take a minute.
func Docs() error {
	const addr = "localhost:6060"
	return sh.RunV("go", "run",
		"golang.org/x/pkgsite/cmd/pkgsite@"+pkgsiteVersion, "-http", addr, ".")
}

// WebWasm builds the web client to WebAssembly (web/app.wasm).
func WebWasm() error {
	return wasmBuild("web/app.wasm")
}

// wasmBuild compiles the web client to WebAssembly at out. -trimpath makes the
// build reproducible; -ldflags="-s -w" drops debug info to shrink the bundle.
func wasmBuild(out string) error {
	return sh.RunWithV(
		map[string]string{"GOARCH": "wasm", "GOOS": "js"},
		"go", "build", "-trimpath", "-ldflags=-s -w", "-o", out, "./cmd/web",
	)
}

// CapturePrune trims closed replay captures. It replays every entry the dev
// server has written into internal/web/testdata/capture and deletes the ones that
// are no longer open findings: the faults that have been fixed, and every entry
// recorded against a different command-log version or card pool, which can never
// be evidence about this tree again. What is left is the list of findings still to
// fix — the directory is that list, not an archive of every failure a dev server
// ever saw. Same contract as `mage corpusPrune`.
func CapturePrune() error {
	return sh.RunV("go", "run", "./magefiles/webcapture")
}

// Web serves the wasm client with live rebuilds. It listens on
// http://localhost:8000 and rebuilds and restarts on any Go or CSS change so edits
// show up live. It also serves the Style gallery at /style and the browser
// scenarios at /ui-test, which no other deployment does. Each restart
// bumps go-app's version; the browser polls for it (see cmd/web devReload),
// reloads, and OnMount resumes the in-progress match. No external watcher needed;
// press Ctrl-C to stop.
//
// Set WEB_SETTLE to a Go duration (e.g. WEB_SETTLE=5s) to debounce rapid bursts
// of file changes: the watcher waits for that quiet period before rebuilding.
// Defaults to 5s when unset.
func Web() error {
	bin := filepath.Join(os.TempDir(), "vex-web-dev")
	// The dev server is the one place the Style gallery (/style) is meant to
	// exist, so this is where it is switched on; the served binary inherits it.
	if err := os.Setenv("VEX_STYLE", "1"); err != nil {
		return err
	}
	// The browser scenarios (/ui-test) are switched on here for the same reason:
	// the dev server is where they are driven, and no other deployment offers them.
	if err := os.Setenv("VEX_UITEST", "1"); err != nil {
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

// UITest drives the browser scenarios headlessly.
//
// It runs internal/web/uitest, which builds web/app.wasm, serves cmd/web on a
// free port with /ui-test switched on, and points a headless Chrome at one page
// per scenario registered in internal/web. The scenarios are defined once, in
// internal/web/uitest_scenarios.go, so a new journey joins both the page a human
// watches and this suite.
//
// It is deliberately NOT part of mage ci:check or ci:test: the uitest build tag
// keeps the package out of ./..., so the shared CI workflow never needs a
// browser, and a multi-megabyte wasm build plus a browser boot does not belong
// in a gate that runs on every save. It sits beside mage profile and mage trace
// as a real-but-ungated target. Budget ~10-15s for the current suite; see
// docs/testing.md.
func UITest() error {
	return sh.RunV("go", "test", "-count=1", "-tags=uitest", "./internal/web/uitest/")
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
