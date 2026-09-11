package web

import (
	"math/rand"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file bridges the engine's synchronous Chooser to the browser: a prompt
// raised while an effect resolves parks that goroutine, the UI renders the
// question, and the player's click sends the answer back.

// chooseReply carries the player's answer to an engine chooser request: the
// chosen id, or ok=false when the player cancels. auto is set when the player
// answered an ordering prompt with Auto-resolve, asking for a random order.
type chooseReply struct {
	id   engine.LocalID
	ok   bool
	auto bool
}

// webChooser adapts the engine's synchronous Chooser to go-app's single UI
// goroutine. The engine calls ChooseCreature from a background action goroutine
// (see game.runAction); it shows the chooser overlay on the UI goroutine, then
// blocks until a candidate is clicked (or the request is cancelled).
type webChooser struct {
	g           *game
	reply       chan chooseReply
	optionReply chan int
	// positionReply carries the battleline position a Deploy placement resolves to.
	positionReply chan int
	// cancel is closed by a manual-mode Cancel to drain the rest of the action: once
	// closed every prompt this effect raises answers itself immediately, so a single
	// click backs the whole action out. runAction remakes it before the next action.
	cancel chan struct{}
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
	r := c.raise(source, prompt, candidates, declinable, false)
	return r.id, r.ok
}

// OrderCreatures implements the engine's Orderer: instead of being asked to pick
// the next id repeatedly, the whole window is ordered here so an Auto-resolve
// button can offer a single random order. It otherwise reproduces the engine's
// default repeated-pick loop — each pick is a card prompt, the last id is forced —
// so ordering a batch by clicking cards is unchanged; Auto-resolve shuffles
// whatever remains and returns it, ending the window in one click.
func (c *webChooser) OrderCreatures(
	source, prompt string,
	ids []engine.LocalID,
) []engine.LocalID {
	remaining := append([]engine.LocalID(nil), ids...)
	ordered := make([]engine.LocalID, 0, len(ids))
	for len(remaining) > 1 {
		r := c.raise(source, prompt, remaining, false, true)
		if r.auto {
			rand.Shuffle(len(remaining), func(i, j int) {
				remaining[i], remaining[j] = remaining[j], remaining[i]
			})
			return append(ordered, remaining...)
		}
		if !r.ok {
			break
		}
		ordered = append(ordered, r.id)
		for i, id := range remaining {
			if id == r.id {
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	return append(ordered, remaining...)
}

// raise shows a card prompt on the UI goroutine and blocks the action goroutine
// until a candidate is clicked, the prompt is declined or Auto-resolved, or it is
// cancelled. ordering marks an Orderer prompt so the controls offer Auto-resolve.
func (c *webChooser) raise(
	source, prompt string,
	candidates []engine.LocalID,
	declinable, ordering bool,
) chooseReply {
	if len(candidates) == 0 {
		return chooseReply{}
	}
	// A manual-mode Cancel already in flight drains without showing the prompt: this
	// and every later prompt of the same action answer themselves so one click backs
	// the action out.
	select {
	case <-c.cancel:
		return chooseReply{}
	default:
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
		c.g.chooserOrdering = ordering
		c.g.chooserPrompt = prompt
		c.g.chooserCandidates = candidates
		c.g.promptSource = source
		c.g.promptCursor, c.g.hasCursor = 0, false
		c.g.btnCursor, c.g.hasBtnCursor = 0, false
		c.g.presentPrompt(candidates, declinable)
	})
	var r chooseReply
	select {
	case r = <-c.reply:
	case <-c.cancel:
		r = chooseReply{ok: false}
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosing = false
		c.g.chooserDeclinable = false
		c.g.chooserOrdering = false
		c.g.chooserPrompt = ""
		c.g.chooserCandidates = nil
		c.g.promptAsButtons = false
		c.g.promptSource = ""
		c.g.promptCursor, c.g.hasCursor = 0, false
		c.g.btnCursor, c.g.hasBtnCursor = 0, false
		c.g.closeZoneForPrompt()
	})
	return r
}

// maxPromptButtons bounds how many out-of-play candidates a prompt lists as
// action-bar buttons before it falls back to the zone viewer. A "look at the top
// N cards" pick (Navigator Ali, Lay of the Land — N is 3) is a short enough list
// to read as buttons; a choice over a whole discard pile is not.
const maxPromptButtons = 6

// presentPrompt decides how a card prompt's candidates are offered. Candidates in
// play highlight where they stand, so nothing is set up here for them. A bounded,
// mandatory pick from the top of the deck (a "look at the top N cards and put them
// back" — Navigator Ali, Eyegore, Lay of the Land) is offered as a short list of
// action-bar buttons, since the cards are hidden until the effect reveals them and
// the set is small. Every other out-of-play pick — a visible pile (discard,
// archives, purge), or an unbounded one (declinable — shuffle any number) — opens
// the zone viewer, which shows the whole pile as full card faces.
func (g *game) presentPrompt(candidates []engine.LocalID, declinable bool) {
	p, label, inPile := g.zoneOfCard(candidates[0])
	if !inPile {
		return
	}
	if label == "Deck" && !declinable && len(candidates) <= maxPromptButtons {
		g.promptAsButtons = true
		return
	}
	g.zonesPlayer, g.promptZone, g.promptZoneScrolled = p, label, false
}

// zoneOfCard finds the out-of-play pile a card sits in, if any, so a prompt over
// that pile knows which viewer row to open or which player's cards to name.
func (g *game) zoneOfCard(id engine.LocalID) (player int, label string, ok bool) {
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
			if containsID(z.ids, id) {
				return p, z.label, true
			}
		}
	}
	return 0, "", false
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
	// A manual-mode Cancel already in flight drains without showing the prompt.
	select {
	case <-c.cancel:
		return 0
	default:
	}
	// The Play creature / Play upgrade buttons already made the creature-as-upgrade
	// choice, so answer the engine's play-as-which prompt from the armed choice
	// rather than raising it a second time in the sidebar.
	if i, ok := c.armedUpgradeChoice(options); ok {
		return i
	}
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
	var i int
	select {
	case i = <-c.optionReply:
	case <-c.cancel:
		i = 0
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosingOption = false
		c.g.optionPrompt = ""
		c.g.optionLabels = nil
		c.g.promptSource = ""
	})
	return i
}

