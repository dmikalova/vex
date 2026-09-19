package web

import (
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is manual mode: the handlers that let a playtester move,
// ready, and exhaust cards, adjust Æmber, chains, keys and the active house by
// hand, and add any card in the game to a hand.

// toggleManual turns the engine's manual mode on or off, lifting house
// restrictions and revealing the manual controls.
func (g *game) toggleManual(ctx app.Context, _ app.Event) {
	// Manual mode may be toggled mid-prompt (g.busy with a chooser waiting): turning
	// it on reveals the Cancel button that escapes a stuck prompt. Only a non-prompt
	// busy state (an effect still animating) blocks the toggle.
	if g.busy && !g.choosing && !g.choosingOption {
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
		g.record(input{Kind: inManualMove, ID: g.sel, Index: int(dest)})
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
	g.record(input{Kind: inManualReady, ID: g.sel})
	g.g.ManualSetExhausted(g.sel, false)
	g.save(ctx)
}

// manualGraft begins host targeting to thread the selected card face up under an
// in-play host (a graft). manualPlaceUnder is its face-down counterpart. Both
// only arm the targeting; the actual attach happens when a host is clicked
// (attachToHost).
func (g *game) manualGraft(_ app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.hostTargeting, g.hostFaceDown = true, false
}

// manualPlaceUnder begins host targeting to place the selected card face down
// under an in-play host.
func (g *game) manualPlaceUnder(_ app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.hostTargeting, g.hostFaceDown = true, true
}

// attachToHost threads the selected card under the clicked host — face up for a
// graft, face down for a place-under — then clears the targeting and selection.
func (g *game) attachToHost(ctx app.Context, host engine.LocalID) {
	if !g.hostTargeting || !g.hasSel || !g.g.Manual() || host == g.sel {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualAttach, Card: host, ID: g.sel, Left: g.hostFaceDown})
	g.g.ManualAttachUnder(host, g.sel, g.hostFaceDown)
	g.hostTargeting = false
	g.clearSelection()
	g.save(ctx)
}

// cancelHostTargeting backs out of a Graft / Place under host pick without
// attaching, leaving the card selected.
func (g *game) cancelHostTargeting(_ app.Context, _ app.Event) {
	g.hostTargeting = false
}

// manualPlay drops the selected hand card straight into play in manual mode,
// without its play effects or bonus Æmber. A creature enters the placement picker
// (reusing the Deploy line) so it can land anywhere in the battleline; with no
// other creatures to place it beside, and for a non-creature, it goes in at once.
func (g *game) manualPlay(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() || g.selKind != selHand {
		return
	}
	if g.g.IsCreature(g.sel) && len(g.g.Battleline(g.g.Owner(g.sel))) > 0 {
		g.manualPlacing = true
		g.choosingPosition = true
		g.positionLine = g.g.Battleline(g.g.Owner(g.sel))
		g.positionRight = false
		g.positionSideChosen = false
		return
	}
	g.manualPlaceInPlay(ctx, 0)
}

// manualPlaceInPlay commits a manual put-into-play at battleline position pos and
// clears the placement picker. It is the manual counterpart to answerPosition: a
// clicked position lands the creature here instead of replying to a prompt.
func (g *game) manualPlaceInPlay(ctx app.Context, pos int) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualPlace, ID: g.sel, Index: pos})
	g.g.ManualPlaceInPlay(g.sel, pos)
	g.manualPlacing = false
	g.choosingPosition = false
	g.positionLine = nil
	g.positionSideChosen = false
	g.clearSelection()
	g.save(ctx)
}

// cancelManualPlace backs out of a manual put-into-play placement without placing,
// leaving the card selected in hand.
func (g *game) cancelManualPlace(_ app.Context, _ app.Event) {
	g.manualPlacing = false
	g.choosingPosition = false
	g.positionLine = nil
	g.positionSideChosen = false
}

// manualToHand sends the selected upgrade or under-card to its owner's hand,
// detaching it from its host first. It is offered only when the selection is
// actually attached (isAttached).
func (g *game) manualToHand(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualDetach, ID: g.sel})
	g.g.ManualDetachToHand(g.sel)
	g.clearSelection()
	g.save(ctx)
}

// isAttached reports whether a card in play is an upgrade of, or placed under, a
// host — the two states ManualDetachToHand can send back to hand. Upgrades are
// found directly (HostOf); under-cards have no back-link reader, so the hosts'
// Under chains are scanned.
func (g *game) isAttached(id engine.LocalID) bool {
	if _, ok := g.g.HostOf(id); ok {
		return true
	}
	for p := range 2 {
		for _, host := range g.g.Battleline(p) {
			if containsID(g.g.Under(host), id) {
				return true
			}
		}
		for _, host := range g.g.Artifacts(p) {
			if containsID(g.g.Under(host), id) {
				return true
			}
		}
	}
	return false
}

