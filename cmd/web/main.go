// Command web serves and runs the Vactrol browser client. The same binary is
// built two ways: compiled to WebAssembly it runs the interactive UI in the
// browser (app.RunWhenOnBrowser blocks there); built natively it serves that
// wasm bundle and the required go-app resources over HTTP.
package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/web"
)

func main() {
	app.Route("/", web.NewLanding)
	app.Route("/play", web.NewGame)
	app.Route("/cards", web.NewGallery)
	app.Route("/rulebook", web.NewRulebook)
	app.Route("/glossary", web.NewGlossary)
	if styleEnabled() {
		app.Route("/style", web.NewStyle)
		app.Route("/clusters", web.NewClusters)
	}
	app.RunWhenOnBrowser()

	version := resourceVersion()
	assets := assetPaths()

	// Serve a fullscreen web app manifest at go-app's manifest path. go-app
	// hardcodes display "standalone"; overriding the route makes an installed PWA
	// launch immersively, hiding the Android status and navigation bars.
	http.HandleFunc("/manifest.webmanifest", serveManifest)

	// Serve /web/ static assets ourselves, ahead of the go-app handler, so we can
	// stream prebuilt brotli/gzip files (mage webAssets) instead of compressing
	// every request. go-app would otherwise serve these itself and re-gzip the
	// multi-megabyte wasm on each fetch.
	http.Handle("/web/", staticAssets(version))

	http.Handle("/", gzipHandler(&app.Handler{
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
		Styles: []string{"/web/app.css"},
		RawHeaders: []string{
			bootStyle,
			appleTouchIcon,
			boardScript,
			galleryScript,
			cardFitScript,
			devReloadScript,
		},
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
	}))

	// Cloud Run injects PORT; fall back to 8000 for local `mage web`.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port
	log.Printf("Vactrol web client on http://localhost%s (build %s)", addr, buildID(version))
	if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}

// gzipHandler wraps h to gzip-encode responses for clients that accept it. It
// covers go-app's dynamically generated resources (the HTML shell, app.js,
// wasm_exec.js, app-worker.js); the large static assets under /web/ are served
// precompressed by staticAssets instead, so they never reach this path.
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

// staticAssets serves files under /web/, preferring the prebuilt brotli or gzip
// sibling (mage webAssets) that matches the client's Accept-Encoding and falling
// back to the raw file. It mirrors go-app's own no-cache + version ETag so an
// unchanged build still revalidates to 304, but adds brotli and skips the
// per-request compression go-app's file server would incur. The dev server has no
// precompressed siblings, so it simply serves the raw files from disk.
func staticAssets(version string) http.Handler {
	etag := `"` + version + `"`
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// path.Clean collapses any ".." before the web/ prefix guard below rejects
		// anything that escaped it, so traversal cannot reach outside web/.
		file := filepath.FromSlash(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
		if !strings.HasPrefix(file, "web"+string(filepath.Separator)) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Add("Vary", "Accept-Encoding")
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", assetContentType(file))
		accept := r.Header.Get("Accept-Encoding")
		if strings.Contains(accept, "br") && serveAsset(w, file+".br", "br") {
			return
		}
		if strings.Contains(accept, "gzip") && serveAsset(w, file+".gz", "gzip") {
			return
		}
		if !serveAsset(w, file, "") {
			http.NotFound(w, r)
		}
	})
}

// serveAsset streams path with the given Content-Encoding (empty for identity),
// reporting whether it existed. The caller has already set the Content-Type to
// the logical asset's type, since a .br/.gz name would resolve to the wrong one.
func serveAsset(w http.ResponseWriter, path, encoding string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		return false
	}
	if encoding != "" {
		w.Header().Set("Content-Encoding", encoding)
	}
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	_, _ = io.Copy(w, f)
	return true
}

