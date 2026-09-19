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

// OnMount resumes the saved match if there is one, else opens the new-game set
// picker so a first-time load chooses its sets rather than silently dealing the
// base set. It runs on the UI goroutine once the component is inserted.
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
	if g.resume(ctx) {
		g.inPlayPrev = g.inPlaySet()
		g.save(ctx)
	} else {
		// Nothing to resume: open the set picker with no game dealt (g.g stays nil).
		// The player picks two sets and only then is the first match dealt, so a
		// fresh visit never assumes the base set.
		g.sidebarCollapsed = mobileViewport()
		g.beginSetup()
	}
	g.installKeyShortcuts()
	g.installScrollTracking()
	g.installSwipeGestures()
	g.installTips()
	g.scrollLogToBottom()
}

// mobileMaxWidth is the viewport width (in px) at or below which the view reads as
// mobile-sized: the 20rem sidebar would cover more than 40% of it, so it is a
// drawer over the board rather than a panel beside it. It matches sidebarTooWide's
// 40% rule (20rem / 0.4 = 50rem, at 16px/rem).
const mobileMaxWidth = 50 * 16.0

// mobileViewport reports whether the view is mobile-sized (horizontally thin), so
// a fresh load hides the sidebar by default and gives the board the width.
// Off-browser (no window to measure) it reports false, the desktop default.
func mobileViewport() bool {
	vw := app.Window().Get("innerWidth").Float()
	return vw > 0 && vw < mobileMaxWidth
}

// logScrollSlack is how far (in pixels) from the bottom of the log still counts
// as "following along", and so still auto-scrolls when new lines arrive.
const logScrollSlack = 48

// OnUpdate keeps the log pinned to its newest line, but only when the player was
// already at the bottom — someone who has scrolled back to read an earlier turn
// is not yanked forward every time a line is appended.
func (g *game) OnUpdate(app.Context) {
	// Before the first deal (the set picker on a fresh load) there is no game to
	// read: every helper below reaches into g.g, so bail until a match exists.
	if g.g == nil {
		return
	}
	g.flyIntoPlay()
	g.scrollPromptZoneIntoView()
	g.scrollCursorIntoView()
	g.scrollUsableRowsIntoView()
	g.focusPickerInput()
	g.refreshToast()
	g.clampOpenDeckList()
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
				el.Call(
					"scrollIntoView",
					map[string]any{
						"block":  "nearest",
						"inline": "nearest",
					},
				)
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
	el.Call(
		"scrollIntoView",
		map[string]any{
			"block":  "nearest",
			"inline": "nearest",
		},
	)
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
		// A selection dropped while the card is still at rest (not mid-play/use)
		// plays the copy out; an action-driven drop leaves the flight animation to
		// carry the card instead.
		if was && !g.busy && g.focusShown.id != 0 {
			g.startFocusExit()
		}
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
	// A live copy supersedes any exit that was still fading.
	g.focusExit = false
	if g.hasFocus && g.focusID == id && g.focusRect == next &&
		math.Abs(g.focusViewW-vw) < 0.01 && math.Abs(g.focusViewH-vh) < 0.01 {
		return false
	}
	if !g.hasFocus || g.focusID != id {
		g.focusParity = !g.focusParity
	}
	g.focusRect, g.focusID, g.hasFocus = next, id, true
	g.focusViewW, g.focusViewH = vw, vh
	g.focusShown = g.focusSnapshotNow(id)
	return true
}

// focusExitDur is how long the deselected copy takes to shrink back to its slot.
// It matches the grow-in it plays backwards (card-focus-out in app.css).
const focusExitDur = 150 * time.Millisecond

