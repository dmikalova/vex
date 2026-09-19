package engine

import (
	"strings"
)

// This file holds the Game object itself — the live match harness that bundles
// the flat GameState with the read-only catalog and the surrounding services
// (player names, choosers, log) — plus the chooser interfaces the engine
// calls when an effect must make a decision. The Game's behaviors are spread
// across the other game_*.go files (turn, play, combat, destruction, and so on).

// Win/economy constants.
const (
	// KeyCost is the amount of Æmber required to forge one key.
	KeyCost = 6
	// KeysToWin is the number of keys a player must forge to win.
	KeysToWin = 3
	// MaxKeys is the most keys a player can hold. A player wins on their third key,
	// but an effect can forge a fourth (colourless) key before the game ends, so the
	// key count and the KeyColors slots run one past KeysToWin.
	MaxKeys = 4
	// HandSize is the number of cards a player draws back up to at end of turn.
	HandSize = 6
)

// Chooser makes target decisions for a player. The engine calls it whenever an
// effect must pick a creature. Implementations must be deterministic so games can
// be reproduced from a seed.
type Chooser interface {
	// ChooseCreature returns one id from candidates and true, or false if none.
	// source is the name of the card whose ability is asking (for prompt
	// attribution), or "" when the choice has no card source such as an ordering
	// or turn-structure prompt.
	ChooseCreature(source, prompt string, candidates []LocalID) (LocalID, bool)
}

// FirstChooser always picks the first available candidate. It is the default and
// keeps behavior deterministic for tests and simulation.
type FirstChooser struct{}

// FirstChooser deliberately implements only the base Chooser: every optional
// capability falls back to its engine default, which is what keeps simulation and
// the bot from answering prompts a human would render. A prompt-totality test
// (ADR 0045) must scope itself to the human-facing chooser, not to this one.
var _ Chooser = FirstChooser{}

// ChooseCreature returns the first candidate, or false if the list is empty.
func (FirstChooser) ChooseCreature(_, _ string, candidates []LocalID) (LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	return candidates[0], true
}

// OptionChooser is an optional Chooser capability: choosing one of several
// labeled options, for "choose one" effects. Choosers that do not implement it
// default to the first option.
type OptionChooser interface {
	ChooseOption(source, prompt string, options []string) int
}

// PositionChooser is an optional Chooser capability: choosing where a Deploy
// creature enters its controller's battleline by pointing at the line itself —
// position i lands the creature before line[i] (0 the left flank, len(line) the
// right flank). It lets a client offer click-to-place on the battleline instead
// of one labeled option per gap. A Chooser that does not implement it falls back
// to the OptionChooser channel with those gaps rendered as labeled options.
type PositionChooser interface {
	ChoosePosition(source, prompt string, line []LocalID) int
}

// DeclinableChooser is an optional Chooser capability: an optional card choice
// the player may pass on — the "you may purge a card" and "up to 3 creatures" of
// a card's text. It is separate from ChooseCreature because declining is itself
// an answer, so a sole candidate must still be offered rather than forced.
// A Chooser that does not implement it is asked through the OptionChooser channel
// with the candidate names plus a trailing DoneOption.
type DeclinableChooser interface {
	ChooseCardOrDecline(source, prompt string, candidates []LocalID) (LocalID, bool)
}

// DoneOption labels the pass entry appended to a declinable prompt's fallback
// option list.
const DoneOption = "Done"

// Orderer is an optional Chooser capability: arranging ids into a resolution
// order in a single call, instead of being asked to pick the next id repeatedly.
// A Chooser that implements it takes full control of ordering (see
// Game.orderByChoice); one that does not falls back to repeated ChooseCreature.
type Orderer interface {
	OrderCreatures(source, prompt string, ids []LocalID) []LocalID
}

