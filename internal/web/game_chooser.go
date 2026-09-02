package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file bridges the engine's synchronous Chooser to the browser: a prompt
// raised while an effect resolves parks that goroutine, the UI renders the
// question, and the player's click sends the answer back.

// chooseReply carries the player's answer to an engine chooser request: the
// chosen id, or ok=false when the player cancels.
type chooseReply struct {
	id engine.LocalID
	ok bool
}

// webChooser adapts the engine's synchronous Chooser to go-app's single UI
// goroutine. The engine calls ChooseCreature from a background action goroutine
// (see game.runAction); it shows the chooser overlay on the UI goroutine, then
// blocks until a candidate is clicked (or the request is cancelled).
type webChooser struct {
	g           *game
	reply       chan chooseReply
	optionReply chan int
}

// ChooseCreature posts a chooser request to the UI and waits for the player's
// pick. It returns false when there are no candidates or the player cancels.
func (c *webChooser) ChooseCreature(
	source, prompt string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	return c.ask(source, prompt, candidates, false)
}

// ChooseCardOrDecline implements the engine's DeclinableChooser: an optional card
// choice ("you may destroy another friendly creature", "exhaust up to 3
// creatures") is shown as the same highlighted-card prompt, plus a Done button.
// Without it the engine falls back to a list of card names, which made the player
// read a menu instead of clicking the card in front of them.
func (c *webChooser) ChooseCardOrDecline(
	source, prompt string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	return c.ask(source, prompt, candidates, true)
}

// ask is the shared card-prompt path: it shows the prompt on the UI goroutine and
// blocks the action goroutine until a candidate is clicked, the prompt is
// declined, or it is cancelled.
func (c *webChooser) ask(
	source, prompt string,
	candidates []engine.LocalID,
	declinable bool,
) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	// Discard any stale reply left in the buffer (e.g. from a double click on the
	// previous prompt) so it cannot silently answer this one.
	select {
	case <-c.reply:
	default:
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosing = true
		c.g.chooserDeclinable = declinable
		c.g.chooserPrompt = prompt
		c.g.chooserCandidates = candidates
		c.g.promptSource = source
		c.g.promptCursor, c.g.hasCursor = 0, false
		c.g.btnCursor, c.g.hasBtnCursor = 0, false
		c.g.openZoneForPrompt(candidates)
	})
	r := <-c.reply
	c.g.dispatch(func(app.Context) {
		c.g.choosing = false
		c.g.chooserDeclinable = false
		c.g.chooserPrompt = ""
		c.g.chooserCandidates = nil
		c.g.promptSource = ""
		c.g.promptCursor, c.g.hasCursor = 0, false
		c.g.btnCursor, c.g.hasBtnCursor = 0, false
		c.g.closeZoneForPrompt()
	})
	return r.id, r.ok
}

// openZoneForPrompt opens the out-of-play zone viewer when a prompt's candidates
// live outside play. The board only draws cards in play, so a prompt over a
// discard pile (World Tree, Witch of the Eye) would otherwise have nothing to
// click; the viewer makes the pile the board, scrolled to the right row.
func (g *game) openZoneForPrompt(candidates []engine.LocalID) {
	first := candidates[0]
	for p := range 2 {
		for _, z := range []struct {
			label string
			ids   []engine.LocalID
		}{
			{"Discard", g.g.Discard(p)},
			{"Archives", g.g.Archives(p)},
			{"Purge", g.g.Purge(p)},
			{"Deck", g.g.Deck(p)},
		} {
			if !containsID(z.ids, first) {
				continue
			}
			g.zonesPlayer, g.promptZone, g.promptZoneScrolled = p, z.label, false
			return
		}
	}
}

// closeZoneForPrompt closes a zone viewer that a prompt opened, leaving one the
// player opened themselves alone.
func (g *game) closeZoneForPrompt() {
	if g.promptZone == "" {
		return
	}
	g.zonesPlayer, g.promptZone = -1, ""
}

// ChooseOption implements the engine's OptionChooser: it posts a labeled
// multiple-choice prompt (e.g. whether to take archived cards into hand) to the
// UI and blocks until the player clicks one of the option buttons. Without this,
// the engine falls back to the first option — which silently auto-took archives.
func (c *webChooser) ChooseOption(source, prompt string, options []string) int {
	// Drop any stale reply so a leftover click cannot answer this prompt.
	select {
	case <-c.optionReply:
	default:
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosingOption = true
		c.g.optionPrompt = prompt
		c.g.optionLabels = options
		c.g.promptSource = source
	})
	i := <-c.optionReply
	c.g.dispatch(func(app.Context) {
		c.g.choosingOption = false
		c.g.optionPrompt = ""
		c.g.optionLabels = nil
		c.g.promptSource = ""
	})
	return i
}

// ---- the click handlers a prompt is answered with ----

func (g *game) chooseCandidate(_ app.Context, id engine.LocalID) {
	if !g.choosing {
		return
	}
	select {
	case g.chooser.reply <- chooseReply{id: id, ok: true}:
	default:
	}
}

// declineChooser answers a declinable card prompt with a pass — the Done button,
// and what Escape means while such a prompt is up.
func (g *game) declineChooser(_ app.Context, _ app.Event) {
	if !g.choosing || !g.chooserDeclinable {
		return
	}
	select {
	case g.chooser.reply <- chooseReply{}:
	default:
	}
}

// chooseOptionIdx answers the current option prompt with option i. The index is
// stable per button position, so a captured value is safe here (unlike per-card
// closures).
func (g *game) chooseOptionIdx(i int) app.EventHandler {
	return func(_ app.Context, _ app.Event) {
		if !g.choosingOption {
			return
		}
		select {
		case g.chooser.optionReply <- i:
		default:
		}
	}
}

// onScorePillClick opens the out-of-play zone viewer for the clicked player. The
// player index is read from the zone counts' data attribute rather than captured
// in a closure, so the single stable handler stays valid across re-renders (go-app
// compares event handlers by function pointer). The viewer is read-only, so it
// stays available during a prompt — e.g. to inspect archives before deciding
// whether to take them into hand.
func (g *game) onScorePillClick(ctx app.Context, _ app.Event) {
	p, err := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("player").String())
	if err != nil {
		return
	}
	g.zonesPlayer = p
}

// closeZones hides the out-of-play zone viewer.
func (g *game) closeZones(_ app.Context, _ app.Event) {
	// A viewer opened by a prompt is the only place its candidates are clickable,
	// so it stays up until the prompt is answered.
	if g.promptZone != "" {
		return
	}
	g.zonesPlayer = -1
}

// stopClick keeps a click inside the zone panel from bubbling up to the
// backdrop's close handler, so only clicks outside the panel dismiss the viewer.
func (g *game) stopClick(_ app.Context, e app.Event) {
	e.Call("stopPropagation")
}

// cancelChooser answers the current creature chooser with "no choice", so a
// player in manual mode can escape a prompt with no clickable candidate.
func (g *game) cancelChooser(_ app.Context, _ app.Event) {
	if !g.choosing {
		return
	}
	select {
	case g.chooser.reply <- chooseReply{ok: false}:
	default:
	}
}
