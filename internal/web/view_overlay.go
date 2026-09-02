package web

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws what covers the board: the end-of-game panel and the overlay
// that lists the out-of-play zones.

// overPanel is the end-of-game result, shown in the controls area where every
// other action lives.
func (g *game) overPanel() app.UI {
	winner := g.g.Winner()
	return app.Div().Class("btn-col over-panel").Body(
		app.Div().Class("section-title").Text(g.g.PlayerName(winner)+" wins!"),
		btn("New game", g.restart, "btn-primary"),
	)
}

// promptZoneID marks the zone row a prompt is asking about, so it can be scrolled
// into view when the viewer opens.
const promptZoneID = "promptzone"

// zonesOverlay shows a single player's zones — their deck, discard, archives,
// and purge piles — as read-only card strips, so cards outside the board can be
// inspected. It is a modal dismissed by the ✕ in the corner or by clicking the
// backdrop outside the panel.
func (g *game) zonesOverlay() app.UI {
	p := g.zonesPlayer
	return app.Div().Class("over-backdrop").OnClick(g.closeZones).Body(
		app.Div().Class("zones-panel").OnClick(g.stopClick).Body(
			app.Div().Class("zones-header").Body(
				app.Button().Class("zones-close").Text("✕").OnClick(g.closeZones),
				app.Div().Class("over-title").Text(g.g.PlayerName(p)+"'s Zones"),
			),
			app.Div().Class("zones-body").Body(
				g.zoneRow("Deck", g.sortByHouseTypeName(g.g.Deck(p))),
				g.zoneRow("Discard", g.g.Discard(p)),
				g.zoneRow("Archives", g.g.Archives(p)),
				g.zoneRow("Purge", g.g.Purge(p)),
			),
		),
	)
}

func (g *game) zoneRow(label string, ids []engine.LocalID) app.UI {
	row := app.Div().Class("zone-row")
	if label == g.promptZone {
		row = row.ID(promptZoneID)
	}
	return row.Body(
		app.Div().Class("row-label").Text(fmt.Sprintf("%s (%d)", label, len(ids))),
		app.If(len(ids) == 0, func() app.UI {
			return app.Div().Class("row-empty").Text("—")
		}).Else(func() app.UI {
			return app.Div().Class("card-strip").Body(
				app.Range(ids).Slice(func(i int) app.UI { return g.renderZoneCard(ids[i]) }),
			)
		}),
	)
}

// renderZoneCard renders a card face for a card in an out-of-play zone. It is
// read-only except when something wants it clicked: a chooser whose candidates
// live in a pile (World Tree recovering a creature from the discard), or manual
// mode moving a card between zones.
func (g *game) renderZoneCard(id engine.LocalID) app.UI {
	var activate func(app.Context, engine.LocalID)
	targetable, dimmed := false, false
	switch {
	case g.choosing:
		// During a prompt the pile is the board: only the candidates are clickable,
		// and everything else dims the same way an unchoosable creature does.
		targetable = containsID(g.chooserCandidates, id)
		dimmed = !targetable
		if targetable {
			activate = g.chooseCandidate
		}
	case g.g.Manual():
		activate = g.selectZoneCard
	}
	c := g.printedCard(id)
	c.Targetable = targetable
	c.Dimmed = dimmed
	c.Selected = g.isSelected(id)
	c.OnActivate = activate
	return c
}

// printedCard is a read-only face built from a card's printed definition, for a
// card that is not on the board — one in a pile, or one in flight out of play.
func (g *game) printedCard(id engine.LocalID) *cardView {
	def := g.g.Def(id)
	return &cardView{
		ID:       id,
		Title:    def.Name,
		HouseCls: houseClasses(def.House),
		Emblem:   houseIconName(def.House),
		TypeIcon: typeIconName(def.Type),
		Stat:     handStat(def),
		Rules:    engine.RenderCardRules(def),
		Kind:     kindLabel(def),
		Trait:    traitLabel(def),
		Rarity:   rarityMarkOf(def.Rarity),
		Maverick: g.isMaverick(id),
	}
}

// shortcuts is the keyboard sheet the ? key opens, in the order it is read:
// moving around the board first, then acting, then the controls that frame a
// game.
var shortcuts = []struct{ keys, what string }{
	{"← → ↑ ↓", "Move the selection between cards and rows"},
	{"j k l ;", "The same four, on the home row (left, down, up, right)"},
	{"1 – 9", "Select the nth card of the selected card's row"},
	{"Tab / Shift+Tab", "Step through what you can act on, or a prompt's choices"},
	{"Enter / Space", "Answer the prompt with what Tab stopped on"},
	{"Esc", "Back out one layer: overlay, prompt, targeting, selection"},
	{"p", "Play the selected card from hand"},
	{"d", "Discard the selected card from hand"},
	{"a", "Use the selected card's Action ability"},
	{"f", "Fight with the selected creature"},
	{"u", "Unstun the selected creature"},
	{"r", "Yes: answer a prompt, take the right flank, or use the card"},
	{"n", "No: refuse a prompt you may decline"},
	{"l", "Take the left flank while placing a creature"},
	{"e", "End the turn (press twice when you could still act)"},
	{"z", "Cycle the out-of-play zone viewer"},
	{"h", "Hide or show the sidebar"},
	{"m", "Toggle manual mode"},
	{"⌘/Ctrl+Z", "Undo"},
	{"⇧⌘/Ctrl+Z", "Redo"},
	{"?", "Open or close this sheet"},
}

// keysOverlay lists every keyboard shortcut, so the ones that are not written on
// a button can still be found. It is dismissed like any other overlay.
func (g *game) keysOverlay() app.UI {
	return app.Div().Class("over-backdrop").OnClick(g.closeKeys).Body(
		app.Div().Class("keys-panel").OnClick(g.stopClick).Body(
			app.Div().Class("zones-header").Body(
				app.Button().Class("zones-close").Text("✕").OnClick(g.closeKeys),
				app.Div().Class("over-title").Text("Keyboard shortcuts"),
			),
			app.Div().Class("keys-grid").Body(
				app.Range(shortcuts).Slice(func(i int) app.UI {
					return app.Div().Class("keys-line").Body(
						app.Span().Class("keys-key").Text(shortcuts[i].keys),
						app.Span().Class("keys-what").Text(shortcuts[i].what),
					)
				}),
			),
		),
	)
}
