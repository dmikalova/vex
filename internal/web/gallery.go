package web

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/card"
	"github.com/dmikalova/vex/internal/deckgen"
	"github.com/dmikalova/vex/internal/engine"
)

// This file is the /cards gallery (ADR 0023): every card in the database as a
// printed face, filterable by a named sidebar of facets (House, Set, Type,
// Rarity, Keyword, Trait), by name and rules-text queries (see gallery_query.go
// for the query grammar), and by power/armor ranges, then ordered by a chosen
// key. Facets are OR within a category and AND across categories. Template cards
// (Master of X) are shown as their materialized variants (Master of 1…5), not the
// placeholder face. Only the cards on or near the viewport are mounted as full
// printed faces; the rest render as fixed-size text placeholders (see window and
// galleryPlaceholder), so a large catalog keeps a correct scroll height and stays
// find-in-page searchable without paying to lay out and paint every face.

// NewGallery returns the root component for the /cards page.
func NewGallery() app.Composer { return &gallery{} }

// galleryCard is one catalog entry: a card's concrete definition, the sets it
// belongs to (native and, separately, by reprint), and pre-lowered name and rules
// haystacks the queries match against so filtering never re-renders per keystroke.
type galleryCard struct {
	def         *engine.CardDefinition
	nativeSets  map[string]bool
	reprintSets map[string]bool
	nameHay     string
	textHay     string
}

// gallery is the /cards page component. It builds the catalog once on mount, then
// every render filters and orders it by the active facets, queries, and ranges.
type gallery struct {
	app.Compo

	cards    []galleryCard
	houses   []engine.House
	sets     []string
	types    []engine.CardType
	rarities []engine.Rarity
	keywords []engine.Keyword
	traits   []engine.Trait

	selHouse   map[engine.House]bool
	selSet     map[string]bool
	selType    map[engine.CardType]bool
	selRarity  map[engine.Rarity]bool
	selKeyword map[engine.Keyword]bool
	selTrait   map[engine.Trait]bool

	excludeReprints bool
	name            string
	text            string
	minPower        string
	maxPower        string
	minArmor        string
	maxArmor        string
	order           string

	// ready is false until the async catalog build finishes; the page renders its
	// shell immediately and fills the grid in once ready.
	ready bool

	// winFrom..winTo is the half-open index range of the shown cards drawn as full
	// printed faces; every other shown card is a lightweight text placeholder (see
	// window and galleryPlaceholder). A scroll or resize re-measures the viewport
	// and slides the window (recomputeWindow), so only the cards on or near screen
	// pay for a heavy face while the rest keep the scroll height and find-in-page
	// text intact.
	winFrom int
	winTo   int
	// dispatch schedules a re-render on the UI goroutine (captured in OnMount), so a
	// scroll or resize callback can slide the window and repaint.
	dispatch func(func(app.Context))
	// scrollFunc is the document scroll listener that slides the window; measureFunc
	// is the requestAnimationFrame callback a render queues to re-measure after the
	// grid relays out (e.g. a filter change). measureQueued guards against queuing
	// more than one frame at a time.
	scrollFunc    app.Func
	measureFunc   app.Func
	measureQueued bool
}

// galleryInitialWindow is how many faces the grid draws before the first
// measurement, so the initial paint mounts a screenful of real faces rather than
// the whole catalog, then recomputeWindow narrows it to what is actually visible.
const galleryInitialWindow = 60

// galleryWindowBuffer is how many rows of cards beyond the viewport, on each side,
// stay mounted as full faces, so a short scroll reveals painted cards rather than
// bare placeholders while the next frame catches up.
const galleryWindowBuffer = 4

// rarityOrder gives each rarity its catalog rank, so the Rarity facet reads in
// rulebook order rather than alphabetically.
var rarityOrder = map[engine.Rarity]int{
	engine.Common:    0,
	engine.Uncommon:  1,
	engine.Rare:      2,
	engine.Special:   3,
	engine.Connected: 4,
}

