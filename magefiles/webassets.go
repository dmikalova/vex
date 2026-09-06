//go:build mage

package main

import (
	"compress/gzip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/svg"
)

// WebAssets minifies and precompresses the static web assets the server ships.
// It writes .br (brotli) and .gz (gzip) siblings next to each servable file under
// web/ so cmd/web streams a prebuilt body instead of compressing per request;
// text assets (CSS, SVG) are minified first. Already-compressed binaries (.png,
// .woff2) are left alone. Run after WebWasm so app.wasm is compressed too. The
// dev server (mage web) skips this step and serves the raw files from disk.
func WebAssets() error {
	if err := cleanPrecompressed(); err != nil {
		return err
	}
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	files, err := servableAssets()
	if err != nil {
		return err
	}
	for _, p := range files {
		if err := buildAsset(m, p); err != nil {
			return fmt.Errorf("web asset %s: %w", p, err)
		}
	}
	return nil
}

// servableAssets lists the files under web/ that are worth precompressing: it
// skips directories, the precompressed siblings themselves, and formats that are
// already compressed (.png, .woff2).
func servableAssets() ([]string, error) {
	var out []string
	err := filepath.WalkDir("web", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		switch strings.ToLower(filepath.Ext(p)) {
		case ".br", ".gz", ".png", ".woff2":
			return nil
		}
		out = append(out, p)
		return nil
	})
	return out, err
}

// buildAsset minifies a text asset (leaving other types as-is) and writes its
// brotli and gzip siblings.
func buildAsset(m *minify.M, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	body := raw
	switch strings.ToLower(filepath.Ext(path)) {
	case ".css":
		body, err = m.Bytes("text/css", raw)
	case ".svg":
		body, err = m.Bytes("image/svg+xml", raw)
	}
	if err != nil {
		return err
	}
	if err := writeBrotli(path+".br", body); err != nil {
		return err
	}
	return writeGzip(path+".gz", body)
}

// writeBrotli writes body brotli-compressed at the maximum level.
func writeBrotli(path string, body []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := brotli.NewWriterLevel(f, brotli.BestCompression)
	if _, err := w.Write(body); err != nil {
		_ = w.Close()
		_ = f.Close()
		return err
	}
	if err := w.Close(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// writeGzip writes body gzip-compressed at the maximum level.
func writeGzip(path string, body []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		_ = f.Close()
		return err
	}
	if _, err := w.Write(body); err != nil {
		_ = w.Close()
		_ = f.Close()
		return err
	}
	if err := w.Close(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// cleanPrecompressed removes stale .br/.gz siblings so a deleted or renamed asset
// leaves no orphaned compressed copy behind.
func cleanPrecompressed() error {
	return filepath.WalkDir("web", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		switch strings.ToLower(filepath.Ext(p)) {
		case ".br", ".gz":
			return os.Remove(p)
		}
		return nil
	})
}
