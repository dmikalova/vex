// Package hotreload runs a build-and-serve loop that rebuilds artifacts and
// restarts a server whenever watched source files change. The `mage web` dev
// command uses it to rebuild the wasm client and restart the web server on every
// edit, so changes appear live in the browser.
package hotreload

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

// Config describes a watch → build → serve loop.
type Config struct {
	// Build regenerates build artifacts before each (re)start of the server. When
	// it fails the current server is left running, so the page is never blank while
	// an edit is mid-flight.
	Build func() error

	// Command and Args are the server to launch after a successful Build. It is
	// started fresh on each rebuild, and killed on the next change and on shutdown.
	Command string
	Args    []string

	// Extensions are the file suffixes whose changes trigger a rebuild (e.g. ".go",
	// ".css"). Files ending in "_test.go" never trigger one.
	Extensions []string

	// Root is the directory tree to watch. Defaults to ".".
	Root string

	// Interval is the file-polling period. Defaults to 500ms.
	Interval time.Duration

	// Settle is a debounce/quiet period: once a change is detected, the watcher
	// waits until the tree has been quiet (no further changes) for this long
	// before rebuilding, so a rapid burst of edits collapses into a single
	// rebuild. Zero (the default) rebuilds immediately on the first change.
	Settle time.Duration
}

// skipDirs are directories never worth walking for source changes.
var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "tmp": true, ".idea": true, ".vscode": true,
}

// Serve builds once, starts the server, then watches Root and rebuilds/restarts
// on every change until interrupted (Ctrl-C). It returns only on a failed initial
// build; the server is killed on shutdown so it never outlives the watcher.
func Serve(cfg Config) error {
	cfg.applyDefaults()

	if err := cfg.Build(); err != nil {
		return err
	}
	server := cfg.start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		kill(server)
		os.Exit(0)
	}()

	last := cfg.newestModTime()
	for range time.Tick(cfg.Interval) {
		t := cfg.newestModTime()
		if !t.After(last) {
			continue
		}
		last = cfg.settle(t)
		fmt.Println("hotreload: change detected, rebuilding…")
		if err := cfg.Build(); err != nil {
			fmt.Println("hotreload: build failed:", err)
			continue // keep the running server up
		}
		kill(server)
		server = cfg.start()
	}
	return nil
}

// settle waits out the debounce window: it keeps polling until the tree has been
// quiet for cfg.Settle, so a burst of rapid edits triggers one rebuild instead of
// many. It returns the newest mod time observed once things are stable. With
// Settle == 0 it returns immediately, rebuilding on the first change.
func (cfg *Config) settle(t time.Time) time.Time {
	for cfg.Settle > 0 {
		time.Sleep(cfg.Settle)
		newer := cfg.newestModTime()
		if !newer.After(t) {
			break // no change during the quiet window: settled
		}
		fmt.Println("hotreload: still changing, waiting to settle…")
		t = newer
	}
	return t
}

func (cfg *Config) applyDefaults() {
	if cfg.Root == "" {
		cfg.Root = "."
	}
	if cfg.Interval == 0 {
		cfg.Interval = 500 * time.Millisecond
	}
}

// start launches the server command with its output wired to this process's own
// stdout and stderr. A start failure is reported and yields a nil command.
func (cfg *Config) start() *exec.Cmd {
	cmd := exec.Command(cfg.Command, cfg.Args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("hotreload: start failed:", err)
		return nil
	}
	return cmd
}

// newestModTime is the latest modification time across the watched files under
// Root, or the zero time when none match.
func (cfg *Config) newestModTime() time.Time {
	var latest time.Time
	_ = filepath.WalkDir(cfg.Root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if !cfg.watches(d.Name()) {
			return nil
		}
		if info, err := d.Info(); err == nil && info.ModTime().After(latest) {
			latest = info.ModTime()
		}
		return nil
	})
	return latest
}

// watches reports whether a file name has a watched extension and is not a Go
// test file.
func (cfg *Config) watches(name string) bool {
	if strings.HasSuffix(name, "_test.go") {
		return false
	}
	for _, ext := range cfg.Extensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}

// kill terminates a running server and waits for it to exit. A nil command is a
// no-op.
func kill(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
	_ = cmd.Wait()
}
