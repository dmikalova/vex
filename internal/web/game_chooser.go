package web

import (
	"math/rand"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is the prompt seam. A decision the engine needs is a Request the
// session is SUSPENDED on (ADR 0040), so "is a prompt up, and what does it ask?"
// is a read of session.Pending rather than state a background chooser pushed at
// the client; and answering one is applying the Command that request accepts.
// The readers come first, then the handlers a click answers through.

// pending returns the decision the session is waiting on, and whether it is
// waiting on one at all — false before the first deal and once the match has
// finished.
func (g *game) pending() (engine.Request, bool) {
	if g.s == nil {
		return engine.Request{}, false
	}
	req, done := g.s.Pending()
	return req, !done
}

// prompt returns the pending request when it is a PROMPT — a decision raised
// while a root action resolves — and false when the session is at rest between
// actions (a RequestAction), finished, or not yet dealt. It is the one place "a
// prompt is up" is decided.
func (g *game) prompt() (engine.Request, bool) {
	req, ok := g.pending()
	if !ok || req.Kind == engine.RequestAction {
		return engine.Request{}, false
	}
	return req, true
}

// atPrompt reports whether a prompt is waiting for an answer. It is what the old
// pair of guards — "an action goroutine is in flight" and "a chooser is up" —
// collapsed into: with the turn loop suspended inside the session there is no
// in-flight action any more, only a pending request that is not a RequestAction.
func (g *game) atPrompt() bool {
	_, ok := g.prompt()
	return ok
}

// choosing reports whether the prompt on screen is answered by clicking a card:
// the candidates highlight on the board and a click says which one. A trigger
// window is one too, whenever the board can express it (reactionByCard).
func (g *game) choosing() bool {
	req, ok := g.prompt()
	if !ok {
		return false
	}
	switch req.Kind {
	case engine.RequestPickCard, engine.RequestPickCardOrDecline:
		return true
	case engine.RequestReaction:
		return g.reactionByCard(req)
	}
	return false
}

// chooserCandidates is the cards a card prompt may be answered with, or nil when
// the prompt on screen is not answered by clicking a card.
func (g *game) chooserCandidates() []engine.LocalID {
	req, ok := g.prompt()
	if !ok {
		return nil
	}
	switch req.Kind {
	case engine.RequestPickCard, engine.RequestPickCardOrDecline:
		return req.Cards
	case engine.RequestReaction:
		if g.reactionByCard(req) {
			return reactionCards(req.Reactions)
		}
	}
	return nil
}

// chooserPrompt is the question the card prompt on screen asks.
func (g *game) chooserPrompt() string {
	req, _ := g.prompt()
	return req.Prompt
}

// promptSource identifies the card driving the prompt on screen, or the zero
// value for a prompt no card is attributable to (an ordering step, a trigger
// window). It is the engine's prompt identity rather than a name, so the source
// is the card itself and not a string that happens to match one.
func (g *game) promptSource() engine.PromptSource {
	req, _ := g.prompt()
	return req.Source
}

// chooserDeclinable marks a prompt the player may pass on — a "you may" or an "up
// to N". It adds the Done button and lets Escape answer the prompt instead of
// being swallowed.
func (g *game) chooserDeclinable() bool {
	req, ok := g.prompt()
	return ok && req.Kind == engine.RequestPickCardOrDecline
}

// chooserOrdering reports whether the card prompt on screen is one step of an
// ORDERING window — arranging several cards into a resolution order — which is
// what puts the Auto-resolve button on it.
//
// The engine has no single "here is the order" request: orderByChoice asks for the
// next card repeatedly, with no source card to attribute the prompt to. That makes
// an unattributed mandatory card pick exactly an ordering step, since every other
// mandatory pick is raised through an ability and carries its source card's name.
// With one candidate the order is already settled and there is nothing to resolve.
func (g *game) chooserOrdering() bool {
	req, ok := g.prompt()
	return ok && req.Kind == engine.RequestPickCard &&
		!req.Source.HasCard && len(req.Cards) >= 2
}

// choosingOption reports whether the prompt on screen is a labeled multiple choice
// answered with buttons: an engine option prompt, a trigger window the board
// cannot express, or the which-ability follow-up to a clicked reaction card.
func (g *game) choosingOption() bool {
	req, ok := g.prompt()
	if !ok {
		return false
	}
	switch req.Kind {
	case engine.RequestOption:
		return true
	case engine.RequestReaction:
		return !g.reactionByCard(req)
	}
	return false
}

// optionPrompt is the question the labeled prompt asks.
func (g *game) optionPrompt() string {
	req, ok := g.prompt()
	if !ok {
		return ""
	}
	if req.Kind == engine.RequestReaction && g.hasAbilityPick {
		return whichAbilityPrompt
	}
	return req.Prompt
}

// optionLabels is the buttons a labeled prompt offers, in the order they are
// shown. An option request offers its own Options; a trigger window offers the
// rendered ability text of whichever of its entries the buttons are choosing
// between (reactionChoices).
func (g *game) optionLabels() []string {
	req, ok := g.prompt()
	if !ok {
		return nil
	}
	switch req.Kind {
	case engine.RequestOption:
		return req.Options
	case engine.RequestReaction:
		choices := g.reactionChoices(req)
		labels := make([]string, len(choices))
		for i, at := range choices {
			labels[i] = req.Reactions[at].Label
		}
		return labels
	}
	return nil
}

// choosingPosition reports whether a battleline placement is being picked: the
// engine's Deploy prompt, or manual mode's put-into-play, which borrows the same
// picker over a line of its own.
func (g *game) choosingPosition() bool {
	if g.manualPlacing {
		return true
	}
	req, ok := g.prompt()
	return ok && req.Kind == engine.RequestPosition
}

// positionLine is the battleline a placement is being picked in.
func (g *game) positionLine() []engine.LocalID {
	if g.manualPlacing {
		return g.manualLine
	}
	req, ok := g.prompt()
	if !ok || req.Kind != engine.RequestPosition {
		return nil
	}
	return req.Cards
}

// ---- trigger windows ----

// whichAbilityPrompt is the follow-up when the clicked card carries more than one
// pending ability: the click said which card, this says which of its abilities.
const whichAbilityPrompt = "Choose which of its abilities resolves next"

// reactionByCard reports whether a trigger window is answered on the board, by
// clicking the card whose ability resolves next. It is, so long as two or more
// distinct source cards can be pointed at and the player is not already being
// asked which of one card's several abilities goes first.
func (g *game) reactionByCard(req engine.Request) bool {
	return !g.hasAbilityPick && len(reactionCards(req.Reactions)) >= 2
}

// reactionCards lists the distinct source cards a window can be answered by
// clicking. It reports none when any entry belongs to no card at all — a duration
// reaction — so such a window stays on the labeled list and no reaction is hidden.
// A source sitting in a pile (WithTriggersFromDiscard) is kept: presentPrompt opens
// that pile's viewer, which is closable while board candidates remain.
func reactionCards(reactions []engine.OrderableReaction) []engine.LocalID {
	var cards []engine.LocalID
	for _, r := range reactions {
		if !r.HasCard {
			return nil
		}
		if !containsID(cards, r.Card) {
			cards = append(cards, r.Card)
		}
	}
	return cards
}

// reactionChoices lists the entries of a trigger window the labeled buttons stand
// for, as indexes into the request's Reactions: every entry when the whole window
// is being answered by label, and just the clicked card's entries once the board
// has answered "which card" and the buttons are answering "which of its abilities".
func (g *game) reactionChoices(req engine.Request) []int {
	out := make([]int, 0, len(req.Reactions))
	for i, r := range req.Reactions {
		if g.hasAbilityPick && (!r.HasCard || r.Card != g.abilityPick) {
			continue
		}
		out = append(out, i)
	}
	return out
}

// pickReactionCard resolves a click on a trigger window's card. One pending
// ability on it answers the window outright; two or more ask which of them goes
// first, by buttons — never silently top-down.
func (g *game) pickReactionCard(ctx app.Context, req engine.Request, id engine.LocalID) {
	var on []int
	for i, r := range req.Reactions {
		if r.HasCard && r.Card == id {
			on = append(on, i)
		}
	}
	switch len(on) {
	case 0:
		return
	case 1:
		g.answer(ctx, engine.Command{
			Kind:  engine.CommandReaction,
			Index: on[0],
		})
	default:
		g.abilityPick, g.hasAbilityPick = id, true
		g.btnCursor, g.hasBtnCursor = 0, false
	}
}

// ---- how a card prompt's candidates are offered ----

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
// the zone viewer, which shows the whole pile as full card faces. A prompt whose
// candidates are split between a pile and the board (a trigger window holding a
// WithTriggersFromDiscard source alongside sources in play) opens the pile too; the
// viewer closes freely, so the board candidates are reachable either way.
func (g *game) presentPrompt(candidates []engine.LocalID, declinable bool) {
	p, label, inPile := g.firstPileCandidate(candidates)
	if !inPile {
		return
	}
	if label == "Deck" && !declinable && len(candidates) <= maxPromptButtons {
		g.promptAsButtons = true
		return
	}
	g.zonesPlayer, g.promptZone, g.promptZoneScrolled = p, label, false
}

// firstPileCandidate finds the out-of-play pile a prompt reaches into, reading the
// candidates in order so a window split across the board and a pile still opens
// the pile rather than leaving its candidates unreachable.
func (g *game) firstPileCandidate(
	candidates []engine.LocalID,
) (player int, label string, ok bool) {
	for _, id := range candidates {
		if p, l, inPile := g.zoneOfCard(id); inPile {
			return p, l, true
		}
	}
	return 0, "", false
}

// zoneOfCard finds the out-of-play pile a card sits in, if any, so a prompt over
// that pile knows which viewer row to open or which player's cards to name. The
// engine owns where a card is (ZoneOf); this only maps the piles that have a
// viewer row to their label, so a card in hand or in play opens nothing.
func (g *game) zoneOfCard(id engine.LocalID) (player int, label string, ok bool) {
	p, zone, in := g.eng().ZoneOf(id)
	if !in {
		return 0, "", false
	}
	switch zone {
	case engine.Discard:
		return p, zoneDiscardLabel, true
	case engine.Archives:
		return p, zoneArchivesLabel, true
	case engine.Purged:
		return p, zonePurgeLabel, true
	case engine.Deck:
		return p, zoneDeckLabel, true
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

// armedUpgradeChoice answers the engine's "play as a creature or an upgrade?"
// prompt from the choice the Play creature / Play upgrade buttons armed. It fires
// only for that exact prompt — the ["Creature", "Upgrade"] option pair — and only
// while a choice is armed, so every other option prompt still asks the player.
func (g *game) armedUpgradeChoice(options []string) (int, bool) {
	if g.upgradeChoice == choiceNone {
		return 0, false
	}
	if len(options) != 2 || options[0] != "Creature" || options[1] != "Upgrade" {
		return 0, false
	}
	if g.upgradeChoice == choiceUpgrade {
		return 1, true
	}
	return 0, true
}

// ---- the click handlers a prompt is answered with ----

// chooseCandidate answers the card prompt with the clicked card.
func (g *game) chooseCandidate(ctx app.Context, id engine.LocalID) {
	req, ok := g.prompt()
	if !ok || !g.choosing() || !containsID(g.chooserCandidates(), id) {
		return
	}
	g.inspecting = false
	// Remember the creature just chosen so a "choose how to use X" verb prompt that
	// follows (Universal Translator uses a creature, then asks how) can lift it and
	// put its use buttons on it rather than in the sidebar.
	g.useTarget, g.hasUseTarget = id, true
	g.recordBadge(id)
	if req.Kind == engine.RequestReaction {
		g.pickReactionCard(ctx, req, id)
		return
	}
	g.answer(ctx, engine.Command{
		Kind: engine.CommandPickCard,
		Card: id,
	})
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
func (g *game) declineChooser(ctx app.Context, _ app.Event) {
	if !g.chooserDeclinable() {
		return
	}
	g.answer(ctx, engine.Command{Kind: engine.CommandDecline})
}

// autoResolveOrder answers an ordering window without arranging it — the
// Auto-resolve button — so the player need not order cards whose order they do not
// care about. It applies a uniformly random legal pick for as long as the session
// keeps asking the same question: the same kind, Source and Prompt as the step the
// button was pressed on.
//
// It is a deliberate heuristic rather than one "here is the order" command. An
// ordering window reaches the client as a RUN of RequestPickCards (the engine's
// orderByChoice asks for the next card repeatedly), and a Command must stay
// comparable for Request.IsLegal, so none can carry a slice of ids. The cost of
// guessing wrong is bounded and harmless: an unrelated, identically-prompted pick
// arriving straight after the window has one extra pick auto-answered, and every
// answer is recorded as the ordinary CommandPickCard it is, so the game log and a
// replay stay faithful.
func (g *game) autoResolveOrder(ctx app.Context, _ app.Event) {
	req, ok := g.prompt()
	if !ok || !g.chooserOrdering() {
		return
	}
	for {
		next, live := g.prompt()
		if !live || next.Kind != req.Kind ||
			next.Source != req.Source || next.Prompt != req.Prompt ||
			len(next.Cards) == 0 {
			break
		}
		pick := engine.Command{
			Kind: engine.CommandPickCard,
			Card: next.Cards[rand.Intn(len(next.Cards))],
		}
		if err := g.s.Apply(pick); err != nil {
			break
		}
	}
	g.settleAfterApply(ctx)
}

// chooseOptionIdx answers the current labeled prompt with its i-th button. The
// index is stable per button position, so a captured value is safe here (unlike
// per-card closures).
func (g *game) chooseOptionIdx(i int) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		g.answerOption(ctx, i)
	}
}

// answerOption applies the command the i-th button of the labeled prompt stands
// for: an option index, or the reaction entry a trigger window's label names.
func (g *game) answerOption(ctx app.Context, i int) {
	req, ok := g.prompt()
	if !ok || !g.choosingOption() {
		return
	}
	if req.Kind == engine.RequestReaction {
		choices := g.reactionChoices(req)
		if i < 0 || i >= len(choices) {
			return
		}
		g.answer(ctx, engine.Command{
			Kind:  engine.CommandReaction,
			Index: choices[i],
		})
		return
	}
	if i < 0 || i >= len(req.Options) {
		return
	}
	g.answer(ctx, engine.Command{
		Kind:  engine.CommandOption,
		Index: i,
	})
}

// choosePositionCandidate answers a Deploy placement prompt with the position the
// clicked battleline creature implies: to its left (its own index) or its right
// (index + 1), per the side the player chose first. It does nothing until a side
// has been chosen, so the line is only answerable once the player has said which
// way the new creature lands.
func (g *game) choosePositionCandidate(ctx app.Context, id engine.LocalID) {
	if !g.choosingPosition() || !g.positionSideChosen {
		return
	}
	pos := indexOfID(g.positionLine(), id)
	if pos < 0 {
		return
	}
	if g.positionRight {
		pos++
	}
	if g.manualPlacing {
		g.manualPlaceInPlay(ctx, pos)
		return
	}
	g.answerPosition(ctx, pos)
}

// chooseDeploySide arms which side of a clicked creature the Deploy creature
// lands on — the first step of placement — and marks the side chosen so the
// battleline becomes clickable for the second step. Left of the leftmost creature
// is the left flank and right of the rightmost is the right flank, so both flanks
// stay reachable without their own buttons.
func (g *game) chooseDeploySide(right bool) app.EventHandler {
	return func(_ app.Context, _ app.Event) {
		if !g.choosingPosition() {
			return
		}
		g.positionRight = right
		g.positionSideChosen = true
	}
}

// deploySideBack undoes the side choice, returning the placement to its first
// step so the player can pick the other side.
func (g *game) deploySideBack(_ app.Context, _ app.Event) {
	if !g.choosingPosition() {
		return
	}
	g.positionSideChosen = false
}

// answerPosition answers a pending placement prompt with a battleline position.
func (g *game) answerPosition(ctx app.Context, pos int) {
	g.inspecting = false
	g.answer(ctx, engine.Command{
		Kind:  engine.CommandPosition,
		Index: pos,
	})
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

// closeZones hides the out-of-play zone viewer. It always closes, even under a
// mandatory prompt whose only candidates are in the pile: the player may need to
// read the board underneath to decide. Closing answers nothing — the prompt stays
// up, the zone stays recorded, and reopening the viewer returns to the same row.
func (g *game) closeZones(_ app.Context, _ app.Event) {
	g.zonesPlayer = -1
}

// stopClick keeps a click inside the zone panel from bubbling up to the
// backdrop's close handler, so only clicks outside the panel dismiss the viewer.
func (g *game) stopClick(_ app.Context, e app.Event) {
	e.Call("stopPropagation")
}