// armedUpgradeChoice answers the engine's "play as a creature or an upgrade?"
// prompt from the choice the Play creature / Play upgrade buttons armed. It fires
// only for that exact prompt — the ["Creature", "Upgrade"] option pair — and only
// while a choice is armed, so every other option prompt still asks the player.
// The armed choice is set on the UI goroutine before the play's goroutine starts,
// so reading it here (off that goroutine) sees the value set before the play ran.
func (c *webChooser) armedUpgradeChoice(options []string) (int, bool) {
	if c.g.upgradeChoice == choiceNone {
		return 0, false
	}
	if len(options) != 2 || options[0] != "Creature" || options[1] != "Upgrade" {
		return 0, false
	}
	if c.g.upgradeChoice == choiceUpgrade {
		return 1, true
	}
	return 0, true
}

// ChoosePosition implements the engine's PositionChooser: instead of a labeled
// option per battleline gap, it lights the line's creatures and lifts the
// creature being placed with its placement verbs on it (deployActions) — click a
// creature to land beside it, to its left or right per the armed direction
// toggle, or take a flank. It returns the position the new creature enters before
// (0 the left flank, len(line) the right flank). The engine renders the prompt
// text; the client speaks the choice through the lifted card and the lit line, so
// the prompt string is unused here.
func (c *webChooser) ChoosePosition(source, _ string, line []engine.LocalID) int {
	// With no other friendly creatures in play there is only one placement (the
	// lone spot), so place it without asking. The engine already skips the prompt
	// for an empty line; this guards the web side against ever raising a choice
	// that has no alternatives.
	if len(line) == 0 {
		return 0
	}
	// A manual-mode Cancel already in flight drains without showing the prompt.
	select {
	case <-c.cancel:
		return 0
	default:
	}
	// Drop any stale reply so a leftover click cannot answer this prompt.
	select {
	case <-c.positionReply:
	default:
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosingPosition = true
		c.g.positionLine = line
		c.g.positionRight = false
		c.g.positionSideChosen = false
		c.g.promptSource = source
	})
	var pos int
	select {
	case pos = <-c.positionReply:
	case <-c.cancel:
		pos = 0
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosingPosition = false
		c.g.positionLine = nil
		c.g.positionRight = false
		c.g.positionSideChosen = false
		c.g.promptSource = ""
	})
	return pos
}

// ---- the click handlers a prompt is answered with ----