// OnMount initializes the empty filter state, then builds the catalog off the
// initial paint (ctx.Async) so the page shell appears immediately and the cards
// stream in when the build — materializing template cards, indexing set
// memberships and search haystacks — finishes.
func (g *gallery) OnMount(ctx app.Context) {
	g.selHouse = map[engine.House]bool{}
	g.selSet = map[string]bool{}
	g.selType = map[engine.CardType]bool{}
	g.selRarity = map[engine.Rarity]bool{}
	g.selKeyword = map[engine.Keyword]bool{}
	g.selTrait = map[engine.Trait]bool{}
	g.order = "name"
	g.winFrom = 0
	g.winTo = galleryInitialWindow
	g.dispatch = ctx.Dispatch
	g.installWindowing()

	ctx.Async(func() {
		cat := buildGalleryCatalog()
		ctx.Dispatch(func(app.Context) {
			g.cards = cat.cards
			g.houses = cat.houses
			g.sets = cat.sets
			g.types = cat.types
			g.rarities = cat.rarities
			g.keywords = cat.keywords
			g.traits = cat.traits
			g.ready = true
		})
	})
}

// installWindowing wires the scroll listener and the frame-measure callback that
// keep the drawn window of full faces aligned with the viewport. The page scrolls
// on the document, and a scroll event does not bubble, so the listener is
// registered in the capture phase. It no-ops off-browser, where there is no page
// to measure and the grid draws its initial window of faces unchanged.
func (g *gallery) installWindowing() {
	if !app.IsClient || g.scrollFunc != nil {
		return
	}
	g.scrollFunc = app.FuncOf(func(app.Value, []app.Value) any {
		g.recomputeWindow()
		return nil
	})
	g.measureFunc = app.FuncOf(func(app.Value, []app.Value) any {
		g.measureQueued = false
		g.recomputeWindow()
		return nil
	})
	app.Window().Get("document").Call("addEventListener", "scroll", g.scrollFunc, true)
}

// OnDismount releases the windowing listeners so a navigation away from the
// gallery leaves nothing bound to the document.
func (g *gallery) OnDismount() {
	if g.scrollFunc != nil {
		app.Window().Get("document").
			Call("removeEventListener", "scroll", g.scrollFunc, true)
		g.scrollFunc.Release()
		g.scrollFunc = nil
	}
	if g.measureFunc != nil {
		g.measureFunc.Release()
		g.measureFunc = nil
	}
}

// OnResize re-measures the window when the viewport changes, since the number of
// columns and the visible row count both follow the width and height.
func (g *gallery) OnResize(app.Context) { g.recomputeWindow() }

// scheduleMeasure asks for one re-measure on the next frame, after the grid a
// render just produced has laid out. It is how a filter change — which reflows the
// grid without a scroll or resize — re-aligns the window. It queues at most one
// frame at a time and no-ops off-browser.
func (g *gallery) scheduleMeasure() {
	if !app.IsClient || g.measureFunc == nil || g.measureQueued {
		return
	}
	g.measureQueued = true
	app.Window().Call("requestAnimationFrame", g.measureFunc)
}

// recomputeWindow measures the viewport against the laid-out grid and, when the
// visible index range has changed, slides the window and repaints. It runs on the
// UI goroutine (a scroll, resize, or frame callback), so it reads and writes the
// window fields directly and only repaints on a real change, which keeps a scroll
// from looping through renders that do not move the window.
func (g *gallery) recomputeWindow() {
	from, to, ok := g.measureWindow()
	if !ok || (from == g.winFrom && to == g.winTo) {
		return
	}
	g.winFrom, g.winTo = from, to
	if g.dispatch != nil {
		g.dispatch(func(app.Context) {})
	}
}

