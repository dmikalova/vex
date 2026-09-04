package engine

import (
	"math/rand"
	"strings"
)

// This file holds the Game object itself — the live match harness that bundles
// the flat GameState with the read-only catalog and the surrounding services
// (player names, choosers, RNG, log) — plus the chooser interfaces the engine
// calls when an effect must make a decision. The Game's behaviors are spread
// across the other game_*.go files (turn, play, combat, destruction, and so on).

// Win/economy constants.
const (
	// KeyCost is the amount of Æmber required to forge one key.
	KeyCost = 6
	// KeysToWin is the number of keys a player must forge to win.
	KeysToWin = 3
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

// Game bundles the flat GameState with the read-only Catalog and the surrounding
// engine services (player names, choosers, RNG, log). Cloning a state for MCTS
// only needs GameState.FastCopy; this wrapper is the live match harness.
type Game struct {
	State   GameState
	Verbose bool
	// Log is the public narration of the match — the same for both players, naming
	// no card in a hidden zone (ADR 0011). It lives on Game, not GameState, because
	// entries are interface values and so cannot satisfy the flat, pointerless,
	// comparable state ADR 0005 requires.
	Log []Record
	// frames is the stack of open attribution frames; see Game.openFrame.
	frames []Frame
	// triggerDepth is how many TriggerAbility resolutions are open. Two Replicators
	// would trigger each other's reap effect forever, so the Rule of Six bounds the
	// chain the same way it bounds a repeated effect.
	triggerDepth int
	// recording is whether outcomes are narrated at all. NewGame turns it on; a bot
	// exploring cloned positions turns it off so the log costs nothing.
	recording bool

	names    [2]string
	choosers [2]Chooser
	cat      *catalog
	rng      *rand.Rand
	// houses[p] is the set of houses in player p's deck — the houses they may choose
	// as their active house. Empty means unknown, in which case any house is allowed
	// (so tests and the AI need not declare deck houses). A frontend sets it so a
	// forced house (Control the Weak) that the player lacks is ignored: cannot
	// overrides must.
	houses [2][]House
	// manual turns on manual mode: house restrictions on playing and using cards
	// are lifted so a UI can rearrange the game freely. See game_manual.go.
	manual bool
	// settling is true while a destruction batch or a state-based sweep is running,
	// so the sweep does not re-enter and split a batch's simultaneous timing.
	settling bool
}

// NewGame creates a new two-player game seeded for deterministic play.
func NewGame(p0Name, p1Name string, seed int64) *Game {
	g := &Game{
		names:     [2]string{p0Name, p1Name},
		choosers:  [2]Chooser{FirstChooser{}, FirstChooser{}},
		cat:       &catalog{},
		rng:       rand.New(rand.NewSource(seed)),
		recording: true,
	}
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
	return g.chooserFor(player).ChooseCreature(source, renderPrompt(source, prompt), candidates)
}

// pickCard resolves a "choose one card" prompt. It uses the same chooser channel
// as creature choices because a prompt is still one visible card chosen from a set;
// callers are responsible for passing the legal card candidates.
func (g *Game) pickCard(player int, source, prompt string, candidates []LocalID) (LocalID, bool) {
	if len(candidates) == 1 {
		return candidates[0], true
	}
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
