package web

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the Style gallery: one page showing every piece of the client's
// visual vocabulary at once. See docs/adr/0014-style-gallery-on-real-components.md
// for why it renders the game's own components instead of its own markup, and
// why it ships in the production bundle rather than hiding behind a build tag.

// stylePersistKey is where the gallery's font choices are kept. It is separate
// from the match's key so clearing a stuck game never costs the fonts, and
// choosing fonts never touches a game in progress.
const stylePersistKey = "vactrol.style"

// fontOrigins are the stylesheet hosts the gallery will load a font from. The
// gallery injects whatever URL is pasted into it as a <link> in the page, so the
// list is a security boundary and not a convenience: without it, any pasted URL
// becomes a stylesheet this origin executes with, which is a stored-XSS shaped
// hole in a page that otherwise has no user input at all.
var fontOrigins = map[string]bool{
	"fonts.googleapis.com": true,
	"fonts.bunny.net":      true,
	"cdn.jsdelivr.net":     true,
}

// NewStyle returns the root component for the Style gallery. It is registered
// only when the page is served from localhost — see cmd/web/main.go.
func NewStyle() app.Composer { return &style{} }

// style is the gallery's root component. It owns a synthetic game to render the
// board surfaces from, and the font selections the page is compared under.
type style struct {
	app.Compo

	// harness is a game built to a fixed display position, so the surfaces that
	// only exist mid-match (the Player bar) have something to draw.
	harness *game

	// fonts are the font families loaded so far, in the order they were added.
	// The first is the page default; the compare strip renders one specimen per
	// entry.
	fonts []string
	// fontURL is the stylesheet URL box's contents, and fontErr the reason the
	// last one was refused ("" when the last add succeeded).
	fontURL string
	fontErr string
	// uiFont is the family the whole page is set in; houseFont overrides it per
	// house, which is the comparison the gallery exists to make. An empty entry
	// means "inherit the UI font".
	uiFont    string
	houseFont [engine.NumHouses]string

	// flashOdd flips on every replay so the animation classes alternate, which is
	// how the client restarts a CSS animation that is already running.
	flashOdd bool
	// autoloop replays the animations continuously rather than on each press.
	autoloop bool
}

// styleSection is one titled block of the page, named so the section links and
// the sections themselves are built from one list and cannot drift apart.
type styleSection struct {
	id, title string
	body      func() app.UI
}

func (s *style) OnMount(ctx app.Context) {
	s.harness = styleHarness()
	var saved struct {
		Fonts  []string
		UIFont string
		House  [engine.NumHouses]string
	}
	_ = ctx.LocalStorage().Get(stylePersistKey, &saved)
	s.fonts, s.uiFont, s.houseFont = saved.Fonts, saved.UIFont, saved.House
	for _, f := range s.fonts {
		linkFont(f)
	}
}

// save writes the font selections back to local storage. Nothing else on the
// page is worth persisting: the specimens are derived from the catalog and the
// animations are momentary.
func (s *style) save(ctx app.Context) {
	_ = ctx.LocalStorage().Set(stylePersistKey, struct {
		Fonts  []string
		UIFont string
		House  [engine.NumHouses]string
	}{s.fonts, s.uiFont, s.houseFont})
}

// styleHarness builds the game the board surfaces are drawn from. It is a real
// engine.Game filled through the engine's own setup calls, then given fixed
// display values for chains, keys and Æmber — no single reachable position shows
// every branch of a Player bar at once, and a gallery that has to be played to
// is a gallery nobody looks at.
func styleHarness() *game {
	g := &game{selHand: -1, zonesPlayer: -1, forgingKey: -1, handSlot: -1}
	g.g = engine.NewGame("Player One", "Player Two", 1)
	g.mavericks = map[engine.LocalID]bool{}
	g.deckHouses = [2][]engine.House{
		{engine.Brobnar, engine.Dis, engine.Logos},
		{engine.Mars, engine.Sanctum, engine.Untamed},
	}

	// Fill the zones so the counts in the Player bar are all different and none
	// is zero: equal counts hide a wiring mistake that swaps two of them.
	all := cards.All()
	sizes := []struct {
		n   int
		add func(engine.CardDefinition, int) engine.LocalID
	}{
		{6, g.g.AddToHand},
		{22, g.g.AddToDeck},
		{9, g.g.AddToDiscard},
		{2, g.g.AddToArchives},
	}
	next := 0
	for player := range 2 {
		for _, z := range sizes {
			for range z.n {
				z.add(all[next%len(all)], player)
				next++
			}
		}
	}

	// Display values, set outright rather than played to.
	g.g.State.ActivePlayer = 0
	g.g.State.ActiveHouse = engine.Brobnar
	g.g.State.Aember = [2]int{7, 3}
	g.g.State.Chains = [2]int{0, 3}
	return g
}

