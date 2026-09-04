package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is the component's lifecycle and its page-level wiring: mounting and
// dismounting, the post-render scrolling, the keyboard shortcuts, and the
// hot-reload hand-off.

// OnMount resumes the saved match if there is one, else deals a fresh game. It
// runs on the UI goroutine once the component is inserted into the page.
func (g *game) OnMount(ctx app.Context) {
	// go-app only calls a component's OnUpdate when its *parent* re-renders it
	// with changed exported fields; a root route like this one flags and
	// re-renders itself, so that path never runs and OnUpdate would otherwise
	// never fire. Deferring it here — after every dispatch, once the DOM has
	// actually been patched — is what actually invokes it.
	g.dispatch = func(fn func(app.Context)) {
		ctx.Dispatch(func(c app.Context) {
			if fn != nil {
				fn(c)
			}
			c.Defer(g.OnUpdate)
		})
	}
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
	g.scrollCursorIntoView()
	g.focusPickerInput()
	// Measured after the cursor scroll, so the rect is the card's resting place
	// rather than wherever it was on the way there.
	if g.measureFocus() {
		g.dispatch(nil)
	}
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

// scrollCursorIntoView keeps whatever Tab is currently pointed at — a card
// prompt's or fight target's cursor, or in ordinary play the selection — inside
// its row's scrolled strip. It is derived from the cursor's position rather than
// fired on a Tab keypress, so it is safe to run on every render: "nearest" only
// moves the scroll when the card is not already visible, so it never fights a
// player who has scrolled a strip on their own.
func (g *game) scrollCursorIntoView() {
	id, ok := g.cursorCardID()
	if !ok {
		return
	}
	el := app.Window().GetElementByID(id)
	if !el.Truthy() {
		return
	}
	el.Call("scrollIntoView", map[string]any{"block": "nearest", "inline": "nearest"})
}

// cursorCardID is the DOM id of the card Tab is currently pointed at.
func (g *game) cursorCardID() (string, bool) {
	if cands, ok := g.tabCandidates(); ok {
		if !g.hasCursor || !containsID(cands, g.promptCursor) {
			return "", false
		}
		return g.cardDOMID(g.promptCursor), true
	}
	if !g.hasSel {
		return "", false
	}
	return g.cardDOMID(g.sel), true
}

// cardDOMID is the DOM id a card renders under, in whichever zone currently
// holds it.
func (g *game) cardDOMID(id engine.LocalID) string {
	if containsID(g.g.Hand(g.active()), id) {
		return handCardID(id)
	}
	return boardCardID(id)
}

// measureFocus records where the selected card sits on screen and how big the
// window is, which is all the lifted copy of it needs to place itself. It reports
// whether anything moved, which is what tells OnUpdate to render again — and, since
// the copy is on its own layer and so reflows nothing, the render after that
// measures the same numbers and the pair settles. It is also called from the
// selection itself, so the common case places the copy on its first render rather
// than a frame later.
func (g *game) measureFocus() bool {
	id, ok := g.focusCardID()
	if !ok {
		was := g.hasFocus
		g.hasFocus = false
		return was
	}
	el := app.Window().GetElementByID(g.cardDOMID(id))
	if !el.Truthy() {
		return false
	}
	r := el.Call("getBoundingClientRect")
	next := cardRect{
		x: r.Get("left").Float(),
		y: r.Get("top").Float(),
		w: r.Get("width").Float(),
		h: r.Get("height").Float(),
	}
	// A card mid-flight or mid-collapse has no size to grow from, and dividing the
	// grow by it would be a division by zero.
	if next.w <= 0 || next.h <= 0 {
		return false
	}
	vw := app.Window().Get("innerWidth").Float()
	vh := app.Window().Get("innerHeight").Float()
	if g.hasFocus && g.focusID == id && g.focusRect == next &&
		g.focusViewW == vw && g.focusViewH == vh {
		return false
	}
	if !g.hasFocus || g.focusID != id {
		g.focusParity = !g.focusParity
	}
	g.focusRect, g.focusID, g.hasFocus = next, id, true
	g.focusViewW, g.focusViewH = vw, vh
	return true
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
			// Ctrl/Cmd+Z undo; Ctrl/Cmd+Shift+Z redo.
			if key == "z" || key == "Z" {
				if e.Get("shiftKey").Bool() {
					g.dispatch(func(ctx app.Context) { g.redoAction(ctx, app.Event{}) })
				} else {
					g.dispatch(func(ctx app.Context) { g.undoAction(ctx, app.Event{}) })
				}
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

// OnResize re-places the lifted card copy, which is positioned from a measurement
// of the board underneath it: without this, resizing the window leaves the copy
// floating over wherever its card used to be.
func (g *game) OnResize(app.Context) {
	if g.measureFocus() {
		g.dispatch(nil)
	}
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
		// While placing a creature r commits the right flank on its own, the same
		// as l for the left — a deliberate press should not need a second key to
		// confirm it. Otherwise r just reaps: it is not a general "yes" (see Space).
		if g.phase == phaseFlank {
			g.playFlank(false)(ctx, app.Event{})
			return
		}
		if g.chooseKeyColorKey(ctx, engine.KeyColorRed) {
			return
		}
		g.keyboardAction = true
		g.reap(ctx, app.Event{})
		return
	case "n":
		g.deny(ctx)
		return
	case "b":
		g.chooseKeyColorKey(ctx, engine.KeyColorBlue)
		return
	case "y":
		g.chooseKeyColorKey(ctx, engine.KeyColorYellow)
		return
	case "Tab":
		if shift {
			g.tabSel(ctx, -1)
		} else {
			g.tabSel(ctx, 1)
		}
		return
	case "Enter":
		g.confirmPrompt(ctx)
		return
	case " ":
		// Space confirms the Tab cursor like Enter, but when there is none to
		// confirm it falls back to affirm's default answer instead of doing nothing.
		if !g.confirmPrompt(ctx) {
			g.affirm(ctx)
		}
		return
	case "?":
		g.keysOpen = !g.keysOpen
		return
	case "h":
		// Showing and hiding the log mutates no game state, so it races with nothing
		// and belongs with the other view-only keys, above the guard below.
		g.toggleSidebar(ctx, app.Event{})
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
		g.keyboardAction = true
		g.play(ctx, app.Event{})
	case "d":
		g.keyboardAction = true
		g.discard(ctx, app.Event{})
	case "a":
		g.keyboardAction = true
		g.useAction(ctx, app.Event{})
	case "f":
		g.startFight(ctx, app.Event{})
	case "u":
		g.keyboardAction = true
		g.unstun(ctx, app.Event{})
	case "z", "Z":
		g.cycleZones()
	case "m":
		g.toggleManual(ctx, app.Event{})
	case "e":
		// A second e confirms an armed end-turn, so ending a turn with moves left is
		// e e rather than a second key to remember.
		g.endTurn(ctx, app.Event{})
	}
}

// affirm is Space's fallback when there is no Tab cursor to confirm: it takes
// the affirmative move for whatever is in front of the player — answering a
// yes/no prompt, or the selected card's main use (play from hand, an
// artifact's action, otherwise reap). Each handler self-guards, so a press
// with nothing to affirm is a no-op.
func (g *game) affirm(ctx app.Context) {
	if g.choosingOption {
		// Only a yes/no prompt has an affirmative answer; a list of alternatives has
		// no option that "yes" could mean.
		if len(g.optionLabels) > 0 && g.optionLabels[0] == "Yes" {
			g.chooseOptionIdx(0)(ctx, app.Event{})
		}
		return
	}
	if g.busy || g.choosing || g.phase == phaseFlank {
		return
	}
	switch g.selKind {
	case selHand:
		g.keyboardAction = true
		g.play(ctx, app.Event{})
	case selYourArtifact:
		g.keyboardAction = true
		g.useAction(ctx, app.Event{})
	default:
		g.keyboardAction = true
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
