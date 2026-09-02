package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the component's lifecycle and its page-level wiring: mounting and
// dismounting, the post-render scrolling, the keyboard shortcuts, and the
// hot-reload hand-off.

// OnMount resumes the saved match if there is one, else deals a fresh game. It
// runs on the UI goroutine once the component is inserted into the page.
func (g *game) OnMount(ctx app.Context) {
	g.dispatch = ctx.Dispatch
	if !g.resume(ctx) {
		g.newMatch()
	}
	g.inPlayPrev = g.inPlaySet()
	g.save(ctx)
	g.installKeyShortcuts()
	g.scrollLogToBottom()
}

// logScrollSlack is how far (in pixels) from the bottom of the log still counts
// as "following along", and so still auto-scrolls when new lines arrive.
const logScrollSlack = 48

// OnUpdate keeps the log pinned to its newest line, but only when the player was
// already at the bottom — someone who has scrolled back to read an earlier turn
// is not yanked forward every time a line is appended.
func (g *game) OnUpdate(app.Context) {
	g.flyIntoPlay()
	g.scrollPromptZoneIntoView()
	g.focusPickerInput()
	el := app.Window().GetElementByID("gamelog")
	if !el.Truthy() {
		return
	}
	top := el.Get("scrollTop").Float()
	height := el.Get("scrollHeight").Float()
	view := el.Get("clientHeight").Float()
	// scrollTop is still the pre-update position, so compare it against the height
	// recorded on the previous update to decide where the player was reading.
	if top+view >= g.logScrollHeight-logScrollSlack {
		el.Set("scrollTop", height)
	}
	g.logScrollHeight = height
}

// scrollPromptZoneIntoView brings the zone row a prompt is asking about into view
// the first time that prompt's viewer renders, so the player does not have to hunt
// down the pile they are being asked to click in.
func (g *game) scrollPromptZoneIntoView() {
	if g.promptZone == "" || g.promptZoneScrolled {
		return
	}
	el := app.Window().GetElementByID(promptZoneID)
	if !el.Truthy() {
		return
	}
	el.Call("scrollIntoView", map[string]any{"block": "center"})
	g.promptZoneScrolled = true
}

// focusPickerInput puts the caret in the manual card picker's search box the first
// time it renders, so the modal opens ready to be typed into.
func (g *game) focusPickerInput() {
	if !g.pickerOpen || g.pickerFocused {
		return
	}
	el := app.Window().GetElementByID(pickerInputID)
	if !el.Truthy() {
		return
	}
	el.Call("focus")
	g.pickerFocused = true
}

// scrollLogToBottom pins the log to its newest line, used when the board is first
// dealt or resumed so the log opens where the game currently is.
func (g *game) scrollLogToBottom() {
	el := app.Window().GetElementByID("gamelog")
	if !el.Truthy() {
		return
	}
	el.Set("scrollTop", el.Get("scrollHeight"))
	g.logScrollHeight = el.Get("scrollHeight").Float()
}

// installKeyShortcuts wires a document-level keydown listener so common actions
// have a single-key shortcut (see onKey). It listens on the document because the
// board has no single focused element to receive the keys.
func (g *game) installKeyShortcuts() {
	if g.keyFunc != nil {
		return
	}
	g.keyFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 {
			return nil
		}
		e := args[0]
		if isTextInput(e.Get("target")) || e.Get("altKey").Bool() {
			return nil
		}
		key := e.Get("key").String()
		if e.Get("ctrlKey").Bool() || e.Get("metaKey").Bool() {
			// Ctrl/Cmd+Z undo; Ctrl/Cmd+Shift+Z or Ctrl/Cmd+Y redo.
			switch key {
			case "z", "Z":
				if e.Get("shiftKey").Bool() {
					g.dispatch(func(ctx app.Context) { g.redoAction(ctx, app.Event{}) })
				} else {
					g.dispatch(func(ctx app.Context) { g.undoAction(ctx, app.Event{}) })
				}
			case "y", "Y":
				g.dispatch(func(ctx app.Context) { g.redoAction(ctx, app.Event{}) })
			}
			return nil
		}
		// Tab and the arrows move the selection, so the browser must not also move
		// focus or scroll the page with them.
		if navigates(key) {
			e.Call("preventDefault")
		}
		shift := e.Get("shiftKey").Bool()
		g.dispatch(func(ctx app.Context) { g.onKey(ctx, key, shift) })
		return nil
	})
	app.Window().Get("document").Call("addEventListener", "keydown", g.keyFunc)
}