// manualExhaust sets the selected card's exhausted flag.
func (g *game) manualExhaust(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualExhaust, ID: g.sel})
	g.g.ManualSetExhausted(g.sel, true)
	g.save(ctx)
}

// A manual per-player control reads its target player (and a stepper its signed
// step) from the button's own data attributes rather than a captured closure. go-app
// compares event handlers by function pointer, and a closure minted at one call site
// — manualAmberDelta(player, delta) — shares that pointer for both players and both
// signs, so the diff never re-binds it: one bar's stepper would drive whichever
// player was rendered first. A stable method that reads the DOM is bound once and
// always acts on the button that was clicked (the same pattern as onScorePillClick).

// stopManualClick stops a manual control's click from bubbling to the score pill
// (whose click opens the zone viewer) and reports whether manual mode is on.
func (g *game) stopManualClick(e app.Event) bool {
	e.Call("stopPropagation")
	return g.g.Manual()
}

// datasetPlayer reads a manual control's target player from its data-player
// attribute, or -1 when unreadable (off-browser, or a malformed attribute).
func (g *game) datasetPlayer(ctx app.Context) int {
	p, err := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("player").String())
	if err != nil {
		return -1
	}
	return p
}

// datasetDelta reads a stepper's signed step from its data-delta attribute, or 0
// (a no-op) when unreadable.
func (g *game) datasetDelta(ctx app.Context) int {
	d, err := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("delta").String())
	if err != nil {
		return 0
	}
	return d
}

// onManualAmberStep adjusts the clicked bar's player's Æmber in manual mode.
func (g *game) onManualAmberStep(ctx app.Context, e app.Event) {
	if !g.stopManualClick(e) {
		return
	}
	g.adjustManualAmber(ctx, g.datasetPlayer(ctx), g.datasetDelta(ctx))
}

// adjustManualAmber applies a manual Æmber step of delta to player.
func (g *game) adjustManualAmber(ctx app.Context, player, delta int) {
	if player < 0 || delta == 0 {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualAmber, Player: player, Delta: delta})
	g.g.ManualAddAmber(player, delta)
	g.save(ctx)
}

// onManualForgeKey opens the key-forge colour picker for the clicked bar's player.
func (g *game) onManualForgeKey(ctx app.Context, e app.Event) {
	if !g.stopManualClick(e) {
		return
	}
	g.openForgeKey(g.datasetPlayer(ctx))
}

// openForgeKey aims the key-forge colour picker at player.
func (g *game) openForgeKey(player int) {
	if player >= 0 {
		g.forgingKey = player
	}
}

// onManualUnforgeKey removes the clicked bar's player's most recently forged key.
func (g *game) onManualUnforgeKey(ctx app.Context, e app.Event) {
	if !g.stopManualClick(e) {
		return
	}
	g.removeManualKey(ctx, g.datasetPlayer(ctx))
}

// removeManualKey unforges player's most recently forged key in manual mode.
func (g *game) removeManualKey(ctx app.Context, player int) {
	if player < 0 {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualUnforge, Player: player})
	g.g.ManualUnforgeKey(player)
	g.save(ctx)
}

// onManualChainsStep adjusts the clicked bar's player's chains in manual mode.
func (g *game) onManualChainsStep(ctx app.Context, e app.Event) {
	if !g.stopManualClick(e) {
		return
	}
	g.adjustManualChains(ctx, g.datasetPlayer(ctx), g.datasetDelta(ctx))
}

// adjustManualChains applies a manual chains step of delta to player.
func (g *game) adjustManualChains(ctx app.Context, player, delta int) {
	if player < 0 || delta == 0 {
		return
	}
	g.beginAction()
	g.record(input{Kind: inManualChains, Player: player, Delta: delta})
	g.g.ManualAddChains(player, delta)
	g.save(ctx)
}

// manualSetHouse switches the active player's active house in manual mode; from
// the house-choice step it also advances play, like picking a house normally.
func (g *game) manualSetHouse(h engine.House) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if !g.g.Manual() {
			return
		}
		g.beginAction()
		g.record(input{Kind: inManualHouse, House: h})
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
		g.record(input{Kind: inManualForgeColor, Player: g.forgingKey, Index: int(c)})
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
	g.pickerNaming = false
	g.pickerQuery = ""
	g.pickerFocused = false
	g.pickerCursor = 0
}

// closePicker dismisses the picker. A picker answering a name-a-card prompt is
// not dismissible: the prompt is waiting on it, so closing it would hang the
// action with nothing on screen to answer.
func (g *game) closePicker(_ app.Context, _ app.Event) {
	if g.pickerNaming {
		return
	}
	g.pickerOpen = false
}

