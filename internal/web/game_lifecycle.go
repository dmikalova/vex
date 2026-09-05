package web

import (
	"math"
	"time"

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
	g.installScrollTracking()
	g.installSwipeGestures()
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
	g.scrollUsableRowsIntoView()
	g.focusPickerInput()
	g.refreshToast()
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

// scrollUsableRowsIntoView, when the end-turn confirm arms, scrolls each row that
// still holds a usable card so at least one jiggling card in that row is in view.
// The confirm warns that moves are left; a card the player has scrolled past is
// pulled back into its strip so the warning points at something they can see. It
// fires once per arming: scrolling every render would fight a player who then
// scrolls away, so once the rows have been revealed it waits for the next confirm.
func (g *game) scrollUsableRowsIntoView() {
	if !g.confirmEndTurn {
		g.confirmScrolled = false
		return
	}
	if g.confirmScrolled {
		return
	}
	g.confirmScrolled = true
	p := g.active()
	// The first jiggling card in a strip is scrolled into view, which carries its
	// whole strip's scroll along, so one card per row satisfies the confirm.
	rows := []struct {
		ids  []engine.LocalID
		kind selKind
		hand bool
	}{
		{g.g.Battleline(p), selYourCreature, false},
		{g.g.Artifacts(p), selYourArtifact, false},
		{g.g.Hand(p), selHand, true},
	}
	for _, row := range rows {
		for _, id := range row.ids {
			if !g.jiggling(id, row.kind) {
				continue
			}
			domID := boardCardID(id)
			if row.hand {
				domID = handCardID(id)
			}
			if el := app.Window().GetElementByID(domID); el.Truthy() {
				el.Call("scrollIntoView", map[string]any{"block": "nearest", "inline": "nearest"})
			}
			break
		}
	}
}

// scrollCursorIntoView brings whatever Tab is currently pointed at — a card
// prompt's or fight target's cursor, or in ordinary play the selection — inside
// its row's scrolled strip. It fires only when the cursor lands on a new card:
// running on every render, it would haul a strip back every time the player
// scrolled away from the selected card, and the player and the client would take
// turns fighting over the scroll position.
func (g *game) scrollCursorIntoView() {
	id, ok := g.cursorCardID()
	if !ok {
		g.cursorScrolled = ""
		return
	}
	if id == g.cursorScrolled {
		return
	}
	el := app.Window().GetElementByID(id)
	if !el.Truthy() {
		return
	}
	el.Call("scrollIntoView", map[string]any{"block": "nearest", "inline": "nearest"})
	g.cursorScrolled = id
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

// toastLinger is how long a minimized-log toast bubble stays up before it clears
// itself.
const toastLinger = 5 * time.Second

// refreshToast surfaces log lines the player would otherwise miss: while the
// sidebar (and its log) is hidden, new lines since the last catch-up group into
// the same bubbles the panel draws and toast over the board. With the sidebar
// open the log itself is on screen, so nothing toasts and the catch-up simply
// tracks the log. It runs after each render (from OnUpdate); once caught up it
// dispatches nothing, so it cannot loop. Turn and phase headers do not toast — a
// bare scene break is not news — but they still close the open bubble.
func (g *game) refreshToast() {
	if !g.sidebarCollapsed {
		g.toastSeen = len(g.g.Log)
		g.toastBubbles = nil
		g.toastOpen = false
		return
	}
	if len(g.g.Log) <= g.toastSeen {
		return
	}
	starts := make(map[int]int, len(g.logGroups))
	for _, m := range g.logGroups {
		starts[m.Start] = m.Player
	}
	changed := false
	for i := g.toastSeen; i < len(g.g.Log); i++ {
		rec := g.g.Log[i]
		if rule, _ := ruleOf(rec); rule != ruleNone {
			g.toastOpen = false
			continue
		}
		if player, ok := starts[i]; ok {
			g.openToastBubble(player)
		} else if !g.toastOpen {
			g.openToastBubble(g.enclosingPlayer(i))
		}
		b := &g.toastBubbles[len(g.toastBubbles)-1]
		b.lines = append(b.lines, rec)
		g.armToastExpiry(b)
		changed = true
	}
	g.toastSeen = len(g.g.Log)
	if changed {
		g.dispatch(nil)
	}
}

// openToastBubble starts a fresh bubble for a new root action, so the toast keeps
// the same one-bubble-per-action grouping the log panel does.
func (g *game) openToastBubble(player int) {
	g.toastBubbles = append(g.toastBubbles, toastBubble{player: player})
	g.toastOpen = true
}

// enclosingPlayer is whose action a bubble that did not start on a group mark
// belongs to — the last group opened at or before that line.
func (g *game) enclosingPlayer(i int) int {
	player, at := -1, -1
	for _, m := range g.logGroups {
		if m.Start <= i && m.Start > at {
			at, player = m.Start, m.Player
		}
	}
	return player
}

// armToastExpiry (re)starts a bubble's countdown, freshening it on every new line
// so an action still resolving does not fade mid-way. When it fires it drops just
// that bubble, unless the pointer is holding the toast open, in which case it
// waits out another linger. A superseded timer finds no bubble with its id and
// does nothing.
func (g *game) armToastExpiry(b *toastBubble) {
	g.toastGen++
	b.gen = g.toastGen
	gen := g.toastGen
	time.AfterFunc(toastLinger, func() {
		g.dispatch(func(app.Context) {
			if g.toastHover || g.toastPinned {
				g.rearmToastExpiry(gen)
				return
			}
			g.dropToastBubble(gen)
		})
	})
}

// rearmToastExpiry keeps a held-open bubble alive: it arms a fresh countdown for
// the same bubble, so a paused toast never expires under the pointer.
func (g *game) rearmToastExpiry(gen int) {
	for i := range g.toastBubbles {
		if g.toastBubbles[i].gen == gen {
			g.armToastExpiry(&g.toastBubbles[i])
			return
		}
	}
}

// dropToastBubble removes the bubble whose countdown just fired. Dropping the
// newest bubble also closes the group, so the next line opens a fresh one.
func (g *game) dropToastBubble(gen int) {
	for i := range g.toastBubbles {
		if g.toastBubbles[i].gen == gen {
			if i == len(g.toastBubbles)-1 {
				g.toastOpen = false
			}
			g.toastBubbles = append(g.toastBubbles[:i], g.toastBubbles[i+1:]...)
			return
		}
	}
}

// pauseToast, resumeToast, and toggleToastPin keep the toast up while the player
// is reading it: hovering freezes every bubble's countdown, leaving lets them run
// again, and a click pins the whole toast open until the next click.
func (g *game) pauseToast(_ app.Context, _ app.Event)  { g.toastHover = true }
func (g *game) resumeToast(_ app.Context, _ app.Event) { g.toastHover = false }

func (g *game) toggleToastPin(_ app.Context, _ app.Event) { g.toastPinned = !g.toastPinned }

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

// installScrollTracking keeps the lifted card copy over the card it was lifted
// from while a strip is scrolled. The copy is fixed to the window, so without this
// scrolling a hand slides the card out from under its own copy. A scroll event
// does not bubble, so the listener is registered in the capture phase and thereby
// sees every strip at once.
func (g *game) installScrollTracking() {
	if g.scrollFunc != nil {
		return
	}
	g.scrollFunc = app.FuncOf(func(app.Value, []app.Value) any {
		g.placeFocus()
		return nil
	})
	app.Window().Get("document").Call("addEventListener", "scroll", g.scrollFunc, true)
}

// swipeEdgeBand is how far (in pixels) from the right edge a touch must begin for
// an open-swipe to count, so an edge drag reveals the sidebar without a swipe that
// starts mid-board doing the same.
const swipeEdgeBand = 32

// swipeMinDistance is the horizontal travel (in pixels) a swipe must cover before
// it toggles the sidebar, so a tap or a short drag does not move it.
const swipeMinDistance = 60

// installSwipeGestures wires document-level touch listeners so a horizontal swipe
// moves the sidebar on a touchscreen: a swipe that starts near the right edge and
// travels left reveals the sidebar, and a swipe that travels right hides it. It
// mirrors the » / « reveal buttons for a phone where the edge is easier to reach
// than the button. A mostly-vertical drag (scrolling a strip or the log) is left
// alone, and an open-swipe must begin in the edge band so a mid-board drag does
// not summon the drawer.
func (g *game) installSwipeGestures() {
	if g.touchStartFunc != nil {
		return
	}
	g.touchStartFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 {
			return nil
		}
		touches := args[0].Get("touches")
		if !touches.Truthy() || touches.Get("length").Int() != 1 {
			g.swipeTracking = false
			return nil
		}
		t := touches.Index(0)
		g.swipeStartX = t.Get("clientX").Float()
		g.swipeStartY = t.Get("clientY").Float()
		g.swipeTracking = true
		return nil
	})
	g.touchEndFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if !g.swipeTracking || len(args) == 0 {
			return nil
		}
		g.swipeTracking = false
		changed := args[0].Get("changedTouches")
		if !changed.Truthy() || changed.Get("length").Int() == 0 {
			return nil
		}
		t := changed.Index(0)
		dx := t.Get("clientX").Float() - g.swipeStartX
		dy := t.Get("clientY").Float() - g.swipeStartY
		// A mostly-vertical drag is a scroll, not a sidebar swipe.
		if math.Abs(dx) < swipeMinDistance || math.Abs(dy) > math.Abs(dx) {
			return nil
		}
		width := app.Window().Get("innerWidth").Float()
		switch {
		case dx < 0 && g.sidebarCollapsed && g.swipeStartX >= width-swipeEdgeBand:
			// Swipe left from the right edge: reveal the hidden sidebar.
			g.dispatch(func(ctx app.Context) { g.toggleSidebar(ctx, app.Event{}) })
		case dx > 0 && !g.sidebarCollapsed:
			// Swipe right: hide the sidebar out to the edge.
			g.dispatch(func(ctx app.Context) { g.toggleSidebar(ctx, app.Event{}) })
		}
		return nil
	})
	doc := app.Window().Get("document")
	doc.Call("addEventListener", "touchstart", g.touchStartFunc, map[string]any{"passive": true})
	doc.Call("addEventListener", "touchend", g.touchEndFunc, map[string]any{"passive": true})
}