// Render lays the gallery out as one scrolling page under a sticky header: the
// primitives first (colours, icons, type), then what they compose into (cards,
// the Player bar), then what moves. Reading top to bottom is reading the design
// system bottom up.
func (s *style) Render() app.UI {
	if s.harness == nil {
		return app.Div().Class("style-page").Text("Building the gallery…")
	}
	sections := []styleSection{
		{"colors", "Colour tokens", s.colorSection},
		{"icons", "Icons", s.iconSection},
		{"type", "Typography", s.typeSection},
		{"houses", "House grid", s.houseSection},
		{"features", "Card features", s.featureSection},
		{"bar", "Player bar", s.barSection},
		{"motion", "Animations", s.motionSection},
	}
	body := []app.UI{s.header(sections)}
	for _, sec := range sections {
		body = append(body, app.Section().ID(sec.id).Class("style-section").Body(
			app.H2().Class("style-h2").Text(sec.title),
			sec.body(),
		))
	}
	page := app.Div().Class("style-page").Body(body...)
	if s.uiFont != "" {
		page = page.Style("font-family", quoteFamily(s.uiFont))
	}
	return page
}

// header is the sticky strip carrying the section links and the font controls.
func (s *style) header(sections []styleSection) app.UI {
	links := make([]app.UI, 0, len(sections))
	for _, sec := range sections {
		links = append(links, app.A().Class("style-nav-link").Href("#"+sec.id).Text(sec.title))
	}
	return app.Header().Class("style-header").Body(
		app.Nav().Class("style-nav").Body(links...),
		s.fontControls(),
	)
}

// fontControls is the font loader and the two selectors it feeds: the page-wide
// UI font, and a per-house override. Trying a face on one house at a time is the
// question the gallery was built to answer, so it lives in the header where it
// stays reachable from every section.
func (s *style) fontControls() app.UI {
	body := []app.UI{
		app.Input().Class("style-input").
			Type("url").
			Placeholder("Font stylesheet URL (Google Fonts, Bunny, jsDelivr)").
			Value(s.fontURL).
			OnChange(s.onFontURL),
		app.Button().Class("style-btn").Text("Load font").OnClick(s.onAddFont),
	}
	if s.fontErr != "" {
		body = append(body, app.Span().Class("style-err").Text(s.fontErr))
	}
	if len(s.fonts) > 0 {
		body = append(body,
			app.Label().Class("style-label").Body(
				app.Text("UI font"),
				s.fontSelect(s.uiFont, s.onUIFont),
			),
			app.Label().Class("style-label").Body(
				app.Text("House font"),
				s.houseFontSelect(),
			),
		)
	}
	return app.Div().Class("style-controls").Body(body...)
}

// fontSelect renders the loaded families as a picker, with an empty first option
// meaning "leave it to the stylesheet".
func (s *style) fontSelect(current string, on func(app.Context, app.Event)) app.UI {
	opts := []app.UI{app.Option().Value("").Text("(default)").Selected(current == "")}
	for _, f := range s.fonts {
		opts = append(opts, app.Option().Value(f).Text(f).Selected(f == current))
	}
	return app.Select().Class("style-select").OnChange(on).Body(opts...)
}

// houseFontSelect pairs a house picker with a family picker. The house being
// edited is read back from the select rather than held in state, so the control
// has no mode of its own to get stuck in.
func (s *style) houseFontSelect() app.UI {
	rows := make([]app.UI, 0, engine.NumHouses)
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		rows = append(rows, app.Div().Class(cx("style-house-font", houseAccent(h))).Body(
			houseIcon(h, "icon-inline"),
			app.Span().Text(h.String()),
			s.fontSelect(s.houseFont[h], s.onHouseFont(h)),
		))
	}
	return app.Div().Class("style-house-fonts").Body(rows...)
}

func (s *style) onFontURL(ctx app.Context, _ app.Event) {
	s.fontURL = ctx.JSSrc().Get("value").String()
}

func (s *style) onUIFont(ctx app.Context, _ app.Event) {
	s.uiFont = ctx.JSSrc().Get("value").String()
	s.save(ctx)
}