// measureWindow reads the laid-out grid and viewport and returns the half-open
// index range of cards on or near screen, given a uniform card box. It reports
// ok=false when there is nothing to measure (off-browser, or before the grid has
// mounted), so the caller leaves the current window in place.
func (g *gallery) measureWindow() (from, to int, ok bool) {
	if !app.IsClient {
		return 0, 0, false
	}
	win := app.Window()
	doc := win.Get("document")
	if !doc.Truthy() {
		return 0, 0, false
	}
	grid := doc.Call("querySelector", ".gallery-grid")
	if !grid.Truthy() {
		return 0, 0, false
	}
	slots := grid.Get("children")
	if slots.Get("length").Int() == 0 {
		return 0, 0, false
	}
	box := slots.Call("item", 0).Call("getBoundingClientRect")
	cardW := box.Get("width").Float()
	cardH := box.Get("height").Float()
	if cardW <= 0 || cardH <= 0 {
		return 0, 0, false
	}
	gridRect := grid.Call("getBoundingClientRect")
	style := win.Call("getComputedStyle", grid)
	rowGap := parsePx(style.Get("rowGap").String())
	colGap := parsePx(style.Get("columnGap").String())
	cols := max(int((gridRect.Get("width").Float()+colGap)/(cardW+colGap)), 1)
	rowH := cardH + rowGap
	if rowH <= 0 {
		return 0, 0, false
	}
	gridTop := gridRect.Get("top").Float()
	innerH := win.Get("innerHeight").Float()
	firstRow := int(-gridTop/rowH) - galleryWindowBuffer
	lastRow := int((innerH-gridTop)/rowH) + galleryWindowBuffer
	if firstRow < 0 {
		firstRow = 0
	}
	if lastRow < firstRow {
		lastRow = firstRow
	}
	return firstRow * cols, (lastRow + 1) * cols, true
}

// window clamps the drawn-face index range to the shown-card count, so a filter
// that shrinks the list can never point the window past its end.
func (g *gallery) window(n int) (from, to int) {
	from, to = g.winFrom, g.winTo
	if from < 0 {
		from = 0
	}
	if from > n {
		from = n
	}
	if to > n {
		to = n
	}
	if to < from {
		to = from
	}
	return from, to
}

// parsePx reads a CSS pixel length like "12px" as a float, treating an empty or
// unparsable value as zero.
func parsePx(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(s), "px"), 64)
	return f
}

// galleryCatalog is the built catalog: every card as a printed entry plus the
// distinct facet values the sidebar offers, all ordered for display.
type galleryCatalog struct {
	cards    []galleryCard
	houses   []engine.House
	sets     []string
	types    []engine.CardType
	rarities []engine.Rarity
	keywords []engine.Keyword
	traits   []engine.Trait
}

// buildGalleryCatalog materializes template cards into their concrete variants,
// indexes each card's set memberships and search haystacks, and collects the
// distinct facet values the sidebar offers. It is pure, so OnMount runs it off
// the UI goroutine.
func buildGalleryCatalog() galleryCatalog {
	var cat galleryCatalog
	reprintsByName := reprintSetsByName()
	seenHouse := map[engine.House]bool{}
	seenSet := map[string]bool{}
	seenType := map[engine.CardType]bool{}
	seenRarity := map[engine.Rarity]bool{}
	seenKeyword := map[engine.Keyword]bool{}
	seenTrait := map[engine.Trait]bool{}

	regs := card.Cards()
	for i := range regs {
		native := nativeSetsOf(regs[i])
		reprint := reprintsByName[regs[i].Def.Name]
		defs := materializedDefs(regs[i])
		for j := range defs {
			def := defs[j]
			cat.cards = append(cat.cards, galleryCard{
				def:         &def,
				nativeSets:  native,
				reprintSets: reprint,
				nameHay:     def.Name,
				textHay:     def.Name + "\n" + engine.RenderCardRules(&def),
			})
			markSeen(seenHouse, &cat.houses, def.House)
			markSeen(seenType, &cat.types, def.Type)
			markSeen(seenRarity, &cat.rarities, def.Rarity)
			for _, k := range def.Keywords {
				markSeen(seenKeyword, &cat.keywords, k)
			}
			for _, t := range def.Traits {
				markSeen(seenTrait, &cat.traits, t)
			}
		}
		for s := range native {
			markSeen(seenSet, &cat.sets, s)
		}
		for s := range reprint {
			markSeen(seenSet, &cat.sets, s)
		}
	}

	sort.Slice(
		cat.houses,
		func(i, j int) bool { return cat.houses[i].String() < cat.houses[j].String() },
	)
	sort.Strings(cat.sets)
	sort.Slice(
		cat.rarities,
		func(i, j int) bool { return rarityOrder[cat.rarities[i]] < rarityOrder[cat.rarities[j]] },
	)
	sort.Slice(
		cat.keywords,
		func(i, j int) bool { return cat.keywords[i].String() < cat.keywords[j].String() },
	)
	sort.Slice(
		cat.traits,
		func(i, j int) bool { return cat.traits[i].String() < cat.traits[j].String() },
	)
	return cat
}

