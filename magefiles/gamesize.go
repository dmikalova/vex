//go:build mage

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/andybalholm/brotli"
	"github.com/magefile/mage/mg"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/svg"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// GameSize reports web and state sizes. It prints the in-memory GameState size,
// the number of implemented cards, and — after building a fresh wasm — the size
// of the WebAssembly bundle and the other static assets, both raw and compressed.
// Text assets are minified before compression so the numbers match what the
// server sends.
func (Tool) GameSize() error {
	mg.Deps(WebWasm)

	stateSize := int64(unsafe.Sizeof(engine.GameState{}))
	fmt.Printf("Game state:        %s (%d bytes per undo snapshot)\n",
		humanBytes(stateSize), stateSize)
	fmt.Printf("Cards implemented: %d\n\n", len(cards.All()))

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	files, err := servableAssets()
	if err != nil {
		return err
	}

	var wasm, others assetTally
	for _, p := range files {
		s, err := assetSize(m, p)
		if err != nil {
			return fmt.Errorf("measure %s: %w", p, err)
		}
		if strings.EqualFold(filepath.Ext(p), ".wasm") {
			wasm.add(s)
		} else {
			others.add(s)
		}
	}

	printTally("WASM bundle", wasm)
	printTally(fmt.Sprintf("Other assets (%d files)", others.count), others)

	total := wasm
	total.count += others.count
	total.raw += others.raw
	total.brotli += others.brotli
	total.gzip += others.gzip
	printTally("Total shipped", total)
	return nil
}

// assetSizes holds the raw, brotli, and gzip byte counts of one shipped asset.
type assetSizes struct {
	raw, brotli, gzip int64
}

// assetTally accumulates the sizes of a group of assets.
type assetTally struct {
	count             int
	raw, brotli, gzip int64
}

func (t *assetTally) add(s assetSizes) {
	t.count++
	t.raw += s.raw
	t.brotli += s.brotli
	t.gzip += s.gzip
}

// assetSize measures one asset as the server ships it: text assets (CSS, SVG) are
// minified first, then the body is compressed with brotli and gzip to count the
// wire size without writing the sibling files.
func assetSize(m *minify.M, path string) (assetSizes, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return assetSizes{}, err
	}
	body := raw
	switch strings.ToLower(filepath.Ext(path)) {
	case ".css":
		body, err = m.Bytes("text/css", raw)
	case ".svg":
		body, err = m.Bytes("image/svg+xml", raw)
	}
	if err != nil {
		return assetSizes{}, err
	}
	return assetSizes{
		raw:    int64(len(body)),
		brotli: brotliSize(body),
		gzip:   gzipSize(body),
	}, nil
}

// brotliSize returns the length of body compressed with brotli at the level
// WebAssets ships, without keeping the compressed bytes.
func brotliSize(body []byte) int64 {
	var n countWriter
	w := brotli.NewWriterLevel(&n, brotli.BestCompression)
	_, _ = w.Write(body)
	_ = w.Close()
	return int64(n)
}

// gzipSize returns the length of body compressed with gzip at the level WebAssets
// ships, without keeping the compressed bytes.
func gzipSize(body []byte) int64 {
	var n countWriter
	w, _ := gzip.NewWriterLevel(&n, gzip.BestCompression)
	_, _ = w.Write(body)
	_ = w.Close()
	return int64(n)
}

// countWriter is an io.Writer that discards its input and counts the bytes.
type countWriter int64

func (c *countWriter) Write(p []byte) (int, error) {
	*c += countWriter(len(p))
	return len(p), nil
}

var _ io.Writer = (*countWriter)(nil)

// printTally prints one group's raw and compressed totals as an indented block.
func printTally(label string, t assetTally) {
	fmt.Printf("%s:\n", label)
	fmt.Printf("  raw:    %10s\n", humanBytes(t.raw))
	fmt.Printf("  gzip:   %10s\n", humanBytes(t.gzip))
	fmt.Printf("  brotli: %10s\n\n", humanBytes(t.brotli))
}

// humanBytes renders a byte count in binary units (KiB, MiB) with one decimal,
// leaving counts under a kibibyte in plain bytes.
func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for n/div >= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGT"[exp])
}