func (g *game) chooseCandidate(_ app.Context, id engine.LocalID) {
	if !g.choosing {
		return
	}
	g.inspecting = false
	// Remember the creature just chosen so a "choose how to use X" verb prompt that
	// follows (Universal Translator uses a creature, then asks how) can lift it and
	// put its use buttons on it rather than in the sidebar.
	g.useTarget, g.hasUseTarget = id, true
	select {
	case g.chooser.reply <- chooseReply{id: id, ok: true}:
	default:
	}
}

// onPromptButtonPick answers a bounded card prompt with the candidate its
// action-bar button stands for. The id is read off the button's own dataset, so
// the single stable handler stays valid across re-renders (go-app compares
// handlers by pointer) — the same pattern the log's card mentions use.
func (g *game) onPromptButtonPick(ctx app.Context, _ app.Event) {
	id, err := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("id").String())
	if err != nil {
		return
	}
	g.chooseCandidate(ctx, engine.LocalID(id))
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

// autoResolveOrder answers an ordering prompt with a request for a random order —
// the Auto-resolve button — so the player need not arrange abilities whose order
// they do not care about.
func (g *game) autoResolveOrder(_ app.Context, _ app.Event) {
	if !g.choosing || !g.chooserOrdering {
		return
	}
	select {
	case g.chooser.reply <- chooseReply{auto: true}:
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

// choosePositionCandidate answers a Deploy placement prompt with the position the
// clicked battleline creature implies: to its left (its own index) or its right
// (index + 1), per the side the player chose first. It does nothing until a side
// has been chosen, so the line is only answerable once the player has said which
// way the new creature lands.
func (g *game) choosePositionCandidate(_ app.Context, id engine.LocalID) {
	if !g.choosingPosition || !g.positionSideChosen {
		return
	}
	pos := -1
	for i, c := range g.positionLine {
		if c == id {
			pos = i
			break
		}
	}
	if pos < 0 {
		return
	}
	if g.positionRight {
		pos++
	}
	g.answerPosition(pos)
}

// chooseDeploySide arms which side of a clicked creature the Deploy creature
// lands on — the first step of placement — and marks the side chosen so the
// battleline becomes clickable for the second step. Left of the leftmost creature
// is the left flank and right of the rightmost is the right flank, so both flanks
// stay reachable without their own buttons.
func (g *game) chooseDeploySide(right bool) app.EventHandler {
	return func(_ app.Context, _ app.Event) {
		if !g.choosingPosition {
			return
		}
		g.positionRight = right
		g.positionSideChosen = true
	}
}

// deploySideBack undoes the side choice, returning the placement to its first
// step so the player can pick the other side.
func (g *game) deploySideBack(_ app.Context, _ app.Event) {
	if !g.choosingPosition {
		return
	}
	g.positionSideChosen = false
}

// answerPosition sends a resolved battleline position back to the parked action
// goroutine.
func (g *game) answerPosition(pos int) {
	g.inspecting = false
	select {
	case g.chooser.positionReply <- pos:
	default:
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
func (g *game) closeZones(ctx app.Context, e app.Event) {
	// A viewer opened by a prompt is the only place its candidates are clickable,
	// so it normally stays up until the prompt is answered. A declinable prompt
	// (Not Finished with You — shuffle any number, including zero) can be finished
	// from here: closing submits the current selection by declining the next pick.
	if g.promptZone != "" {
		if g.chooserDeclinable {
			g.declineChooser(ctx, e)
		}
		return
	}
	g.zonesPlayer = -1
}

// stopClick keeps a click inside the zone panel from bubbling up to the
// backdrop's close handler, so only clicks outside the panel dismiss the viewer.
func (g *game) stopClick(_ app.Context, e app.Event) {
	e.Call("stopPropagation")
}

// cancelChooser backs the whole action out of a prompt in manual mode: it closes
// the chooser's cancel channel so this prompt and any that follow it answer
// themselves, draining the effect to completion. runAction then rolls the action
// back to the snapshot beginAction recorded. It is the only way out of a prompt
// with no clickable candidate (a mandatory card prompt or an option prompt).
func (g *game) cancelChooser(_ app.Context, _ app.Event) {
	if !g.choosing && !g.choosingOption && !g.choosingPosition {
		return
	}
	if g.cancelling {
		return
	}
	g.cancelling = true
	close(g.chooser.cancel)
}