// OrderableReaction is one entry the active player may pick as the next to
// resolve in a trigger window: either a card's triggered ability (HasCard, Card
// its source) or a duration reaction from the lasting registry (Full Moon,
// Charge!, Crystal Hive), which has no card in play. Label is the rendered display
// string for both — the source card's name and ability text, or the duration
// reaction's rendered effect. Card abilities and duration reactions are offered in
// one flat labeled list so the whole window orders as one.
type OrderableReaction struct {
	Card    LocalID
	HasCard bool
	Label   string
}

// ReactionChooser is an optional Chooser capability: picking which reaction in a
// trigger window resolves next, from a flat labeled list of the window's pending
// reactions — card abilities and duration reactions together. The engine asks
// repeatedly, dropping the picked reaction each time, until one remains (which is
// forced). Returning an out-of-range index leaves the remaining reactions in their
// gathered order. A Chooser that does not implement it resolves the whole window in
// its gathered order. The engine never asks when every pending reaction is
// identical, since their order cannot matter.
type ReactionChooser interface {
	ChooseReaction(prompt string, reactions []OrderableReaction) int
}

// BadgeChooser is an optional Chooser capability: a client that previews the
// status a creature it is about to choose will receive — the "3 damage" a
// Festering Touch pick deals, the ward an Imperium "ward N" places. The effect
// sends the badge before its choose loop and the zero badge after, so the client
// can badge each candidate as it is picked and animate the badges away when the
// loop ends. It is display-only; a Chooser that does not implement it ignores it.
type BadgeChooser interface {
	PreviewBadge(badge SelectionBadge)
}

// ActionChooser is an optional Chooser capability: choosing the active player's
// next ROOT action from the legal set (ADR 0039) — the suspension point a canonical
// turn loop asks through, where the other capabilities answer one choice within an
// action. It has no engine default: a Chooser that does not implement it
// (FirstChooser, the bot, the sim) drives root actions its own way rather than
// through the loop, so only an interactive driver installs one.
type ActionChooser interface {
	ChooseAction(actions []Command) Command
}