// startFocusExit begins the shrink-back of a just-deselected copy and arms a timer
// to clear it once the animation has run. The generation tag means a fresh
// selection during the fade (which clears focusExit itself) leaves this timer to
// find nothing to do.
func (g *game) startFocusExit() {
	g.focusExit = true
	g.focusExitGen++
	gen := g.focusExitGen
	time.AfterFunc(focusExitDur, func() {
		g.dispatch(func(app.Context) {
			if g.focusExitGen == gen {
				g.focusExit = false
			}
		})
	})
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

// toastLeave is how long a bubble's leave animation runs — a quarter second to
// fade, then a quarter second to collapse the space it held (see the
// .log-toast-item--leaving rule) — before it is taken out of the toast.
const toastLeave = 500 * time.Millisecond

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

// dropToastBubble starts the leave of the bubble whose countdown just fired: it
// fades and collapses in place (see toastLeave) before removeToastBubble takes it
// out. Dropping the newest bubble also closes the group, so the next line opens a
// fresh one. A bubble already leaving stays put.
func (g *game) dropToastBubble(gen int) {
	for i := range g.toastBubbles {
		if g.toastBubbles[i].gen == gen {
			if g.toastBubbles[i].leaving {
				return
			}
			if i == len(g.toastBubbles)-1 {
				g.toastOpen = false
			}
			g.toastBubbles[i].leaving = true
			time.AfterFunc(toastLeave, func() {
				g.dispatch(func(app.Context) { g.removeToastBubble(gen) })
			})
			return
		}
	}
}

// removeToastBubble takes a faded-out bubble out of the toast once its leave
// animation has run.
func (g *game) removeToastBubble(gen int) {
	for i := range g.toastBubbles {
		if g.toastBubbles[i].gen == gen {
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

func (g *game) toggleToastPin(
	_ app.Context,
	_ app.Event,
) {
	g.toastPinned = !g.toastPinned
}

// dismissToast clears the toast outright, catching the seen floor up to the log
// so its lines do not toast again. It stops the click from also toggling the pin.
func (g *game) dismissToast(_ app.Context, e app.Event) {
	e.Call("stopPropagation")
	g.clearToast()
}

// clearToast drops every toast bubble and catches the seen floor up to the log so
// the cleared lines do not toast again.
func (g *game) clearToast() {
	g.toastBubbles = nil
	g.toastOpen = false
	g.toastPinned = false
	g.toastSeen = len(g.g.Log)
}

// installKeyShortcuts wires a document-level keydown listener so common actions
// have a single-key shortcut (see onKey). It listens on the document because the
// board has no single focused element to receive the keys. Which group of
// shortcuts is live depends on what is on screen: the card picker and the
// new-game set picker each take the keys while they are up, and the board takes
// the rest.
func (g *game) installKeyShortcuts() {
	if g.keyFunc != nil {
		return
	}
	g.keyFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 {
			return nil
		}
		e := args[0]
		switch {
		case g.pickerOpen:
			g.pickerKey(e)
		case g.awaitingSetup:
			g.setupKey(e)
		case isTextInput(e.Get("target")) || e.Get("altKey").Bool():
			// Typing in a field, or holding Alt for a browser menu: not ours.
		default:
			g.boardKey(e)
		}
		return nil
	})
	app.Window().Get("document").Call("addEventListener", "keydown", g.keyFunc)
}

// pickerKey handles the keys the card picker owns while it is open, even with the
// caret in its search box (where boardKey's text-input guard would otherwise skip
// them), so the list is driven from the keyboard and Tab cannot escape the modal
// into the board behind it. Every other key is swallowed.
func (g *game) pickerKey(e app.Value) {
	key := e.Get("key").String()
	switch key {
	case "Tab", "Enter", "ArrowUp", "ArrowDown", "Escape":
		e.Call("preventDefault")
		shift := e.Get("shiftKey").Bool()
		g.dispatch(func(ctx app.Context) { g.onPickerKey(ctx, key, shift) })
	}
}

// setupKey handles the keys the new-game set picker takes. The picker is a strip
// of ordinary buttons, so the browser owns Tab (move focus between sets) and
// Enter/Space (activate the focused one). Only the shortcut sheet (?) and backing
// out (Escape) are taken here; everything else falls through untouched so native
// focus still works.
func (g *game) setupKey(e app.Value) {
	switch e.Get("key").String() {
	case "Escape":
		// Close the shortcut sheet if it is up; otherwise back out of the picker,
		// but only when there is a game behind it to return to (not on first load).
		if g.keysOpen {
			g.dispatch(func(ctx app.Context) { g.closeKeys(ctx, app.Event{}) })
		} else if g.g != nil {
			g.dispatch(func(ctx app.Context) { g.dismiss(ctx) })
		}
	case "?":
		g.dispatch(func(ctx app.Context) { g.toggleKeys(ctx, app.Event{}) })
	}
}

// boardKey handles the shortcuts that drive the board itself: the Ctrl/Cmd
// history and new-game chords, then the plain single keys onKey interprets.
func (g *game) boardKey(e app.Value) {
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
		// Ctrl/Cmd+G opens the new-game set picker.
		if key == "g" || key == "G" {
			e.Call("preventDefault")
			g.dispatch(func(ctx app.Context) { g.openSetup(ctx, app.Event{}) })
		}
		return
	}
	// Tab and the arrows move the selection, so the browser must not also move
	// focus or scroll the page with them.
	if navigates(key) {
		e.Call("preventDefault")
	}
	shift := e.Get("shiftKey").Bool()
	g.dispatch(func(ctx app.Context) { g.onKey(ctx, key, shift) })
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
	app.Window().
		Get("document").
		Call("addEventListener", "scroll", g.scrollFunc, true)
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
// alone, an open-swipe must begin in the edge band so a mid-board drag does not
// summon the drawer, and a swipe that begins on a horizontally-scrollable card row
// never opens the sidebar so scrolling a row of creatures is not mistaken for one.
func (g *game) installSwipeGestures() {
	if g.touchStartFunc != nil {
		return
	}
	g.touchStartFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 {
			return nil
		}
		// Any touch marks a touchscreen, so log mentions open on one tap rather than
		// fighting the synthetic mouseenter (see onLogCardHover).
		g.isTouch = true
		touches := args[0].Get("touches")
		if !touches.Truthy() || touches.Get("length").Int() != 1 {
			g.swipeTracking = false
			return nil
		}
		t := touches.Index(0)
		g.swipeStartX = t.Get("clientX").Float()
		g.swipeStartY = t.Get("clientY").Float()
		g.swipeTracking = true
		target := args[0].Get("target")
		// A press that begins on the player bar drives its stat tooltip (see
		// installTips), so it must not also swipe the sidebar.
		if target.Truthy() && target.Call("closest", ".score-pill").Truthy() {
			g.swipeTracking = false
			return nil
		}
		// Note whether the touch began on the toast, so its end can flick it away
		// instead of moving the sidebar.
		g.toastSwipeStart = target.Truthy() &&
			target.Call("closest", ".log-toast").Truthy()
		// Note whether the touch began inside a horizontally-scrollable card row, so
		// scrolling a row of creatures or artifacts is not mistaken for an open-swipe.
		g.swipeOnStrip = target.Truthy() &&
			target.Call("closest", ".card-strip").Truthy()
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
		// A horizontal swipe that began on the toast flicks it away, either way.
		if g.toastSwipeStart && len(g.toastBubbles) > 0 {
			g.dispatch(func(app.Context) { g.clearToast() })
			return nil
		}
		width := app.Window().Get("innerWidth").Float()
		switch {
		case dx < 0 && g.sidebarCollapsed && g.swipeStartX >= width-swipeEdgeBand && !g.swipeOnStrip:
			// Swipe left from the right edge: reveal the hidden sidebar. A swipe that
			// began on a scrollable card row is left to scroll that row instead.
			g.dispatch(
				func(ctx app.Context) { g.toggleSidebar(ctx, app.Event{}) },
			)
		case dx > 0 && !g.sidebarCollapsed:
			// Swipe right: hide the sidebar out to the edge.
			g.dispatch(
				func(ctx app.Context) { g.toggleSidebar(ctx, app.Event{}) },
			)
		}
		return nil
	})
	doc := app.Window().Get("document")
	doc.Call(
		"addEventListener",
		"touchstart",
		g.touchStartFunc,
		map[string]any{"passive": true},
	)
	doc.Call(
		"addEventListener",
		"touchend",
		g.touchEndFunc,
		map[string]any{"passive": true},
	)
}

