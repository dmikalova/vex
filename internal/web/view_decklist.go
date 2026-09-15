package web

import (
	"fmt"
	"sort"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// This file draws the deck list: the static roster of a player's generated deck
// — its three houses each with their twelve cards — shown from a deck icon on the
// Player bar as a wide popover. It reads the retained roster (ADR 0025), never the
// live piles, so it leaks nothing about draw order.

// deckListVisible reports whether viewer may see owner's deck list. Every game is
// the open format today, so it is always true; a future sealed or hidden-list
// format flips this one predicate without touching the render path.
func (g *game) deckListVisible(_, _ int) bool { return true }

// deckTip is the Player-bar deck icon that opens the deck list, or nil when there
// is no roster (the style gallery) or the viewer may not read it. Desktop hover
// opens the popover; a tap pins it open (deckOpen) and a tap outside clears it.
func (g *game) deckTip(player int) app.UI {
	if g.rosters[player].Empty() || !g.deckListVisible(g.active(), player) {
		return nil
	}
	return app.Span().
		Class(cx("deck-tip", ifCls(g.deckOpen[player], "deck-tip--open"))).
		OnMouseEnter(g.onDeckHover).
		Body(
			app.Span().Class("deck-tip-btn").
				OnClick(g.onDeckToggle(player)).
				Body(icon("deck-list", "icon-inline")),
			g.deckListPopover(player),
		)
}

// onDeckHover places the popover on screen as soon as a hover begins to open it,
// before the wide list can spill off a viewport edge.
func (g *game) onDeckHover(ctx app.Context, _ app.Event) {
	g.placePopover(ctx.JSSrc())
}

// clampOpenDeckList places every tap-opened popover on screen after a re-render,
// since a tap opens it with no hover event to place against.
func (g *game) clampOpenDeckList() {
	if g.deckOpen == ([2]bool{}) {
		return
	}
	doc := app.Window().Get("document")
	if !doc.Truthy() {
		return
	}
	open := doc.Call("querySelectorAll", ".deck-tip--open")
	for i := range open.Get("length").Int() {
		g.placePopover(open.Call("item", i))
	}
}

// placePopover positions the deck-list popover for one .deck-tip. tip is the
// .deck-tip container.
func (g *game) placePopover(tip app.Value) {
	placeFloating(tip, ".deck-list")
}

// placeFloating positions a popover (a child of anchor matched by sel) with
// position: fixed, so it escapes the player bar's overflow clip the way a floating
// tip does. It opens away from the bar the anchor sits in — up from a bottom-half
// anchor, down from a top-half one — centred on the anchor and clamped inside the
// window; the popover's own max-height and scroll handle a list taller than the
// gap to the edge. It is shared by the deck list and the zone-count rosters.
func placeFloating(anchor app.Value, sel string) {
	if !anchor.Truthy() {
		return
	}
	pop := anchor.Call("querySelector", sel)
	if !pop.Truthy() {
		return
	}
	const gap, margin = 6.0, 8.0
	a := anchor.Call("getBoundingClientRect")
	style := pop.Get("style")
	// Measure the popover at the origin — its natural size, unaffected by a prior
	// placement — then set left/top from the anchor.
	style.Set("left", "0")
	style.Set("top", "0")
	p := pop.Call("getBoundingClientRect")
	vw := app.Window().Get("innerWidth").Float()
	vh := app.Window().Get("innerHeight").Float()
	w := p.Get("width").Float()
	mid := a.Get("left").Float() + a.Get("width").Float()/2
	left := mid - w/2
	if left < margin {
		left = margin
	} else if left+w > vw-margin {
		if left = vw - margin - w; left < margin {
			left = margin
		}
	}
	var top float64
	if a.Get("top").Float() > vh/2 {
		if top = a.Get("top").Float() - gap - p.Get("height").Float(); top < margin {
			top = margin
		}
	} else {
		top = a.Get("bottom").Float() + gap
	}
	style.Set("left", px(left))
	style.Set("top", px(top))
	fitCondensedNames(pop)
}

// fitCondensedNames condenses each deck-list card name horizontally (scaleX) so a
// long name shrinks to fit its column instead of truncating with an ellipsis,
// mirroring the card banner's title fit (cmd/web cardFitScript). It runs here, on
// every popover placement, because the popover is shown by CSS hover with no
// re-render to trigger that client-side observer. A 0.7 floor keeps a very long
// name legible and leaves the rest to the ellipsis.
func fitCondensedNames(pop app.Value) {
	els := pop.Call("querySelectorAll", ".deck-list-name-text")
	for i := 0; i < els.Get("length").Int(); i++ {
		el := els.Call("item", i)
		style := el.Get("style")
		style.Set("transform", "")
		style.Set("width", "")
		avail := el.Get("clientWidth").Float()
		natural := el.Get("scrollWidth").Float()
		if avail <= 0 || natural <= avail {
			continue
		}
		scale := max(0.7, avail/natural)
		style.Set("transform", fmt.Sprintf("scaleX(%.4f)", scale))
		style.Set("width", fmt.Sprintf("%.2f%%", 100/scale))
	}
}

// onDeckToggle pins one player's deck list open on a tap and closes it on a second
// tap, so a touchscreen with no hover can still read it. Each side toggles on its
// own, leaving the other player's list as it was. Desktop hover works regardless.
// It routes through dispatch so OnUpdate runs afterward and places the pinned
// popover (clampOpenDeckList) — a plain handler's re-render never fires OnUpdate.
func (g *game) onDeckToggle(player int) app.EventHandler {
	return func(_ app.Context, _ app.Event) {
		g.dispatch(func(app.Context) { g.deckOpen[player] = !g.deckOpen[player] })
	}
}

// deckListPopover renders a player's roster as one column per house — the house
// header then its twelve cards, each a type icon, rarity shape, and name with any
// Maverick/Legacy mark — sorted within a house by type then name (ADR 0025).
func (g *game) deckListPopover(player int) app.UI {
	roster := g.rosters[player]
	cols := make([]app.UI, 0, len(roster.Houses))
	for _, hr := range roster.Houses {
		cards := hr.Cards[:]
		sort.SliceStable(cards, func(i, j int) bool {
			if ri, rj := typeRank(cards[i].Def.Type), typeRank(cards[j].Def.Type); ri != rj {
				return ri < rj
			}
			return cards[i].Def.Name < cards[j].Def.Name
		})
		rows := make([]app.UI, 0, len(cards)+1)
		rows = append(rows,
			app.Div().Class(cx("deck-list-colhead", houseClasses(hr.House))).Body(
				houseIcon(hr.House, "icon-inline"),
				app.Span().Class("deck-list-house").Text(hr.House.String()),
			),
		)
		for _, c := range cards {
			rows = append(rows, deckListRow(c))
		}
		cols = append(cols, app.Div().Class("deck-list-col").Body(rows...))
	}
	suffix := ""
	if roster.Set != "" {
		suffix = " • " + roster.Set
	}
	return app.Div().Class("deck-list").Body(
		app.Div().Class("deck-list-head").Body(
			app.Span().Class(playerNameCls(player)).Text(g.g.PlayerName(player)),
			app.Text(suffix),
		),
		app.Div().Class("deck-list-cols").Body(cols...),
	)
}

// deckListRow renders one roster card: its type and rarity marks (kept tight
// together) then its name, then a tight, evenly-spaced trailing group of its
// Maverick/Legacy marks followed by its bonus icons. Rarity is the card's own, so
// a legacy card shows the same rarity as its home-set printing.
func deckListRow(c match.RosterCard) app.UI {
	marks := []app.UI{}
	if name := typeIconName(c.Def.Type); name != "" {
		marks = append(marks, icon(name, "icon-mark"))
	}
	if mark := deckRarityIcon(c.Def.Rarity); mark != nil {
		marks = append(marks, mark)
	}
	cells := []app.UI{app.Span().Class("deck-list-marks").Body(marks...)}
	cells = append(cells, app.Span().Class("deck-list-name").Body(
		app.Span().Class("deck-list-name-text").Text(c.Def.Name),
	))
	// Provenance marks first, then bonus icons, all in one tight group so they
	// share the same even spacing.
	tail := []app.UI{}
	if c.Maverick {
		tail = append(tail, icon("maverick", "icon-mark", "icon-outline"))
	}
	if c.Legacy {
		tail = append(tail, icon("legacy", "icon-mark", "icon-outline"))
	}
	for _, b := range c.Def.Bonuses {
		tail = append(tail, icon(bonusIconStem(b), "icon-mark", "icon-outline"))
	}
	if len(tail) > 0 {
		cells = append(cells, app.Span().Class("deck-list-tail").Body(tail...))
	}
	return app.Div().Class("deck-list-row").Body(cells...)
}

// typeRank orders card types for the deck list: creatures first as the bulk of a
// house, then artifacts and upgrades, then one-shot Tactics last.
func typeRank(t engine.CardType) int {
	switch t {
	case engine.Creature:
		return 0
	case engine.Artifact:
		return 1
	case engine.Upgrade:
		return 2
	case engine.Tactic:
		return 3
	}
	return 4
}