// OnResize re-places the lifted card copy, which is positioned from a measurement
// of the board underneath it: without this, resizing the window leaves the copy
// floating over wherever its card used to be.
func (g *game) OnResize(app.Context) { g.placeFocus() }

// placeFocus re-measures the board beneath the lifted card copy and redraws when
// the copy no longer sits where its card does.
func (g *game) placeFocus() {
	if g.measureFocus() {
		g.dispatch(nil)
	}
}

// wheelOverFocus scrolls the strip the lifted card copy is drawn over, since the
// copy now takes the pointer (it has to, to be the source of its own drag) and
// would otherwise stop a wheel dead on the very card the player is reading.
func (g *game) wheelOverFocus(_ app.Context, e app.Event) {
	id, ok := g.focusCardID()
	if !ok {
		return
	}
	el := app.Window().GetElementByID(g.cardDOMID(id))
	if !el.Truthy() {
		return
	}
	strip := el.Call("closest", ".card-strip")
	if !strip.Truthy() {
		return
	}
	// A strip only scrolls sideways, and a plain mouse wheel only has a dy to give
	// it — which is the trade the browser itself makes over a horizontal scroller.
	d := e.Get("deltaX").Float()
	if d == 0 {
		d = e.Get("deltaY").Float()
	}
	strip.Set("scrollLeft", strip.Get("scrollLeft").Float()+d)
	e.PreventDefault()
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

// OnDismount removes the document listeners and frees their wrapped functions.
func (g *game) OnDismount() {
	if g.keyFunc != nil {
		app.Window().Get("document").Call("removeEventListener", "keydown", g.keyFunc)
		g.keyFunc.Release()
		g.keyFunc = nil
	}
	if g.scrollFunc != nil {
		app.Window().Get("document").Call("removeEventListener", "scroll", g.scrollFunc, true)
		g.scrollFunc.Release()
		g.scrollFunc = nil
	}
	if g.touchStartFunc != nil {
		app.Window().Get("document").Call("removeEventListener", "touchstart", g.touchStartFunc)
		g.touchStartFunc.Release()
		g.touchStartFunc = nil
	}
	if g.touchEndFunc != nil {
		app.Window().Get("document").Call("removeEventListener", "touchend", g.touchEndFunc)
		g.touchEndFunc.Release()
		g.touchEndFunc = nil
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
	// While the set picker is up every game key is inert; only Escape works, to
	// back out of it.
	if g.awaitingSetup {
		if key == "Escape" {
			g.dismiss(ctx)
		}
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
	case g.awaitingSetup:
		g.awaitingSetup = false
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