// tipGap is the space, in pixels, between a floating tip label and the element it
// names.
const tipGap = 4.0

// installTips shows a single floating label for whatever element the pointer is
// over that carries a data-tip. One fixed element (#tip-float, in Render) is
// filled from that element's tip, so it escapes the player bar's overflow clip —
// which a per-element CSS ::after bubble could not — and works for both a mouse
// hover and a finger dragged along the bar. A touch press inside a .score-pill
// also suppresses the sidebar edge-swipe. Positioning is done straight on the DOM
// rather than through a render, so following the pointer costs no re-render.
func (g *game) installTips() {
	if g.tipDownFunc != nil {
		return
	}
	doc := app.Window().Get("document")
	g.installHoverTips(doc)
	g.installTouchTips(doc)
	passive := map[string]any{"passive": true}
	doc.Call("addEventListener", "pointerover", g.tipOverFunc, passive)
	doc.Call("addEventListener", "pointerout", g.tipOutFunc, passive)
	doc.Call("addEventListener", "pointerdown", g.tipDownFunc, passive)
	doc.Call("addEventListener", "pointermove", g.tipMoveFunc, passive)
	doc.Call("addEventListener", "pointerup", g.tipUpFunc, passive)
	doc.Call("addEventListener", "pointercancel", g.tipUpFunc, passive)
	g.installSelCursor(doc, passive)
}