// onHouseFont returns the change handler for one house's font override. The
// house is bound here rather than read from the event because the select only
// carries the family.
func (s *style) onHouseFont(h engine.House) func(app.Context, app.Event) {
	return func(ctx app.Context, _ app.Event) {
		s.houseFont[h] = ctx.JSSrc().Get("value").String()
		s.save(ctx)
	}
}

// onAddFont loads the pasted stylesheet and remembers the family it names.
func (s *style) onAddFont(ctx app.Context, _ app.Event) {
	family, err := fontFamily(s.fontURL)
	if err != nil {
		s.fontErr = err.Error()
		return
	}
	s.fontErr = ""
	for _, f := range s.fonts {
		if f == family {
			return
		}
	}
	s.fonts = append(s.fonts, family)
	linkFont(family)
	s.fontURL = ""
	s.save(ctx)
}

// linkFont appends a stylesheet <link> for a family, using the provider's own
// CSS API rather than the pasted URL. Rebuilding the URL from the allowlisted
// host and the extracted family means the string that reaches the document is
// one this code composed, not one the box was handed.
func linkFont(family string) {
	if family == "" {
		return
	}
	href := "https://fonts.googleapis.com/css2?family=" + familyParam(family) + "&display=swap"
	doc := app.Window().Get("document")
	link := doc.Call("createElement", "link")
	link.Set("rel", "stylesheet")
	link.Set("href", href)
	doc.Get("head").Call("appendChild", link)
}

// familyParam encodes a family name for a font provider's family= parameter. The
// separator between words is a literal "+", which is why the whole name cannot
// just be escaped: escaping turns the separator into %2B and the provider
// answers 400. Each word is escaped on its own and the separators are added
// after, so a name still cannot smuggle anything into the URL.
func familyParam(family string) string {
	words := strings.Fields(family)
	for i, w := range words {
		words[i] = url.QueryEscape(w)
	}
	return strings.Join(words, "+")
}

// fontFamily extracts the family a font stylesheet URL names, rejecting any URL
// that is not an https stylesheet from an allowlisted provider. The family is
// what the page actually needs — the URL is only how the browser is told to
// fetch it — so pulling the family out here is what lets the fetch be rebuilt
// from known-good parts.
func fontFamily(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", errors.New("that is not a URL")
	}
	if u.Scheme != "https" {
		return "", errors.New("the stylesheet URL must be https")
	}
	if !fontOrigins[u.Host] {
		return "", errors.New(u.Host + " is not an allowed font provider")
	}
	family := rawParam(u.RawQuery, "family")
	// Google's css2 API appends axes after a colon (Inter:wght@400;700).
	family, _, _ = strings.Cut(family, ":")
	family = strings.TrimSpace(strings.ReplaceAll(family, "+", " "))
	if family == "" {
		return "", errors.New("the URL names no family= to load")
	}
	return family, nil
}

// rawParam reads one query parameter out of a raw query string. url.Values is no
// use here: Go rejects a query containing a semicolon, and a semicolon is
// exactly what Google's css2 API puts between font weights
// (family=Inter:wght@400;700), so the commonest font URL there is would come
// back with no family at all.
func rawParam(query, name string) string {
	for _, pair := range strings.Split(query, "&") {
		k, v, ok := strings.Cut(pair, "=")
		if !ok || k != name {
			continue
		}
		if unescaped, err := url.QueryUnescape(v); err == nil {
			return unescaped
		}
		return v
	}
	return ""
}

// quoteFamily wraps a family name for a CSS font-family value and appends the
// page's own fallback, so a font that fails to load leaves the page readable
// rather than falling back to the browser default.
func quoteFamily(family string) string {
	return `"` + strings.ReplaceAll(family, `"`, "") + `", system-ui, sans-serif`
}

// styleCard renders one specimen as a captioned card face. A gap — a query that
// matched nothing — is drawn as an empty slot saying so, because which
// combinations the loaded sets do not cover is information the grid exists to
// show.
func (s *style) styleCard(sp specimen) app.UI {
	caption := app.Div().Class("style-caption").Body(
		app.Span().Class("style-caption-q").Text(sp.Caption),
		app.Span().Class("style-caption-a").Text(specimenLabel(sp)),
	)
	if !sp.found() {
		return app.Div().Class("style-specimen style-specimen--gap").Body(
			app.Div().Class("style-gap").Text("—"),
			caption,
		)
	}
	face := printedFace(sp.Def)
	face.Bar = barKeywordsOf(sp.Def)
	body := app.Div().Class("style-specimen").Body(face, caption)
	if f := s.houseFont[sp.Def.House]; f != "" {
		body = body.Style("--font-house", quoteFamily(f))
	}
	return body
}

