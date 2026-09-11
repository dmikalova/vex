package web

import (
	"fmt"
	"sort"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws the active player's hand: the row itself and the cards in it,
// which are the only cards that can be dragged onto the board.

func (g *game) renderHand() app.UI {
	p := g.active()
	ids := g.sortedHand(p)
	return app.Div().Class("board-row").Body(
		app.Div().Class("row-label").Body(
			app.Span().Class("row-label-zone").Text("Hand "),
			app.Text(fmt.Sprintf("%d", len(ids))),
			icon("zone-hand", "row-label-icon"),
		),
		app.Div().Class("card-strip").Body(
			app.Range(ids).Slice(func(i int) app.UI { return g.renderHandCard(ids[i]) }),
		),
	)
}

// sortedHand returns the player's hand ids ordered by house, then card type, then
// name, so the hand reads consistently. It sorts a copy — the engine's own hand
// order (which play/discard index into) is untouched, so selection still maps to
// the right card.
func (g *game) sortedHand(p int) []engine.LocalID {
	return g.sortByHouseTypeName(g.g.Hand(p))
}

// sortedArtifacts returns the player's artifact-row ids ordered by house, then
// name, so the artifact line reads consistently instead of by play order. It
// sorts a copy — the engine's own order is untouched, so selection still maps to
// the right card.
func (g *game) sortedArtifacts(p int) []engine.LocalID {
	return g.sortByHouseTypeName(g.g.Artifacts(p))
}

// sortByHouseTypeName returns a copy of ids ordered by house, then card type in
// the deck list's order (typeRank — creatures, artifacts, upgrades, Tactics last;
// ADR 0025), then name — the stable reading order shared by the hand, the deck
// view, and the deck list. The deck in particular must not reveal its shuffled
// order, so it is always sorted.
func (g *game) sortByHouseTypeName(ids []engine.LocalID) []engine.LocalID {
	ids = append([]engine.LocalID(nil), ids...)
	sort.SliceStable(ids, func(i, j int) bool {
		a, b := g.g.Def(ids[i]), g.g.Def(ids[j])
		if a.House != b.House {
			return a.House < b.House
		}
		if ra, rb := typeRank(a.Type), typeRank(b.Type); ra != rb {
			return ra < rb
		}
		return a.Name < b.Name
	})
	return ids
}

func (g *game) renderHandCard(id engine.LocalID) app.UI {
	def := g.g.Def(id)
	activate, targetable, dimmed := g.cardVisual(id, selHand)
	draggable := !g.busy && !g.choosing && !g.choosingOption &&
		g.phase == phaseMain && g.playableFromHand(id)
	// A hand card is a printed face plus the hand's interaction wiring.
	face := printedFace(def)
	face.ID = id
	face.DOMID = handCardID(id)
	face.Maverick = g.isMaverick(id)
	face.Legacy = g.isLegacy(id)
	face.Selected = g.isSelected(id)
	face.Targetable = targetable
	face.Dimmed = dimmed
	face.Jiggle = g.jiggling(id, selHand)
	face.OnActivate = activate
	face.Draggable = draggable
	face.OnDragStart = g.startHandDrag
	face.OnDragEnd = g.endHandDrag
	face.OnHover = g.hoverCard
	face.OnHoverOut = g.hoverClear
	face.OnContextMenu = g.liftCard
	return face
}