// showTip fills the floating label from el's data-tip and places it centred above
// el — flipping below when it would leave the top of the window, and clamped so
// it never runs off either side.
func showTip(doc, el app.Value) {
	tip := el.Get("dataset").Get("tip")
	if !tip.Truthy() {
		return
	}
	float := doc.Call("getElementById", "tip-float")
	if !float.Truthy() {
		return
	}
	float.Set("textContent", tip.String())
	float.Get("classList").Call("add", "tip-float--on")
	style := float.Get("style")
	r := el.Call("getBoundingClientRect")
	mid := r.Get("left").Float() + r.Get("width").Float()/2
	style.Set("left", px(mid))
	style.Set("top", px(r.Get("top").Float()-tipGap))
	style.Set("transform", "translate(-50%, -100%)")
	// Now that the label has a measured box, clamp it inside the window.
	const margin = 8.0
	fr := float.Call("getBoundingClientRect")
	vw := app.Window().Get("innerWidth").Float()
	if left := fr.Get("left").Float(); left < margin {
		mid += margin - left
		style.Set("left", px(mid))
	} else if right := fr.Get("right").Float(); right > vw-margin {
		mid -= right - (vw - margin)
		style.Set("left", px(mid))
	}
	if fr.Get("top").Float() < margin {
		style.Set("top", px(r.Get("bottom").Float()+tipGap))
		style.Set("transform", "translate(-50%, 0)")
	}
}

// hideTip lowers the floating label.
func hideTip(doc app.Value) {
	if float := doc.Call("getElementById", "tip-float"); float.Truthy() {
		float.Get("classList").Call("remove", "tip-float--on")
	}
}

// tipUnder returns the data-tip element under a viewport point, or a null value
// when the point is off every tip.
func tipUnder(doc app.Value, x, y float64) app.Value {
	el := doc.Call("elementFromPoint", x, y)
	if !el.Truthy() {
		return app.Null()
	}
	return el.Call("closest", "[data-tip]")
}

// installHoverTips builds the mouse-hover pair. Touch drives installTouchTips
// instead, so this pair ignores a non-mouse pointer. A pointerout whose pointer is
// still on a tip (moving onto a child, or across to the next tip) leaves the label
// up — the matching pointerover re-places it — so only leaving tips entirely hides.
func (g *game) installHoverTips(doc app.Value) {
	g.tipOverFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 || args[0].Get("pointerType").String() != "mouse" {
			return nil
		}
		if t := args[0].Get("target"); t.Truthy() {
			if el := t.Call("closest", "[data-tip]"); el.Truthy() {
				showTip(doc, el)
			}
		}
		return nil
	})
	g.tipOutFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 || args[0].Get("pointerType").String() != "mouse" {
			return nil
		}
		if rel := args[0].Get("relatedTarget"); rel.Truthy() &&
			rel.Call("closest", "[data-tip]").Truthy() {
			return nil
		}
		hideTip(doc)
		return nil
	})
}

// installTouchTips builds the press-drag-release trio, which keeps the label under
// a finger dragged along the bar for as long as the press lasts.
func (g *game) installTouchTips(doc app.Value) {
	g.tipDownFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 {
			return nil
		}
		g.closeDeckOnTapAway(args[0])
		tip := tipUnder(
			doc,
			args[0].Get("clientX").Float(),
			args[0].Get("clientY").Float(),
		)
		if !tip.Truthy() {
			return nil
		}
		g.tipTracking = true
		showTip(doc, tip)
		return nil
	})
	g.tipMoveFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if !g.tipTracking || len(args) == 0 {
			return nil
		}
		tip := tipUnder(
			doc,
			args[0].Get("clientX").Float(),
			args[0].Get("clientY").Float(),
		)
		if !tip.Truthy() {
			hideTip(doc)
			return nil
		}
		showTip(doc, tip)
		return nil
	})
	g.tipUpFunc = app.FuncOf(func(_ app.Value, _ []app.Value) any {
		if !g.tipTracking {
			return nil
		}
		g.tipTracking = false
		hideTip(doc)
		return nil
	})
}

