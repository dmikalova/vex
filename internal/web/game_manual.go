package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is manual mode: the handlers that let a playtester move,
// ready, and exhaust cards, adjust Æmber, chains, keys and the active house by
// hand, and add any card in the game to a hand.

// toggleManual turns the engine's manual mode on or off, lifting house
// restrictions and revealing the manual controls.
func (g *game) toggleManual(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.choosingOption {
		return
	}
	g.g.SetManual(!g.g.Manual())
	g.save(ctx)
}

// manualMove moves the selected card to a resting zone, ignoring the normal rules.
func (g *game) manualMove(dest engine.ManualZone) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if !g.hasSel || !g.g.Manual() {
			return
		}
		g.beginAction()
		g.g.ManualMove(g.sel, dest)
		g.clearSelection()
		g.save(ctx)
	}
}

// manualReady clears the selected card's exhausted flag.
func (g *game) manualReady(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.beginAction()
	g.g.ManualSetExhausted(g.sel, false)
	g.save(ctx)
}

// manualExhaust sets the selected card's exhausted flag.
func (g *game) manualExhaust(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.beginAction()
	g.g.ManualSetExhausted(g.sel, true)
	g.save(ctx)
}

// manualAmberDelta adjusts a player's Æmber in manual mode. It stops the click
// from bubbling to the score pill (whose click opens the zone viewer).
func (g *game) manualAmberDelta(player, delta int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		e.Call("stopPropagation")
		if !g.g.Manual() {
			return
		}
		g.beginAction()
		g.g.ManualAddAmber(player, delta)
		g.save(ctx)
	}
}

// manualForgeKey opens the key-forge colour picker for player in manual mode.
func (g *game) manualForgeKey(player int) app.EventHandler {
	return func(_ app.Context, e app.Event) {
		e.Call("stopPropagation")
		if !g.g.Manual() {
			return
		}
		g.forgingKey = player
	}
}

// manualUnforgeKey removes player's most recently forged key in manual mode.
func (g *game) manualUnforgeKey(player int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		e.Call("stopPropagation")
		if !g.g.Manual() {
			return
		}
		g.beginAction()
		g.g.ManualUnforgeKey(player)
		g.save(ctx)
	}
}

// manualChainsDelta adjusts a player's chains in manual mode. It stops the click
// from bubbling to the score pill (whose click opens the zone viewer).
func (g *game) manualChainsDelta(player, delta int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		e.Call("stopPropagation")
		if !g.g.Manual() {
			return
		}
		g.beginAction()
		g.g.ManualAddChains(player, delta)
		g.save(ctx)
	}
}

// manualSetHouse switches the active player's active house in manual mode; from
// the house-choice step it also advances play, like picking a house normally.
func (g *game) manualSetHouse(h engine.House) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if !g.g.Manual() {
			return
		}
		g.beginAction()
		g.g.ManualSetActiveHouse(h)
		if g.phase == phaseHouse {
			g.phase = phaseMain
		}
		g.save(ctx)
	}
}

// pickForgeColor forges the chosen colour for the player whose picker is open.
func (g *game) pickForgeColor(c engine.KeyColor) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if g.forgingKey < 0 {
			return
		}
		g.beginAction()
		g.g.ManualForgeKeyColor(g.forgingKey, c)
		g.forgingKey = -1
		g.save(ctx)
	}
}

// cancelForgeKey closes the key-forge picker without forging.
func (g *game) cancelForgeKey(_ app.Context, _ app.Event) { g.forgingKey = -1 }

// chooseKeyColorKey answers a key-forge prompt — manual mode's colour picker or
// an ordinary forge choice offered as key-colour option buttons — with color, so
// r/b/y can forge a key directly instead of hunting the matching button. It
// reports whether a matching prompt was up to answer, so the caller can fall
// back to its own binding for the key when none was (r doubles as "affirm").
func (g *game) chooseKeyColorKey(ctx app.Context, color engine.KeyColor) bool {
	if g.forgingKey >= 0 {
		for _, c := range g.remainingKeyColors(g.forgingKey) {
			if c == color {
				g.pickForgeColor(c)(ctx, app.Event{})
				return true
			}
		}
		return false
	}
	if g.choosingOption && g.keyColorOptions() {
		for i, label := range g.optionLabels {
			if keyColorByName(label) == color {
				g.chooseOptionIdx(i)(ctx, app.Event{})
				return true
			}
		}
	}
	return false
}

// selectZoneCard selects a card shown in the zone viewer (in manual mode) and
// closes the viewer, so the manual controls can act on it (e.g. move it to hand).
func (g *game) selectZoneCard(_ app.Context, id engine.LocalID) {
	g.sel, g.selKind, g.selHand, g.hasSel = id, selOther, -1, true
	g.zonesPlayer = -1
	g.status = ""
}

// openPicker opens the fuzzy card picker to add an arbitrary card to hand.
func (g *game) openPicker(_ app.Context, _ app.Event) {
	g.pickerOpen = true
	g.pickerQuery = ""
	g.pickerFocused = false
}

func (g *game) closePicker(_ app.Context, _ app.Event) { g.pickerOpen = false }

// pickerInput records the search box's text as the player types.
func (g *game) pickerInput(ctx app.Context, _ app.Event) {
	g.pickerQuery = ctx.JSSrc().Get("value").String()
}

// addPickedCard adds the clicked picker row's card to the active player's hand.
// The card is named by the row's data attribute rather than captured per row:
// go-app compares handlers by pointer, so a closure per row goes stale the moment
// the search filters the list, and the click adds whatever card used to sit there.
func (g *game) addPickedCard(ctx app.Context, _ app.Event) {
	def, ok := g.defByName[ctx.JSSrc().Get("dataset").Get("card").String()]
	if !ok {
		return
	}
	g.beginAction()
	player := g.active()
	if _, added := g.g.ManualAddCard(*def, player); added {
		// The catalog is not part of the saved state, so record the add for a
		// reload to replay; undo rolls the state back but not the registration.
		g.manualAdds = append(g.manualAdds, manualAdd{Name: def.Name, Player: player})
	}
	g.pickerOpen = false
	g.save(ctx)
}

// isInPlay reports whether a card is on either player's battleline or artifact row.
func (g *game) isInPlay(id engine.LocalID) bool {
	for p := 0; p < 2; p++ {
		if containsID(g.g.Battleline(p), id) || containsID(g.g.Artifacts(p), id) {
			return true
		}
	}
	return false
}