// navigates reports whether a key moves or answers the selection, and so must be
// taken from the browser before it moves focus or scrolls instead. The home-row
// synonyms are plain letters the browser does nothing with, so only the real
// navigation keys and Space are claimed.
func navigates(key string) bool {
	switch key {
	case "Tab", "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", " ":
		return true
	}
	return false
}

// OnDismount removes the keyboard listener and frees its wrapped function.
func (g *game) OnDismount() {
	if g.keyFunc != nil {
		app.Window().Get("document").Call("removeEventListener", "keydown", g.keyFunc)
		g.keyFunc.Release()
		g.keyFunc = nil
	}
}

// isTextInput reports whether a key event's target is a text-entry element, so
// typing (e.g. in the manual card picker) is not hijacked by shortcuts.
func isTextInput(target app.Value) bool {
	if !target.Truthy() {
		return false
	}
	switch target.Get("tagName").String() {
	case "INPUT", "TEXTAREA":
		return true
	}
	return target.Get("isContentEditable").Truthy()
}

// onKey maps a keyboard shortcut to an action for the current selection and
// phase. Each action handler self-guards (busy, phase, selection), so a key in
// the wrong context is a harmless no-op. Escape, r, n and the navigation keys are
// handled first because backing out, answering, and moving between candidates
// must all work while a prompt blocks every other key.
func (g *game) onKey(ctx app.Context, key string, shift bool) {
	if g.g == nil {
		return
	}
	switch key {
	case "Escape":
		g.dismiss(ctx)
		return
	case "r":
		g.affirm(ctx)
		return
	case "n":
		g.deny(ctx)
		return
	case "Tab":
		if shift {
			g.tabSel(ctx, -1)
		} else {
			g.tabSel(ctx, 1)
		}
		return
	case "Enter", " ":
		g.confirmPrompt(ctx)
		return
	case "?":
		g.keysOpen = !g.keysOpen
		return
	}
	if g.busy || g.choosing || g.choosingOption {
		return
	}
	// 1-9 pick the nth card of the row the selection is in.
	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
		g.selectNth(ctx, int(key[0]-'0'))
		return
	}
	switch key {
	case "ArrowLeft", "j":
		g.moveSel(ctx, 0, -1)
	case "ArrowRight", ";":
		g.moveSel(ctx, 0, 1)
	case "ArrowUp":
		g.moveSel(ctx, -1, 0)
	case "ArrowDown", "k":
		g.moveSel(ctx, 1, 0)
	case "l":
		// l is the home-row "up", except while placing a creature, where the flanks
		// are the only thing to move between and l is the left one.
		if g.phase == phaseFlank {
			g.playFlank(true)(ctx, app.Event{})
			return
		}
		g.moveSel(ctx, -1, 0)
	case "p":
		g.play(ctx, app.Event{})
	case "d":
		g.discard(ctx, app.Event{})
	case "a":
		g.useAction(ctx, app.Event{})
	case "f":
		g.startFight(ctx, app.Event{})
	case "u":
		// A stunned creature's only use is to shed the stun, which the engine spends
		// through the same reap path.
		if g.selKind == selYourCreature && g.g.Stunned(g.sel) {
			g.reap(ctx, app.Event{})
		}
	case "z", "Z":
		g.cycleZones()
	case "h":
		g.toggleSidebar(ctx, app.Event{})
	case "m":
		g.toggleManual(ctx, app.Event{})
	case "e":
		// A second e confirms an armed end-turn, so ending a turn with moves left is
		// e e rather than a second key to remember.
		g.endTurn(ctx, app.Event{})
	}
}