// closeDeckOnTapAway unpins a deck list when the press lands outside every deck
// icon, so a touchscreen can dismiss it the way a click-away or hover-out would.
func (g *game) closeDeckOnTapAway(e app.Value) {
	if g.deckOpen == ([2]bool{}) {
		return
	}
	if t := e.Get("target"); t.Truthy() && t.Call("closest", ".deck-tip").Truthy() {
		return
	}
	g.deckOpen = [2]bool{}
	if g.dispatch != nil {
		g.dispatch(nil)
	}
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
		app.Window().
			Get("document").
			Call("removeEventListener", "keydown", g.keyFunc)
		g.keyFunc.Release()
		g.keyFunc = nil
	}
	if g.scrollFunc != nil {
		app.Window().
			Get("document").
			Call("removeEventListener", "scroll", g.scrollFunc, true)
		g.scrollFunc.Release()
		g.scrollFunc = nil
	}
	if g.touchStartFunc != nil {
		app.Window().
			Get("document").
			Call("removeEventListener", "touchstart", g.touchStartFunc)
		g.touchStartFunc.Release()
		g.touchStartFunc = nil
	}
	if g.touchEndFunc != nil {
		app.Window().
			Get("document").
			Call("removeEventListener", "touchend", g.touchEndFunc)
		g.touchEndFunc.Release()
		g.touchEndFunc = nil
	}
	if g.tipDownFunc != nil {
		doc := app.Window().Get("document")
		doc.Call("removeEventListener", "pointerover", g.tipOverFunc)
		doc.Call("removeEventListener", "pointerout", g.tipOutFunc)
		doc.Call("removeEventListener", "pointerdown", g.tipDownFunc)
		doc.Call("removeEventListener", "pointermove", g.tipMoveFunc)
		doc.Call("removeEventListener", "pointerup", g.tipUpFunc)
		doc.Call("removeEventListener", "pointercancel", g.tipUpFunc)
		g.tipOverFunc.Release()
		g.tipOutFunc.Release()
		g.tipDownFunc.Release()
		g.tipMoveFunc.Release()
		g.tipUpFunc.Release()
		g.tipOverFunc, g.tipOutFunc = nil, nil
		g.tipDownFunc, g.tipMoveFunc, g.tipUpFunc = nil, nil, nil
	}
	if g.selCursorFunc != nil {
		app.Window().Get("document").
			Call("removeEventListener", "pointermove", g.selCursorFunc)
		g.selCursorFunc.Release()
		g.selCursorFunc = nil
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

// deny is the "no" key (n): it answers a yes/no prompt in the negative, sheds an
// opening hand at the mulligan prompt, and passes on a prompt the player may
// decline. Anything else is Escape's job, so a press with nothing to refuse is a
// no-op.
func (g *game) deny(ctx app.Context) {
	if g.choosingOption {
		for i, label := range g.optionLabels {
			if isDecliningOption(label) {
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
	// A picker answering a name-a-card prompt is not dismissible: the blocked effect
	// is waiting on a name, so Escape would leave the action stuck with nothing on
	// screen to answer it.
	case g.pickerOpen && !g.pickerNaming:
		g.pickerOpen = false
	case g.zonesPlayer >= 0:
		g.zonesPlayer = -1
	case g.awaitingSetup:
		g.awaitingSetup = false
	case g.forgingKey >= 0:
		g.forgingKey = -1
	case g.inspecting:
		// A peek lift is dropped ahead of the prompt it was raised over, so the first
		// Escape puts the card down and the next backs out of the prompt itself.
		g.clearSelection()
	case g.choosing:
		// In manual mode Escape cancels any prompt, backing the whole action out;
		// otherwise only an optional prompt is escapable, by declining it, since a
		// real chooser is mandatory.
		if g.g.Manual() {
			g.cancelChooser(ctx, app.Event{})
		} else if g.chooserDeclinable {
			g.declineChooser(ctx, app.Event{})
		}
	case g.choosingOption:
		// An option prompt has no decline; only manual mode can back out of it.
		if g.g.Manual() {
			g.cancelChooser(ctx, app.Event{})
		}
	case g.phase == phaseFlank || g.phase == phaseFightTarget:
		g.cancelTargeting(ctx, app.Event{})
	case g.hostTargeting:
		g.cancelHostTargeting(ctx, app.Event{})
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