// Game bundles the flat GameState with the read-only Catalog and the surrounding
// engine services (player names, choosers, RNG, log). Cloning a state for MCTS
// only needs GameState.FastCopy; this wrapper is the live match harness.
type Game struct {
	// State is the flat mutable game state; FastCopy clones it for MCTS rollouts.
	State GameState
	// Verbose, when set, prints each log entry to stdout as it is narrated.
	Verbose bool
	// Log is the public narration of the match — the same for both players, naming
	// no card in a hidden zone (ADR 0011). It lives on Game, not GameState, because
	// entries are interface values and so cannot satisfy the flat, pointerless,
	// comparable state ADR 0005 requires.
	Log []Record
	// frames is the stack of open attribution frames; see Game.openFrame.
	frames []Frame
	// recording is whether outcomes are narrated at all. NewGame turns it on; a bot
	// exploring cloned positions turns it off so the log costs nothing.
	recording bool

	// Engine services around the state: player names, per-player choosers, and the
	// read-only card catalog. The match RNG is not here — it lives flat in
	// GameState.PRNG so a snapshot captures it and replay is bit-exact (ADR 0039).
	names    [2]string
	choosers [2]Chooser
	cat      *catalog
	// houses[p] is the set of houses in player p's deck — the houses they may choose
	// as their active house. Empty means unknown, in which case any house is allowed
	// (so tests and the AI need not declare deck houses). A frontend sets it so a
	// forced house (Control the Weak) that the player lacks is ignored: cannot
	// overrides must.
	houses [2][]House
	// nameableNames is every card name a player may name (Etan's Jar), injected by
	// the match because the engine cannot read the card database itself (ADR 0003).
	// Empty falls back to the names present in this match.
	nameableNames []string
	// manual turns on manual mode: house restrictions on playing and using cards
	// are lifted so a UI can rearrange the game freely. See game_manual.go.
	manual bool
	// settling is true while a destruction batch or a state-based sweep is running,
	// so the sweep does not re-enter and split a batch's simultaneous timing.
	settling bool
	// deferringLeaves is true while a simultaneous batch is running, so a card's
	// "Leaves Play:" window is gathered as it goes but held until every card in the
	// batch has moved. deferredLeaves is where those windows wait. Both are runtime
	// scheduling, not game state: they never outlive the batch, so they stay off
	// GameState and out of the undo snapshot (ADR 0005).
	deferringLeaves bool
	deferredLeaves  []triggeredAbility
	// destroyingSource names the card whose effect is carrying out the current
	// destruction, so the batch it targets narrates as one grouped line ("Strange
	// Gizmo destroys A, B, and C") instead of a passive line per creature. The next
	// destruction batch consumes and clears it, so state-based deaths that follow
	// (a creature that lost a buff) and combat damage — which have no such agent —
	// stay passive.
	destroyingSource    LocalID
	hasDestroyingSource bool
	// shuffleBatch, while a batch is open, collects the cards an effect shuffles
	// into a deck so they narrate as one grouped line per owner ("Lost in the Woods
	// shuffles A and B into P2's deck") instead of a passive line per creature.
	shuffleBatch    []LocalID
	batchingShuffle bool
	// destroyingWindow holds every creature enrolled in the currently open
	// "Destroyed:" window — those being destroyed together, whether by the batch
	// that opened the window or by a Destroyed ability that destroyed more creatures
	// mid-window. A nested destruction that re-selects one of them (Harbinger of
	// Doom's "Destroyed: destroy each creature" re-selects Harbinger itself) skips
	// it, so its Destroyed abilities fire once and it is discarded once. The batch
	// that opened the window owns it and clears the field once every enrolled
	// creature has reached the discard pile.
	destroyingWindow []LocalID
	// destroyWindowOpen is true while a "Destroyed:" window is open, so a nested
	// destruction (a Destroyed ability destroying more creatures) enrolls its
	// creatures in the same window rather than opening its own — a whole chain of
	// deaths shares one simultaneous window (KeyForge timing; ADR 0013).
	destroyWindowOpen bool
	// destroyWindowController is the player who orders the open window's Destroyed
	// abilities — the one who caused the destruction that opened it.
	destroyWindowController int
	// destroyPending is the open window's queue of Destroyed abilities still to
	// resolve. A Destroyed ability that destroys more creatures appends theirs here,
	// so the resolve loop re-gathers and keeps going until the queue drains.
	destroyPending []triggeredAbility
	// powerComputing is the stack of creatures whose Power is mid-computation, so a
	// variable "X" power that reads a neighbor's power cannot recurse forever when
	// two such creatures reference each other (two Picaroons that both lost
	// Changeling to Grey Aberrant). A creature asked for its power while already on
	// the stack contributes 0 — an undeterminable value is 0 (KeyForge). It is
	// runtime scratch, always balanced by Power's defer, so it is not part of the
	// snapshotted state.
	powerComputing []LocalID
}

// NewGame creates a new two-player game seeded for deterministic play.
func NewGame(p0Name, p1Name string, seed int64) *Game {
	g := &Game{
		names:     [2]string{p0Name, p1Name},
		choosers:  [2]Chooser{FirstChooser{}, FirstChooser{}},
		cat:       &catalog{},
		recording: true,
	}
	g.State.PRNG = PRNG{State: uint64(seed)}
	g.State.Winner = -1
	return g
}

// SetChooser installs a custom chooser for a player (nil resets to the default).
func (g *Game) SetChooser(player int, c Chooser) { g.choosers[player] = c }

// SetPlayerHouses records the houses in a player's deck — the houses they may
// choose from. A frontend sets it so a forced active house the player does not
// have is ignored (cannot overrides must). When unset, any house is allowed.
func (g *Game) SetPlayerHouses(player int, houses []House) {
	g.houses[player] = append([]House(nil), houses...)
}

// PlayerName returns a player's display name.
func (g *Game) PlayerName(player int) string { return g.names[player] }