// specimenRow lays specimens out as a wrapping strip.
func (s *style) specimenRow(specs []specimen) app.UI {
	out := make([]app.UI, 0, len(specs))
	for _, sp := range specs {
		out = append(out, s.styleCard(sp))
	}
	return app.Div().Class("style-row").Body(out...)
}

// colorSection shows the palette every surface is built from: the page-wide
// tokens, then each house's four colours. The house swatches are labelled with
// the token name rather than the hex, because the token is what code names.
func (s *style) colorSection() app.UI {
	var body []app.UI
	rootRow := make([]app.UI, 0, len(styleTokens))
	for _, t := range styleTokens {
		rootRow = append(rootRow, app.Div().Class("style-swatch").Body(
			app.Div().Class("style-chip").Style("background", "var("+t+")"),
			app.Span().Class("style-mono").Text(t),
		))
	}
	body = append(body, app.Div().Class("style-row").Body(rootRow...))
	for h := engine.HouseNone; int(h) < engine.NumHouses; h++ {
		chips := make([]app.UI, 0, len(houseTokens))
		for _, t := range houseTokens {
			chips = append(chips, app.Div().Class("style-swatch").Body(
				app.Div().Class("style-chip").Style("background", "var("+t+")"),
				app.Span().Class("style-mono").Text(t),
			))
		}
		body = append(body, app.Div().Class(cx("style-palette", houseClasses(h))).Body(
			app.Span().Class("style-palette-name").Body(
				houseIcon(h, "icon-inline"), app.Text(h.String()),
			),
			app.Div().Class("style-row").Body(chips...),
		))
	}
	return app.Div().Body(body...)
}

// styleTokens are the page-wide colour custom properties, and houseTokens the
// four a house palette redefines. They are listed here so the gallery can show
// them; the tests hold the list to what web/app.css actually defines.
var styleTokens = []string{
	"--bg", "--surface", "--surface-2", "--line", "--line-2",
	"--ink", "--ink-2", "--p0", "--p1", "--warning", "--danger",
}

var houseTokens = []string{"--nm", "--nm-fg", "--tp", "--tp-fg", "--edge"}

// iconSection shows every icon asset at the size the board draws it. The stems
// come from galleryIcons, which a test holds equal to web/assets, so an icon
// added to the app cannot quietly stay out of the gallery.
func (s *style) iconSection() app.UI {
	out := make([]app.UI, 0, len(galleryIcons))
	for _, name := range galleryIcons {
		out = append(out, app.Div().Class("style-swatch").Body(
			icon(name, "icon-outline"),
			app.Span().Class("style-mono").Text(name),
		))
	}
	return app.Div().Class("style-row style-row--icons").Body(out...)
}

// typeSection is the font comparison: the same specimen rendered once per loaded
// family, side by side. Everywhere else on the page shows the active selection;
// this is the only place that shows the candidates together, which is what
// choosing between two faces actually requires.
func (s *style) typeSection() app.UI {
	sp := firstMatch("Longest rules text", longestRules())
	families := append([]string{""}, s.fonts...)
	out := make([]app.UI, 0, len(families))
	for _, f := range families {
		label := f
		if label == "" {
			label = "(stylesheet default)"
		}
		col := app.Div().Class("style-specimen").Body(
			printedFace(sp.Def),
			app.Div().Class("style-caption").Body(
				app.Span().Class("style-caption-q").Text(label),
			),
		)
		if f != "" {
			col = col.Style("font-family", quoteFamily(f))
		}
		out = append(out, col)
	}
	return app.Div().Class("style-row").Body(out...)
}

// houseSection is the House-by-Card-type table. Empty cells are left in place
// rather than collapsed away.
func (s *style) houseSection() app.UI {
	houses, types, cells := houseGrid()
	head := []app.UI{app.Div().Class("style-grid-corner")}
	for _, ct := range types {
		head = append(head, app.Div().Class("style-grid-head").Body(
			icon(typeIconName(ct), "icon-inline"), app.Text(ct.String()),
		))
	}
	rows := []app.UI{app.Div().Class("style-grid-row").Body(head...)}
	for i, h := range houses {
		row := []app.UI{app.Div().Class(cx("style-grid-head", houseAccent(h))).Body(
			houseIcon(h, "icon-inline"), app.Text(h.String()),
		)}
		for j := range types {
			row = append(row, s.styleCard(cells[i*len(types)+j]))
		}
		rows = append(rows, app.Div().Class("style-grid-row").Body(row...))
	}
	return app.Div().Class("style-grid").Body(rows...)
}