// OnAppUpdate reloads the page onto a freshly built wasm bundle, so a dev edit
// appears without a manual refresh (the same hand-off the game page uses).
func (g *gallery) OnAppUpdate(ctx app.Context) { ctx.Reload() }

// markSeen adds v to the ordered slice the first time it is seen, so a facet lists
// only values some card actually has, each once.
func markSeen[T comparable](seen map[T]bool, out *[]T, v T) {
	if seen[v] {
		return
	}
	seen[v] = true
	*out = append(*out, v)
}

// nativeSetsOf is the set of set names a registered card is native to, from every
// provenance tag it carries.
func nativeSetsOf(rc card.RegisteredCard) map[string]bool {
	out := map[string]bool{}
	for _, p := range rc.Provenance {
		out[p.Set.Name] = true
	}
	return out
}

// reprintSetsByName maps a card name to the set names that reprint it, so a card
// can be filtered into a set it joins only as a reprint (ADR 0021).
func reprintSetsByName() map[string]map[string]bool {
	out := map[string]map[string]bool{}
	for _, r := range card.ReprintRefs() {
		if out[r.Name] == nil {
			out[r.Name] = map[string]bool{}
		}
		out[r.Name][r.Set.Name] = true
	}
	return out
}

// materializedDefs expands a registered card into the concrete definitions the
// gallery shows: a plain card is itself; a template card is every distinct variant
// its materializer produces, so the gallery lists Master of 1…5 rather than the
// "Master of X" placeholder. Variants are found by sampling the materializer until
// it stops yielding new names.
func materializedDefs(rc card.RegisteredCard) []engine.CardDefinition {
	if rc.Materializer == nil {
		return []engine.CardDefinition{rc.Def}
	}
	ctx := deckgen.SlotContext{
		House:  rc.Def.House,
		Rarity: rc.Def.Rarity,
	}
	r := rand.New(rand.NewSource(1))
	seen := map[string]engine.CardDefinition{}
	stale := 0
	for i := 0; i < 1000 && stale < 200; i++ {
		d := rc.Materializer.Materialize(ctx, r)
		if _, ok := seen[d.Name]; ok {
			stale++
			continue
		}
		seen[d.Name] = d
		stale = 0
	}
	out := make([]engine.CardDefinition, 0, len(seen))
	for name := range seen {
		out = append(out, seen[name])
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// matches reports whether a catalog entry passes every active facet (OR within a
// category, AND across categories), both text queries, and the power/armor ranges.
func (g *gallery) matches(c galleryCard, nameTerms, textTerms []queryTerm) bool {
	if len(g.selHouse) > 0 && !g.selHouse[c.def.House] {
		return false
	}
	if len(g.selSet) > 0 && !g.matchesSet(c) {
		return false
	}
	if len(g.selType) > 0 && !g.selType[c.def.Type] {
		return false
	}
	if len(g.selRarity) > 0 && !g.selRarity[c.def.Rarity] {
		return false
	}
	if len(g.selKeyword) > 0 && !anySelected(g.selKeyword, c.def.Keywords) {
		return false
	}
	if len(g.selTrait) > 0 && !anySelected(g.selTrait, c.def.Traits) {
		return false
	}
	if !inRange(c.def.Power, g.minPower, g.maxPower) {
		return false
	}
	if !inRange(c.def.Armor, g.minArmor, g.maxArmor) {
		return false
	}
	return matchesQuery(nameTerms, c.nameHay) && matchesQuery(textTerms, c.textHay)
}

// matchesSet reports whether a card belongs to any selected set, counting reprint
// memberships unless the exclude-reprints toggle is on.
func (g *gallery) matchesSet(c galleryCard) bool {
	for s := range g.selSet {
		if c.nativeSets[s] {
			return true
		}
		if !g.excludeReprints && c.reprintSets[s] {
			return true
		}
	}
	return false
}

// anySelected reports whether any of a card's values is selected in the facet set.
func anySelected[T comparable](sel map[T]bool, have []T) bool {
	for _, v := range have {
		if sel[v] {
			return true
		}
	}
	return false
}

// inRange reports whether v falls within the optional [min, max] bounds parsed
// from their input strings; an empty or unparsable bound does not constrain.
func inRange(v int, minStr, maxStr string) bool {
	if lo, ok := atoi(minStr); ok && v < lo {
		return false
	}
	if hi, ok := atoi(maxStr); ok && v > hi {
		return false
	}
	return true
}

// atoi parses an integer from a string, reporting false for empty or invalid
// input so it reads as "no bound".
func atoi(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (g *gallery) toggleHouse(h engine.House) app.EventHandler {
	return func(app.Context, app.Event) { toggleKey(g.selHouse, h) }
}

func (g *gallery) toggleSet(s string) app.EventHandler {
	return func(app.Context, app.Event) { toggleKey(g.selSet, s) }
}

func (g *gallery) toggleType(t engine.CardType) app.EventHandler {
	return func(app.Context, app.Event) { toggleKey(g.selType, t) }
}

func (g *gallery) toggleRarity(r engine.Rarity) app.EventHandler {
	return func(app.Context, app.Event) { toggleKey(g.selRarity, r) }
}

func (g *gallery) toggleKeyword(k engine.Keyword) app.EventHandler {
	return func(app.Context, app.Event) { toggleKey(g.selKeyword, k) }
}

func (g *gallery) toggleTrait(t engine.Trait) app.EventHandler {
	return func(app.Context, app.Event) { toggleKey(g.selTrait, t) }
}

// toggleKey flips a key in a selection set, deleting it when it was already on so
// an unselected facet leaves no empty entry behind.
func toggleKey[K comparable](m map[K]bool, k K) {
	if m[k] {
		delete(m, k)
	} else {
		m[k] = true
	}
}

// bindText returns a handler that stores an input's value into the given field.
func bindText(dst *string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) { *dst = ctx.JSSrc().Get("value").String() }
}

func (g *gallery) onExcludeReprints(ctx app.Context, _ app.Event) {
	g.excludeReprints = ctx.JSSrc().Get("checked").Bool()
}

// Render draws the filter sidebar and the ordered grid of matching card faces.
func (g *gallery) Render() app.UI {
	// A render reflows the grid (a filter may have changed the count), so ask for a
	// re-measure on the next frame to re-align the drawn-face window with the view.
	g.scheduleMeasure()
	nameTerms := parseQuery(g.name)
	textTerms := parseQuery(g.text)
	shown := make([]galleryCard, 0, len(g.cards))
	for _, c := range g.cards {
		if g.matches(c, nameTerms, textTerms) {
			shown = append(shown, c)
		}
	}
	g.orderCards(shown)

	return app.Div().Class("doc-page", "gallery-page").Body(
		// The icon-outline filter is injected here too: the doc pages do not carry
		// the game frame that defines it, so without this every card icon on the
		// gallery would lose its border.
		app.Raw(iconOutlineFilter),
		docHeader("/cards"),
		app.Section().Class("doc-section").Body(
			app.H2().Class("doc-h2").Text("Cards"),
			app.Div().Class("gallery-layout").Body(
				g.sidebar(),
				// A full-height splitter the reader drags to resize the sidebar; the
				// drag is wired in galleryScript, which sets --gallery-sidebar-w.
				app.Div().Class("gallery-divider"),
				app.Div().Class("gallery-main").Body(
					app.If(!g.ready, func() app.UI {
						return app.Div().Class("gallery-count").Text("Loading cards…")
					}).Else(func() app.UI {
						from, to := g.window(len(shown))
						return app.Div().Body(
							app.Div().Class("gallery-count").
								Text(fmt.Sprintf("%d of %d cards", len(shown), len(g.cards))),
							app.Div().Class("gallery-grid").Body(
								app.Range(shown).Slice(func(i int) app.UI {
									// On/near-screen cards draw the heavy printed face; the rest
									// draw a fixed-size text placeholder, so the scroll height and
									// find-in-page text stay intact without mounting every face.
									if i >= from && i < to {
										return app.Div().
											Class("gallery-card").
											Body(printedFace(shown[i].def))
									}
									return galleryPlaceholder(shown[i].def)
								}),
							),
						)
					}),
				),
			),
		),
	)
}

// galleryPlaceholder is the lightweight stand-in an off-screen card draws instead
// of its heavy printed face: the card's name and rules text in a box the same
// fixed size as a face. It keeps the grid's geometry uniform — so the scroll
// height stays exact and the drawn-face window's index math holds — and keeps the
// card's text in the DOM, so the browser's find-in-page (Ctrl+F) still matches a
// card whose full face has not been mounted.
func galleryPlaceholder(def *engine.CardDefinition) app.UI {
	return app.Div().Class("gallery-card", "gallery-card--placeholder").Body(
		app.Div().Class("gallery-placeholder-name").Text(def.Name),
		app.Div().Class("gallery-placeholder-text").Text(engine.RenderCardRules(def)),
	)
}

// orderCards sorts the shown cards by the chosen order key, always breaking ties
// on name so the grid is stable.
func (g *gallery) orderCards(cs []galleryCard) {
	byName := func(i, j int) bool { return cs[i].def.Name < cs[j].def.Name }
	switch g.order {
	case "power":
		sort.SliceStable(cs, func(i, j int) bool {
			if cs[i].def.Power != cs[j].def.Power {
				return cs[i].def.Power > cs[j].def.Power
			}
			return byName(i, j)
		})
	case "armor":
		sort.SliceStable(cs, func(i, j int) bool {
			if cs[i].def.Armor != cs[j].def.Armor {
				return cs[i].def.Armor > cs[j].def.Armor
			}
			return byName(i, j)
		})
	case "house":
		sort.SliceStable(cs, func(i, j int) bool {
			if cs[i].def.House != cs[j].def.House {
				return cs[i].def.House.String() < cs[j].def.House.String()
			}
			return byName(i, j)
		})
	case "type":
		sort.SliceStable(cs, func(i, j int) bool {
			if cs[i].def.Type != cs[j].def.Type {
				return cs[i].def.Type.String() < cs[j].def.Type.String()
			}
			return byName(i, j)
		})
	default:
		sort.SliceStable(cs, byName)
	}
}

// sidebar draws the named filter groups down the left so the grid keeps its
// vertical space.
func (g *gallery) sidebar() app.UI {
	return app.Div().Class("gallery-sidebar").Body(
		g.filterGroup("Order", app.Select().Class("gallery-order").OnChange(bindText(&g.order)).Body(
			orderOption("name", "Name", g.order),
			orderOption("house", "House", g.order),
			orderOption("type", "Type", g.order),
			orderOption("power", "Power", g.order),
			orderOption("armor", "Armor", g.order),
		)),
		g.filterGroup("Name", app.Input().Class("gallery-search").Type("text").
			Placeholder("card name").Value(g.name).OnInput(bindText(&g.name))),
		g.filterGroup("Card text", app.Input().Class("gallery-search").Type("text").
			Placeholder(`term term, a|b, -x, "phrase"`).Value(g.text).OnInput(bindText(&g.text))),
		g.filterGroup("House", g.houseChips()),
		g.filterGroup("Type", g.typeChips()),
		g.filterGroup("Rarity", g.rarityChips()),
		app.If(len(g.sets) > 0, func() app.UI {
			return g.filterGroup("Set", app.Div().Body(
				g.setChips(),
				app.Label().Class("gallery-check").Body(
					app.Input().Type("checkbox").Checked(g.excludeReprints).
						OnChange(g.onExcludeReprints),
					app.Text("Exclude reprints"),
				),
			))
		}),
		app.If(len(g.keywords) > 0, func() app.UI {
			return g.filterGroup("Keywords", g.keywordChips())
		}),
		g.filterGroup("Power", g.rangeInputs(&g.minPower, &g.maxPower)),
		g.filterGroup("Armor", g.rangeInputs(&g.minArmor, &g.maxArmor)),
		app.If(len(g.traits) > 0, func() app.UI {
			return g.collapsibleGroup("Traits", g.traitChips())
		}),
	)
}

// filterGroup wraps a named filter section with its label.
func (g *gallery) filterGroup(name string, body app.UI) app.UI {
	return app.Div().Class("gallery-group").Body(
		app.Div().Class("gallery-group-label").Text(name),
		body,
	)
}

// collapsibleGroup is a filter section the reader can fold away, used for the long
// Traits list so it does not dominate the sidebar.
func (g *gallery) collapsibleGroup(name string, body app.UI) app.UI {
	return app.Details().Class("gallery-group", "gallery-collapsible").Body(
		app.Summary().Class("gallery-group-label").Text(name),
		body,
	)
}

// orderOption is one entry of the ordering select, marked selected when it is the
// active key.
func orderOption(value, label, active string) app.UI {
	o := app.Option().Value(value).Text(label)
	if value == active {
		o = o.Selected(true)
	}
	return o
}

// rangeInputs draws the paired min/max number boxes for a numeric facet.
func (g *gallery) rangeInputs(lo, hi *string) app.UI {
	return app.Div().Class("gallery-range").Body(
		app.Input().Class("gallery-num").Type("number").Placeholder("min").
			Value(*lo).OnInput(bindText(lo)),
		app.Span().Class("gallery-range-dash").Text("–"),
		app.Input().Class("gallery-num").Type("number").Placeholder("max").
			Value(*hi).OnInput(bindText(hi)),
	)
}

func (g *gallery) houseChips() app.UI {
	return app.Div().Class("gallery-facets").Body(
		app.Range(g.houses).Slice(func(i int) app.UI {
			h := g.houses[i]
			return app.Button().
				Class(cx("gallery-chip", "gallery-chip--house", houseAccent(h),
					ifCls(g.selHouse[h], "gallery-chip--on"))).
				OnClick(g.toggleHouse(h)).
				Body(houseIcon(h, "icon-inline"), app.Text(h.String()))
		}),
	)
}

func (g *gallery) typeChips() app.UI {
	return app.Div().Class("gallery-facets").Body(
		app.Range(g.types).Slice(func(i int) app.UI {
			t := g.types[i]
			return app.Button().
				Class(cx("gallery-chip", ifCls(g.selType[t], "gallery-chip--on"))).
				OnClick(g.toggleType(t)).
				Body(icon(typeIconName(t), "icon-inline", "icon-outline"), app.Text(t.String()))
		}),
	)
}

func (g *gallery) rarityChips() app.UI {
	return app.Div().Class("gallery-facets").Body(
		app.Range(g.rarities).Slice(func(i int) app.UI {
			r := g.rarities[i]
			return app.Button().
				Class(cx("gallery-chip", ifCls(g.selRarity[r], "gallery-chip--on"))).
				OnClick(g.toggleRarity(r)).
				Body(app.Span().Class("gallery-rarity-mark").Body(rarityChipMarks(r)...),
					app.Text(string(r)))
		}),
	)
}

func (g *gallery) setChips() app.UI {
	return app.Div().Class("gallery-facets").Body(
		app.Range(g.sets).Slice(func(i int) app.UI {
			s := g.sets[i]
			return app.Button().
				Class(cx("gallery-chip", ifCls(g.selSet[s], "gallery-chip--on"))).
				OnClick(g.toggleSet(s)).
				Body(app.Span().Class(cx("gallery-set-emblem", setAccent(s))),
					app.Text(s))
		}),
	)
}

// rarityChipMarks is a rarity's mark for its filter chip: the single rarity shape,
// so the chips match the shapes on the card face and deck list.
func rarityChipMarks(r engine.Rarity) []app.UI {
	name := rarityMarkOf(r).iconName()
	if name == "" {
		return nil
	}
	return []app.UI{icon(name, "icon-inline", "icon-outline")}
}

func (g *gallery) keywordChips() app.UI {
	return app.Div().Class("gallery-facets").Body(
		app.Range(g.keywords).Slice(func(i int) app.UI {
			k := g.keywords[i]
			return app.Button().
				Class(cx("gallery-chip", ifCls(g.selKeyword[k], "gallery-chip--on"))).
				OnClick(g.toggleKeyword(k)).
				Text(k.String())
		}),
	)
}

func (g *gallery) traitChips() app.UI {
	return app.Div().Class("gallery-facets").Body(
		app.Range(g.traits).Slice(func(i int) app.UI {
			t := g.traits[i]
			return app.Button().
				Class(cx("gallery-chip", ifCls(g.selTrait[t], "gallery-chip--on"))).
				OnClick(g.toggleTrait(t)).
				Text(t.String())
		}),
	)
}