// assetContentType maps an asset path to the Content-Type of its logical
// resource, independent of any .br/.gz compression suffix.
func assetContentType(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".css":
		return "text/css; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".wasm":
		return "application/wasm"
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".json", ".webmanifest":
		return "application/manifest+json"
	}
	if ct := mime.TypeByExtension(filepath.Ext(path)); ct != "" {
		return ct
	}
	return "application/octet-stream"
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
  /* go-app spins Icon.Default, which is the opaque PWA square. Swap in the bare
     Æmber gem (transparent) so only the gem spins, with no box corners to sweep
     the progress label; the PWA and favicon icons keep their opaque background. */
  #app-wasm-loader-icon { content: url("/web/assets/aember.svg"); }
</style>`

// card strips (convenient when a battleline runs off-screen) and keeps the game
// log pinned to its newest entry. A player bar (.score-pill) scrolls sideways the
// same way, so a squeezed bar's stats and keys are reached by wheeling over it.
const boardScript = `<script>
(function () {
  document.addEventListener('wheel', function (e) {
    var strip = e.target && e.target.closest
      ? e.target.closest('.card-strip, .score-pill') : null;
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

  // Publish the floating dock's live width as --dock-reserve so the lower player
  // bar ends at the dock's current edge, not at the widest a prompt could open it
  // to. A ResizeObserver tracks the dock as its content changes size; a
  // MutationObserver re-finds it whenever go-app re-renders it, and clears the
  // reserve when it is gone. Only a floating (fixed) dock reserves room — spanned
  // across the bottom in portrait it is static and the bar keeps the full width.
  var root = document.documentElement;
  function syncDockReserve() {
    var dock = document.querySelector('.control-dock--floating');
    var reserve = 0;
    if (dock && getComputedStyle(dock).position === 'fixed') {
      reserve = dock.getBoundingClientRect().width;
    }
    root.style.setProperty('--dock-reserve', reserve + 'px');
  }
  var dockSize = new ResizeObserver(syncDockReserve);
  new MutationObserver(function () {
    dockSize.disconnect();
    var dock = document.querySelector('.control-dock--floating');
    if (dock) { dockSize.observe(dock); }
    syncDockReserve();
  }).observe(document.documentElement, { childList: true, subtree: true });
  syncDockReserve();

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

// galleryScript makes the /cards filter sidebar resizable: dragging the
// .gallery-divider anywhere along its height sets --gallery-sidebar-w on the
// layout, which the sidebar's width reads. Pointer capture keeps the drag alive
// off the thin handle, and the width is clamped to the sidebar's min/max.
const galleryScript = `<script>
(function () {
  var W = '--gallery-sidebar-w', MIN = 176, MAX = 512;
  document.addEventListener('pointerdown', function (e) {
    var d = e.target && e.target.closest ? e.target.closest('.gallery-divider') : null;
    if (!d) { return; }
    var layout = d.closest('.gallery-layout');
    var side = layout && layout.querySelector('.gallery-sidebar');
    if (!side) { return; }
    e.preventDefault();
    d.classList.add('gallery-divider--drag');
    if (d.setPointerCapture) { d.setPointerCapture(e.pointerId); }
    var left = side.getBoundingClientRect().left;
    function move(ev) {
      var w = Math.max(MIN, Math.min(MAX, ev.clientX - left));
      layout.style.setProperty(W, w + 'px');
    }
    function up() {
      d.classList.remove('gallery-divider--drag');
      document.removeEventListener('pointermove', move);
      document.removeEventListener('pointerup', up);
    }
    document.addEventListener('pointermove', move);
    document.addEventListener('pointerup', up);
  });
})();
</script>`

// cardFitScript keeps a card's icon strip (.card-icons) and name banner
// (.card-name-text) on a single line: when either's content is wider than the
// space, it condenses it with scaleX instead of wrapping, clipping, or stepping a
// title to a preset size. The title is measured and fit here rather than estimated
// in Go, so it condenses to the exact rendered width regardless of which letters
// it uses. It re-fits on load, on resize, and whenever the DOM changes (go-app
// re-renders cards), coalescing bursts into one animation frame.
const cardFitScript = `<script>
(function () {
  function fitBand(band) {
    var inner = band.firstElementChild;
    if (!inner) { return; }
    inner.style.transform = '';
    var cs = getComputedStyle(band);
    var avail = band.clientWidth - parseFloat(cs.paddingLeft) - parseFloat(cs.paddingRight);
    var w = inner.scrollWidth;
    if (w > avail && w > 0) { inner.style.transform = 'scaleX(' + (avail / w) + ')'; }
  }
  function fitTitle(el) {
    // reset the prior fit so the measurement below reads the natural width.
    el.style.transform = ''; el.style.width = ''; el.style.flexShrink = '';
    var avail = el.clientWidth, natural = el.scrollWidth;
    if (avail <= 0 || natural <= avail) { return; }
    // grow the layout width by the inverse of the scale (and pin flex-shrink) so
    // the ellipsis does not fire on a title the scaleX already fits; a 0.7 floor
    // keeps a very long title legible and leaves the rest to the ellipsis.
    var scale = Math.max(0.7, avail / natural);
    el.style.transform = 'scaleX(' + scale + ')';
    el.style.width = (100 / scale) + '%';
    el.style.flexShrink = '0';
  }
  function fitFrame(card) {
    // offsetWidth is the layout width, unaffected by the lift's scale transform, so
    // the icons take their enlarged size on the first pass rather than the mid-grow
    // size getBoundingClientRect would read (which only corrects on the next render).
    var ow = card.offsetWidth;
    if (ow <= 0) { return; }
    // size the left-edge icons off the card's width so their scale is the same on a
    // board card and an enlarged copy (the fractions reproduce the CSS fallback rems
    // on a 9rem board card).
    card.style.setProperty('--edge-icon', (ow * 0.091).toFixed(2) + 'px');
    card.style.setProperty('--edge-house', (ow * 0.117).toFixed(2) + 'px');
    // the two-tone frame splits at fixed card percentages by default; measure the
    // text box and rewrite the split inline so the right half steps down to the
    // box's top and the left half to its bottom, bracketing the box no matter how
    // tall the card is drawn. The rects are scale-invariant ratios, so the split is
    // right even mid-grow.
    var text = card.querySelector('.card-text');
    var cr = card.getBoundingClientRect();
    if (!text || cr.height <= 0) { return; }
    var tr = text.getBoundingClientRect();
    // nudge each split a corner-radius (~6px) past the box edge so the step clears
    // the card's rounded corner: the right half a touch lower, the left a touch higher.
    var pad = 6 / cr.height * 100;
    card.style.setProperty('--nm-right', ((tr.top - cr.top) / cr.height * 100 + pad).toFixed(2) + '%');
    card.style.setProperty('--nm-left', ((tr.bottom - cr.top) / cr.height * 100 - pad).toFixed(2) + '%');
  }
  var scheduled = false;
  function schedule() {
    if (scheduled) { return; }
    scheduled = true;
    requestAnimationFrame(function () {
      scheduled = false;
      document.querySelectorAll('.card-icons').forEach(fitBand);
      document.querySelectorAll('.card-name-text').forEach(fitTitle);
      document.querySelectorAll('.card').forEach(fitFrame);
    });
  }
  window.addEventListener('load', schedule);
  window.addEventListener('resize', schedule);
  // childList only, so our own style writes do not retrigger the observer.
  new MutationObserver(schedule).observe(document.documentElement, { childList: true, subtree: true });
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
  // When a freshly built worker takes control, reload once so the new wasm/CSS is
  // actually shown — this is what lets an already-open tab reset itself instead of
  // waiting for a manual hard refresh.
  var reloaded = false;
  navigator.serviceWorker.addEventListener('controllerchange', function () {
    if (reloaded) { return; }
    reloaded = true;
    window.location.reload();
  });
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
