// Command web serves and runs the Vactrol browser client. The same binary is
// built two ways: compiled to WebAssembly it runs the interactive UI in the
// browser (app.RunWhenOnBrowser blocks there); built natively it serves that
// wasm bundle and the required go-app resources over HTTP.
package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/web"
)

func main() {
	app.Route("/", web.NewGame)
	app.RunWhenOnBrowser()

	version := resourceVersion()

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
			SVG:     "/web/assets/favicon.svg",
			Default: "/web/assets/favicon.svg",
		},
		// Plain CSS served as a static file from web/ — no CDN, no build step. The dev
		// server serves it from disk, so editing web/app.css and refreshing the browser
		// applies changes with the server left running (no restart).
		Styles:     []string{"/web/app.css"},
		RawHeaders: []string{bootStyle, boardScript, devReloadScript},
		// Byte size of app.wasm so go-app's loader accurately computes loading
		// progress even when serving gzip-compressed / chunked responses.
		WasmContentLength: wasmContentLength(),
		// Version keys go-app's service-worker cache. Left empty it defaults to the
		// app.wasm hash, so a CSS-only edit (wasm unchanged) would keep serving the
		// cached stylesheet. Hashing app.css too makes every asset edit bump it.
		Version: version,
		// The short build id the client shows, so a page and the server that built
		// it can be matched by eye.
		Env: map[string]string{"VACTROL_BUILD": buildID(version)},
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

// resourceVersion hashes the served static assets so the go-app Handler version
// changes whenever the wasm bundle or the stylesheet does, busting the service
// worker's precache on every edit.
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
	for _, p := range []string{"web/app.wasm", "web/app.css"} {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))[:12]
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
    { "src": "/web/assets/favicon.svg", "type": "image/svg+xml", "sizes": "any", "purpose": "any" }
  ]
}`

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