// chooserFor returns the chooser for a player, defaulting to FirstChooser.
func (g *Game) chooserFor(player int) Chooser {
	if ch := g.choosers[player]; ch != nil {
		return ch
	}
	return FirstChooser{}
}

// renderPrompt resolves the SelfName placeholder in a prompt to the name of the
// card asking, so a runtime prompt reads like the card's printed text ("fully
// heal Chuff Ape", not "fully heal {self}"). An unattributed prompt is left as is.
func renderPrompt(source, prompt string) string {
	if source == "" {
		return prompt
	}
	return strings.ReplaceAll(prompt, SelfName, source)
}

// pickCreature resolves a "choose one creature" prompt. When only one candidate
// is available the choice is forced, so it is taken automatically without
// consulting the chooser; otherwise the player's chooser decides (and may
// decline). Callers guard the empty case before calling.
func (g *Game) pickCreature(
	player int,
	source, prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	if len(candidates) == 1 {
		return candidates[0], true
	}
	// Boundary: settle before presenting the choice, so the player never chooses
	// among creatures one of which is already dead (ADR 0029).
	g.settleDestroyed(player)
	return g.chooserFor(player).ChooseCreature(source, renderPrompt(source, prompt), candidates)
}

// pickCard resolves a "choose one card" prompt. It uses the same chooser channel
// as creature choices because a prompt is still one visible card chosen from a set;
// callers are responsible for passing the legal card candidates.
func (g *Game) pickCard(player int, source, prompt string, candidates []LocalID) (LocalID, bool) {
	if len(candidates) == 1 {
		return candidates[0], true
	}
	// Boundary: settle before presenting the choice (ADR 0029).
	g.settleDestroyed(player)
	return g.chooserFor(player).ChooseCreature(source, renderPrompt(source, prompt), candidates)
}

// pickOptional resolves a "choose a card, or stop" prompt: the player may take one
// of candidates or decline. Unlike pickCreature it never short-circuits a sole
// candidate — passing is a legal answer, so forcing the last card would take the
// choice away. A chooser that cannot express a decline (no DeclinableChooser) is
// asked through the option channel with the candidate names plus DoneOption, which
// is the shape every optional prompt used to have.
func (g *Game) pickOptional(
	player int,
	source, prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	// Boundary: settle before presenting the choice (ADR 0029).
	g.settleDestroyed(player)
	if dc, ok := g.chooserFor(player).(DeclinableChooser); ok {
		return dc.ChooseCardOrDecline(source, renderPrompt(source, prompt), candidates)
	}
	options := make([]string, len(candidates)+1)
	for i, id := range candidates {
		options[i] = g.Name(id)
	}
	options[len(candidates)] = DoneOption
	if i := g.chooseOption(player, source, prompt, options); i < len(candidates) {
		return candidates[i], true
	}
	return 0, false
}

// orderByChoice asks controller to arrange ids into a resolution order by picking
// the next one repeatedly (the final id is forced, so it is never prompted). With
// 0 or 1 ids there is nothing to order and ids is returned unchanged; a rejected
// pick falls back to the remaining order. The default FirstChooser keeps the
// original order, so ordering only becomes interactive under a real UI.
func (g *Game) orderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	if len(ids) <= 1 {
		return ids
	}
	// Boundary: settle before presenting the choice, so the player never orders
	// among creatures one of which is already dead (ADR 0029).
	g.settleDestroyed(controller)
	if o, ok := g.chooserFor(controller).(Orderer); ok {
		return o.OrderCreatures("", prompt, ids)
	}
	remaining := append([]LocalID(nil), ids...)
	ordered := make([]LocalID, 0, len(ids))
	for len(remaining) > 1 {
		chosen, ok := g.chooserFor(controller).ChooseCreature("", prompt, remaining)
		if !ok {
			break
		}
		ordered = append(ordered, chosen)
		for i, id := range remaining {
			if id == chosen {
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	return append(ordered, remaining...)
}