// affirm is the one "yes, do it" key (r): it takes the affirmative move for
// whatever is in front of the player — answering a yes/no prompt, taking the
// right flank while placing a creature, or the selected card's main use (play
// from hand, an artifact's action, otherwise reap). Each handler self-guards, so
// a press with nothing to affirm is a no-op.
func (g *game) affirm(ctx app.Context) {
	if g.choosingOption {
		// Only a yes/no prompt has an affirmative answer; a list of alternatives has
		// no option that "yes" could mean.
		if len(g.optionLabels) > 0 && g.optionLabels[0] == "Yes" {
			g.chooseOptionIdx(0)(ctx, app.Event{})
		}
		return
	}
	if g.busy || g.choosing {
		return
	}
	switch {
	case g.phase == phaseFlank:
		g.playFlank(false)(ctx, app.Event{})
	case g.selKind == selHand:
		g.play(ctx, app.Event{})
	case g.selKind == selYourArtifact:
		g.useAction(ctx, app.Event{})
	default:
		g.reap(ctx, app.Event{})
	}
}

// deny is the "no" key (n): it answers a yes/no prompt in the negative and
// passes on a prompt the player may decline. Anything else is Escape's job, so a
// press with nothing to refuse is a no-op.
func (g *game) deny(ctx app.Context) {
	if g.choosingOption {
		for i, label := range g.optionLabels {
			if label == "No" {
				g.chooseOptionIdx(i)(ctx, app.Event{})
				return
			}
		}
		return
	}
	if g.choosing && g.chooserDeclinable {
		g.declineChooser(ctx, app.Event{})
	}
}

// dismiss backs out of whatever is currently open, innermost first: an overlay,
// then a pending confirmation, then a mid-action targeting step, then the
// selection. It backs out one layer per press so Escape never does more than the
// player expects.
func (g *game) dismiss(ctx app.Context) {
	switch {
	case g.keysOpen:
		g.keysOpen = false
	case g.menuOpen:
		g.menuOpen = false
	case g.pickerOpen:
		g.pickerOpen = false
	case g.zonesPlayer >= 0 && g.promptZone == "":
		g.zonesPlayer = -1
	case g.confirmRestart:
		g.confirmRestart = false
	case g.forgingKey >= 0:
		g.forgingKey = -1
	case g.choosing:
		// An optional prompt is declined; otherwise only a manual-mode prompt is
		// escapable, since a real chooser is mandatory.
		if g.chooserDeclinable {
			g.declineChooser(ctx, app.Event{})
		} else if g.g.Manual() {
			g.cancelChooser(ctx, app.Event{})
		}
	case g.phase == phaseFlank || g.phase == phaseFightTarget:
		g.cancelTargeting(ctx, app.Event{})
	case g.confirmEndTurn:
		g.confirmEndTurn = false
	case g.hasSel:
		g.clearSelection()
	}
}

// cycleZones steps the out-of-play zone viewer on by one press: your zones, then
// your opponent's, then closed. One key walks every pile, so no shortcut has to
// be remembered per player.
func (g *game) cycleZones() {
	switch g.zonesPlayer {
	case -1:
		g.zonesPlayer = g.active()
	case g.active():
		g.zonesPlayer = 1 - g.active()
	default:
		g.zonesPlayer = -1
	}
}

// OnAppUpdate fires when go-app detects a freshly built wasm bundle. It persists
// the current match and reloads the page onto the new build, which OnMount then
// resumes — the hot-reload path that keeps the game state.
func (g *game) OnAppUpdate(ctx app.Context) {
	g.save(ctx)
	ctx.Reload()
}
