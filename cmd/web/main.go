// Command web serves and runs the Vactrol browser client. The same binary is
// built two ways: compiled to WebAssembly it runs the interactive UI in the
// browser (app.RunWhenOnBrowser blocks there); built natively it serves that
// wasm bundle and the required go-app resources over HTTP.
package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/web"
)

func main() {
	app.Route("/", web.NewGame)
	app.Route("/rulebook", web.NewRulebook)
	app.Route("/glossary", web.NewGlossary)
	if styleEnabled() {
		app.Route("/style", web.NewStyle)
	}
	app.RunWhenOnBrowser()

	version := resourceVersion()
	assets := assetPaths()

	// Serve a fullscreen web app manifest at go-app's manifest path. go-app
	// hardcodes display "standalone"; overriding the route makes an installed PWA
	// launch immersively, hiding the Android status and navigation bars.
	http.HandleFunc("/manifest.webmanifest", serveManifest)

	http.Handle("/", &app.Handler{
		Name:            "Vactrol",
		ShortName:       "Vactrol",
		Title:           "Vactrol",
		Description:     "Vactrol — a KeyForge-style card game, playable in the browser.",
		BackgroundColor: "#1c1c1b",
		ThemeColor:      "#1c1c1b",
		Icon: app.Icon{
			// Raster icons at the standard PWA sizes make the app installable on
			// Android/Chrome, which wants a 192 and a 512 PNG; the SVG stays as the
			// scalable favicon, and the maskable 512 adapts to each platform's shape.
			SVG:      "/web/assets/favicon.svg",
			Default:  "/web/assets/icon-192.png",
			Large:    "/web/assets/icon-512.png",
			Maskable: "/web/assets/icon-512.png",
		},
		// Plain CSS served as a static file from web/ — no CDN, no build step. The dev
		// server serves it from disk, so editing web/app.css and refreshing the browser
		// applies changes with the server left running (no restart).
		Styles:     []string{"/web/app.css"},
		RawHeaders: []string{bootStyle, appleTouchIcon, boardScript, devReloadScript},
		// The icons are fetched one <img> at a time as the board draws, so without
		// precaching them a client that has the wasm cached but no server draws a board
		// of broken images. The service worker serves a cached copy first and only
		// re-fetches when Version changes — which every asset edit does — so an offline
		// client keeps its icons and an online one never keeps a stale set.
		CacheableResources: assets,
		// Byte size of app.wasm so go-app's loader accurately computes loading
		// progress even when serving gzip-compressed / chunked responses.
		WasmContentLength: wasmContentLength(),
		// Version keys go-app's service-worker cache. Left empty it defaults to the
		// app.wasm hash, so a CSS-only edit (wasm unchanged) would keep serving the
		// cached stylesheet. Hashing app.css too makes every asset edit bump it.
		Version: version,
		// The short build id the client shows, so a page and the server that built
		// it can be matched by eye.
		Env: map[string]string{
			"VACTROL_BUILD": buildID(version),
			// Passed down so the wasm client registers the same routes the server
			// serves; without it the gallery's page would be served and render blank.
			styleEnv: os.Getenv(styleEnv),
		},
	})

	// Cloud Run injects PORT; fall back to 8000 for local `mage web`.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port
	log.Printf("Vactrol web client on http://localhost%s (build %s)", addr, buildID(version))
	if err := http.ListenAndServe(addr, gzipHandler(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

// gzipHandler wraps h to gzip-encode responses for clients that accept it. The
// biggest asset is app.wasm, which compresses to roughly half its size, so this
// is the largest lever on first-load download time.
func gzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")
		gz := gzip.NewWriter(w)
		defer func() { _ = gz.Close() }()
		h.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, gz: gz}, r)
	})
}

// gzipResponseWriter streams a handler's output through a gzip.Writer. It drops
// Content-Length (the compressed length differs) on the first write so the
// server falls back to chunked transfer encoding.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz          *gzip.Writer
	wroteHeader bool
}

func (g *gzipResponseWriter) WriteHeader(status int) {
	if g.wroteHeader {
		return
	}
	g.wroteHeader = true
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(status)
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.wroteHeader {
		g.WriteHeader(http.StatusOK)
	}
	return g.gz.Write(b)
}

// resourceVersion hashes the served static assets — the wasm bundle, the
// stylesheet, and every icon — so the go-app Handler version changes whenever any
// of them does, busting the service worker's precache on every edit.
// buildID shortens a resource version to the few characters a human needs to
// tell one build from the next.
func buildID(version string) string {
	if len(version) <= 4 {
		return version
	}
	return version[:4]
}

func resourceVersion() string {
	h := sha256.New()
	for _, p := range append([]string{"web/app.wasm", "web/app.css"}, assetFiles()...) {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))[:12]
}