// featureSection shows one card per face feature and per rarity mark.
func (s *style) featureSection() app.UI {
	return app.Div().Body(
		app.H3().Class("style-h3").Text("Rarity marks"),
		s.specimenRow(raritySpecimens()),
		app.H3().Class("style-h3").Text("Keywords and stats"),
		s.specimenRow(featureSpecimens()),
	)
}

// barSection renders both Player bars from the harness, which is the one surface
// the gallery cannot build from a definition alone.
func (s *style) barSection() app.UI {
	return app.Div().Class("style-bars").Body(
		s.harness.scorePill(0),
		s.harness.scorePill(1),
	)
}

// motionSection replays the card animations on demand. Each is the real card
// face with its animation flag set, so what plays here is what plays on the
// board; the replay buttons flip the parity class, which is how the client
// restarts an animation that is already running.
func (s *style) motionSection() app.UI {
	sp := firstMatch("Creature", func(d *engine.CardDefinition) bool {
		return d.Type == engine.Creature
	})
	kinds := []struct {
		name string
		set  func(*cardView)
	}{
		{"Enter play", func(c *cardView) { c.Enter = true }},
		{"Fight (up)", func(c *cardView) { c.Fight = true }},
		{"Fight (down)", func(c *cardView) { c.Fight, c.FightDown = true, true }},
		{"Take damage", func(c *cardView) { c.Hit = true }},
		{"Stun", func(c *cardView) { c.Stunned, c.StunFlash = true, true }},
		{"Exhaust", func(c *cardView) { c.Exhausted, c.ExhaustFlash = true, true }},
	}
	out := make([]app.UI, 0, len(kinds))
	for _, k := range kinds {
		face := printedFace(sp.Def)
		face.FlashOdd = s.flashOdd
		k.set(face)
		out = append(out, app.Div().Class("style-specimen").Body(
			face,
			app.Div().Class("style-caption").Body(
				app.Span().Class("style-caption-q").Text(k.name),
			),
		))
	}
	return app.Div().Body(
		app.Div().Class("style-controls").Body(
			app.Button().Class("style-btn").Text("Replay all").OnClick(s.onReplay),
			app.Label().Class("style-label").Body(
				app.Input().Type("checkbox").Checked(s.autoloop).OnChange(s.onAutoloop),
				app.Text("Loop"),
			),
		),
		app.Div().Class("style-row").Body(out...),
	)
}

func (s *style) onReplay(_ app.Context, _ app.Event) { s.flashOdd = !s.flashOdd }

// onAutoloop starts or stops the continuous replay. The loop is a chain of
// deferred replays rather than a ticker so it stops as soon as the box is
// cleared, without a goroutine outliving the page.
func (s *style) onAutoloop(ctx app.Context, _ app.Event) {
	s.autoloop = ctx.JSSrc().Get("checked").Bool()
	if s.autoloop {
		s.loop(ctx)
	}
}

func (s *style) loop(ctx app.Context) {
	if !s.autoloop {
		return
	}
	s.flashOdd = !s.flashOdd
	ctx.After(styleLoopPeriod, func(ctx app.Context) { s.loop(ctx) })
}

// styleLoopPeriod is how long a looped animation is left to finish before it is
// restarted. The longest card animation is well under half a second.
const styleLoopPeriod = 900 * time.Millisecond

// galleryIcons lists every icon asset the gallery shows, which is every asset in
// web/assets except the favicon (a browser chrome image, not part of the board's
// vocabulary). A test holds this equal to the directory.
var galleryIcons = []string{
	"aember", "chains", "damage", "exhausted", "forge",
	"house-brobnar", "house-dis", "house-logos", "house-mars",
	"house-sanctum", "house-shadows", "house-untamed",
	"key", "key-blue", "key-red", "key-yellow",
	"maverick", "power", "power-counter-minus", "power-counter-plus",
	"rarity-diamond", "redo", "restart", "shield", "stun",
	"type-action", "type-artifact", "type-creature", "type-upgrade",
	"undo", "wrench",
	"zone-archives", "zone-deck", "zone-discard", "zone-hand", "zone-purge",
}
