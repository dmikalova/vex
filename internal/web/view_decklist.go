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
		Class(cx("deck-tip", ifCls(g.deckOpen == player, "deck-tip--open"))).
		OnMouseEnter(g.onDeckHover).
		Body(
			app.Span().Class("deck-tip-btn").
				OnClick(g.onDeckToggle(player)).
				Body(icon("deck-list", "icon-inline")),
			g.deckListPopover(player),
		)
}

// onDeckHover clamps the popover on screen as soon as a hover begins to open it,
// before the wide list can spill off the viewport edge.
func (g *game) onDeckHover(ctx app.Context, _ app.Event) {
	g.clampPopover(ctx.JSSrc())
}

// clampOpenDeckList clamps a tap-opened popover on screen after a re-render,
// since a tap opens it with no hover event to clamp against.
func (g *game) clampOpenDeckList() {
	if g.deckOpen < 0 {
		return
	}
	doc := app.Window().Get("document")
	if !doc.Truthy() {
		return
	}
	g.clampPopover(doc.Call("querySelector", ".deck-tip--open"))
}

// clampPopover keeps a centered deck-list popover on screen: it measures at the
// centered position and, only if an edge runs off the viewport, nudges the list
// back inward — favoring the left wall so the list is never cut off there. tip is
// the .deck-tip container.
func (g *game) clampPopover(tip app.Value) {
	if !tip.Truthy() {
		return
	}
	list := tip.Call("querySelector", ".deck-list")
	if !list.Truthy() {
		return
	}
	style := list.Get("style")
	style.Set("transform", "translateX(-50%)")
	rect := list.Call("getBoundingClientRect")
	const margin = 8.0
	vw := app.Window().Get("innerWidth").Float()
	left := rect.Get("left").Float()
	right := rect.Get("right").Float()
	shift := 0.0
	switch {
	case left < margin:
		shift = margin - left
	case right > vw-margin:
		shift = -(right - (vw - margin))
		if left+shift < margin {
			shift = margin - left
		}
	}
	if shift != 0 {
		style.Set("transform", fmt.Sprintf("translateX(calc(-50%% + %.0fpx))", shift))
	}
}

// onDeckToggle pins the deck list open on a tap and closes it on a second tap, so
// a touchscreen with no hover can still read it. Desktop hover works regardless.
func (g *game) onDeckToggle(player int) app.EventHandler {
	return func(_ app.Context, _ app.Event) {
		if g.deckOpen == player {
			g.deckOpen = -1
			return
		}
		g.deckOpen = player
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
	header := g.g.PlayerName(player)
	if roster.Set != "" {
		header += " • " + roster.Set
	}
	return app.Div().Class("deck-list").Body(
		app.Div().Class("deck-list-head").Text(header),
		app.Div().Class("deck-list-cols").Body(cols...),
	)
}

// deckListRow renders one roster card: its type and rarity marks (kept tight
// together) then its name, with any Maverick/Legacy mark. Rarity is the card's
// own, so a legacy card shows the same rarity as its home-set printing.
func deckListRow(c match.RosterCard) app.UI {
	marks := []app.UI{}
	if name := typeIconName(c.Def.Type); name != "" {
		marks = append(marks, icon(name, "icon-mark"))
	}
	if mark := deckRarityIcon(c.Def.Rarity); mark != nil {
		marks = append(marks, mark)
	}
	cells := []app.UI{app.Span().Class("deck-list-marks").Body(marks...)}
	cells = append(cells, app.Span().Class("deck-list-name").Text(c.Def.Name))
	if c.Maverick {
		cells = append(cells, icon("maverick", "icon-mark", "icon-outline"))
	}
	if c.Legacy {
		cells = append(cells, icon("legacy", "icon-mark", "icon-outline"))
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