// assetFiles lists every file under web/assets, sorted, so hashing them yields
// the same version for the same bytes on every start.
func assetFiles() []string {
	var out []string
	_ = filepath.WalkDir("web/assets", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

// assetPaths turns the asset files into the URL paths they are served under, for
// the service worker's precache list.
func assetPaths() []string {
	files := assetFiles()
	out := make([]string, len(files))
	for i, p := range files {
		out[i] = "/" + filepath.ToSlash(p)
	}
	return out
}

// wasmContentLength returns the uncompressed byte size of web/app.wasm as a string
// so go-app can calculate download progress when serving responses with gzip.
func wasmContentLength() string {
	info, err := os.Stat("web/app.wasm")
	if err != nil {
		return ""
	}
	return strconv.FormatInt(info.Size(), 10)
}

// serveManifest writes a fullscreen web app manifest, overriding the one go-app
// generates (which hardcodes display "standalone"). ServeMux routes this exact
// path here rather than to the go-app handler registered at "/".
func serveManifest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	_, _ = w.Write([]byte(webManifest))
}

// webManifest mirrors the Handler's name, colors, and icon but sets display to
// fullscreen so an installed PWA takes the whole screen.
const webManifest = `{
  "short_name": "Vactrol",
  "name": "Vactrol",
  "description": "Vactrol — a KeyForge-style card game, playable in the browser.",
  "scope": "/",
  "start_url": "/",
  "background_color": "#1c1c1b",
  "theme_color": "#1c1c1b",
  "display": "fullscreen",
  "display_override": ["fullscreen", "standalone"],
  "icons": [
    { "src": "/web/assets/favicon.svg", "type": "image/svg+xml", "sizes": "any", "purpose": "any" },
    { "src": "/web/assets/icon-192.png", "type": "image/png", "sizes": "192x192", "purpose": "any" },
    { "src": "/web/assets/icon-512.png", "type": "image/png", "sizes": "512x512", "purpose": "any" },
    { "src": "/web/assets/app-icon.svg", "type": "image/svg+xml", "sizes": "any", "purpose": "maskable" },
    { "src": "/web/assets/icon-192.png", "type": "image/png", "sizes": "192x192", "purpose": "maskable" },
    { "src": "/web/assets/icon-512.png", "type": "image/png", "sizes": "512x512", "purpose": "maskable" }
  ]
}`

// appleTouchIcon points iOS Safari at the home-screen icon it uses when the app
// is added to the home screen; go-app's Icon struct has no Apple-touch field, so
// the link tag is injected into the head directly.
const appleTouchIcon = `<link rel="apple-touch-icon" href="/web/assets/apple-touch-icon.png">`

// bootStyle is an inline <head> stylesheet that paints the dark background
// immediately, before the external app.css link finishes loading. Without it a
// refresh flashes the browser's default white page (and the go-app loader over
// it) until the stylesheet applies.
const bootStyle = `<style>
  html, body { margin: 0; background-color: #1c1c1b; color: #f7f1ff; }
  #app-wasm-loader, .goapp-app-info { background-color: #1c1c1b; color: #8b888f; }
</style>`

// card strips (convenient when a battleline runs off-screen) and keeps the game
// log pinned to its newest entry.
const boardScript = `<script>
(function () {
  document.addEventListener('wheel', function (e) {
    var strip = e.target && e.target.closest ? e.target.closest('.card-strip') : null;
    if (strip && strip.scrollWidth > strip.clientWidth && e.deltaY !== 0) {
      strip.scrollLeft += e.deltaY;
      e.preventDefault();
    }
  }, { passive: false });
  // Keep the log pinned to the newest entry only while the player is already at
  // the bottom; if they scroll up, leave their position alone.
  var stick = true;
  document.addEventListener('scroll', function (e) {
    var log = e.target;
    if (!log || log.id !== 'gamelog') { return; }
    stick = (log.scrollTop + log.clientHeight) >= (log.scrollHeight - 4);
  }, true);
  new MutationObserver(function () {
    var log = document.getElementById('gamelog');
    if (log && stick) { log.scrollTop = log.scrollHeight; }
  }).observe(document.documentElement, { childList: true, subtree: true });

  // Drag hand cards onto the board: seed the drag (Firefox needs data on it) and
  // mark the board a valid drop target so the drop fires. The play logic runs in
  // Go via the card's OnDragStart and the board's OnDrop.
  document.addEventListener('dragstart', function (e) {
    var card = e.target && e.target.closest ? e.target.closest('.card') : null;
    if (card && card.getAttribute('draggable') === 'true' && e.dataTransfer) {
      e.dataTransfer.setData('text/plain', '');
      e.dataTransfer.effectAllowed = 'move';
    }
  });
  document.addEventListener('dragover', function (e) {
    if (!e.target || !e.target.closest) { return; }
    // Only the play area between the score pills is a drop target — not the score
    // pills or the hand, so releasing there does not count as playing the card.
    if (!e.target.closest('.play-zone')) { return; }
    e.preventDefault();
    if (e.dataTransfer) { e.dataTransfer.dropEffect = 'move'; }
  });
})();
</script>`

// devReloadScript polls the service worker for a new build. `mage web` restarts
// the server on every edit, which bumps go-app's version; the poll makes an open
// tab re-fetch app-worker.js, and go-app fires OnAppUpdate → ctx.Reload() so code
// and CSS changes appear without a manual refresh. Harmless in a static deploy —
// with a fixed version, update() finds nothing new.
const devReloadScript = `<script>
(function () {
  if (!('serviceWorker' in navigator)) { return; }
  setInterval(function () {
    navigator.serviceWorker.getRegistration().then(function (r) { if (r) { r.update(); } });
  }, 1500);
})();
</script>`

// styleEnv is the variable that turns the Style gallery on. mage web sets it, so
// the gallery is there whenever the client is being developed and absent from
// every other deployment.
const styleEnv = "VACTROL_STYLE"

// styleEnabled reports whether the Style gallery's route should exist. The check
// has to run on both sides of the build and agree: the server must register the
// route or it serves no page at all (go-app 404s an unregistered path), and the
// wasm client must register it or the served page renders nothing. app.Getenv
// bridges the two — it reads the process environment on the server and the Env
// map the server passed down on the client — so one variable decides both.
//
// It is an environment switch rather than a build tag because go-app routes on
// the client: a tag would have to exclude the gallery from the wasm bundle every
// player downloads, which means the gallery would not be compiled by default and
// would rot exactly as the //go:build todo card stubs do. See
// docs/adr/0014-style-gallery-on-real-components.md.
func styleEnabled() bool { return app.Getenv(styleEnv) == "1" }