// pickerInput records the search box's text as the player types. Refiltering
// resets the cursor to the first row so Enter adds the top match.
func (g *game) pickerInput(ctx app.Context, _ app.Event) {
	g.pickerQuery = ctx.JSSrc().Get("value").String()
	g.pickerCursor = 0
}

// pickerMatches is the filtered card pool the picker shows: every card whose name
// contains the query, in pool order. Query and name are both normalized, so the
// match ignores case and reads the Æmber ligature and its ASCII "ae" form as the
// same text. While the picker is answering a name-a-card prompt the pool is
// narrowed to the names that prompt offered.
func (g *game) pickerMatches() []engine.CardDefinition {
	q := normalizeSearch(strings.TrimSpace(g.pickerQuery))
	offered := g.pickerOfferedNames()
	var matches []engine.CardDefinition
	for _, d := range g.allDefs {
		if offered != nil && !offered[d.Name] {
			continue
		}
		if q == "" || strings.Contains(normalizeSearch(d.Name), q) {
			matches = append(matches, d)
		}
	}
	return matches
}

// pickerOfferedNames is the set of names a name-a-card prompt offered, or nil when
// the picker is adding a card to hand and every implemented card is pickable.
func (g *game) pickerOfferedNames() map[string]bool {
	if !g.pickerNaming {
		return nil
	}
	offered := make(map[string]bool, len(g.optionLabels))
	for _, label := range g.optionLabels {
		offered[label] = true
	}
	return offered
}

// movePickerCursor steps the highlighted row by delta, clamped to the list so the
// cursor never runs off either end.
func (g *game) movePickerCursor(delta int) {
	n := len(g.pickerMatches())
	if n == 0 {
		g.pickerCursor = 0
		return
	}
	g.pickerCursor += delta
	if g.pickerCursor < 0 {
		g.pickerCursor = 0
	}
	if g.pickerCursor >= n {
		g.pickerCursor = n - 1
	}
}

// onPickerKey handles the keys the card picker owns while it is open, so Tab,
// Enter and the arrows drive the list instead of the board behind it.
func (g *game) onPickerKey(ctx app.Context, key string, shift bool) {
	switch key {
	case "Escape":
		g.closePicker(ctx, app.Event{})
	case "Enter":
		g.addCursorCard(ctx)
	case "Tab":
		if shift {
			g.movePickerCursor(-1)
		} else {
			g.movePickerCursor(1)
		}
	case "ArrowUp":
		g.movePickerCursor(-1)
	case "ArrowDown":
		g.movePickerCursor(1)
	}
}

// addCursorCard commits the highlighted picker row.
func (g *game) addCursorCard(ctx app.Context) {
	matches := g.pickerMatches()
	if g.pickerCursor < 0 || g.pickerCursor >= len(matches) {
		return
	}
	g.commitPickedCard(ctx, matches[g.pickerCursor])
}

// addPickedCard commits the clicked picker row's card. The card is named by the
// row's data attribute rather than captured per row: go-app compares handlers by
// pointer, so a closure per row goes stale the moment the search filters the list,
// and the click commits whatever card used to sit there.
func (g *game) addPickedCard(ctx app.Context, _ app.Event) {
	def, ok := g.defByName[ctx.JSSrc().Get("dataset").Get("card").String()]
	if !ok {
		return
	}
	g.commitPickedCard(ctx, *def)
}

// commitPickedCard resolves a picked row: it answers the name-a-card prompt the
// picker was opened for, or, in manual mode, adds the card to hand.
func (g *game) commitPickedCard(ctx app.Context, def engine.CardDefinition) {
	if g.pickerNaming {
		g.nameCard(ctx, def.Name)
		return
	}
	g.addCardDef(ctx, def)
}

// nameCard answers the open name-a-card prompt with the picked name and closes
// the picker, so the blocked effect resumes.
func (g *game) nameCard(ctx app.Context, name string) {
	for i, label := range g.optionLabels {
		if label != name {
			continue
		}
		g.pickerOpen, g.pickerNaming = false, false
		g.chooseOptionIdx(i)(ctx, app.Event{})
		return
	}
}

// addCardDef puts a card definition into the active player's hand and records the
// add so a reload can replay it, then closes the picker.
func (g *game) addCardDef(ctx app.Context, def engine.CardDefinition) {
	g.beginAction()
	player := g.active()
	if _, added := g.g.ManualAddCard(def, player); added {
		// Record the add so a reload replays the registration and the rebuilt
		// catalog hands out the same id the rest of the log refers to.
		g.record(input{Kind: inManualAddCard, Name: def.Name, Player: player})
	}
	g.pickerOpen = false
	g.save(ctx)
}

// isInPlay reports whether a card is on either player's battleline or artifact row.
func (g *game) isInPlay(id engine.LocalID) bool {
	for p := range 2 {
		if containsID(g.g.Battleline(p), id) || containsID(g.g.Artifacts(p), id) {
			return true
		}
	}
	return false
}
