// Package sim plays whole random, legal Vactrol games to shake out engine bugs
// that per-card unit tests never reach — a card leaking Æmber, a move-between-
// zones path that duplicates or drops a card, a turn that corrupts the flat state.
// Every game is decoded from a byte script and checked against the engine's own
// invariants (engine.Game.InvariantError) after each step.
//
// One Simulate function backs all three entry points in sim_test.go: the
// coverage-guided FuzzPlay fuzz target, a fixed-seed property test wired into
// `mage test`, and the long-running `mage soak`. Because a game is a pure function
// of its script, any failure is a deterministic input that replays exactly and can
// be promoted to a unit test.
//
// This package lives outside internal/engine on purpose: it drives the engine
// through its public API and the real card database (via internal/match), and the
// engine's 100% coverage gate applies only to engine, not here.
package sim

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// Caps bound a single simulated game so every run terminates — the "put a limit on
// everything" rule. A very long real game stays far under these: 50 turns per
// player is a marathon, and 100 decisions in one turn is already unreachable.
const (
	maxTurns            = 100 // 50 turns per player
	maxDecisionsPerTurn = 100
)

// Simulate plays one random legal game driven entirely by script, checking the
// engine's flat-state invariants after setup and after every action. It returns a
// non-nil error — the first invariant violation or a recovered engine panic — with
// enough context to reproduce, or nil when the game played out cleanly. The whole
// game is a pure function of script, so a returned failure replays deterministically
// from the same bytes.
func Simulate(script []byte) error {
	_, err := simulate(script, false)
	return err
}

// simulate is Simulate with the played-out game handed back, so a debug replay can
// read the game log that led to a failure. With verbose set the game records that
// log, which a soak or fuzz run does not want to pay for.
func simulate(script []byte, verbose bool) (g *engine.Game, err error) {
	d := &decoder{script: script}
	seed := int64(d.uint64())
	g, houses := match.New("P0", "P1", seed)
	g.Verbose = verbose
	ch := &scriptChooser{d: d}
	g.SetChooser(0, ch)
	g.SetChooser(1, ch)
	g.SetPlayerHouses(0, houses[0])
	g.SetPlayerHouses(1, houses[1])

	// An engine panic (including an -tags assert invariant panic) becomes a normal
	// error so the caller can report the offending script instead of crashing the run.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during simulation: %v\n%s", r, debug.Stack())
		}
	}()

	if e := g.InvariantError(); e != nil {
		return g, fmt.Errorf("invariant violated at setup: %w", e)
	}

	// An exhausted script ends the game: with no decisions left no player acts, and
	// a player who never acts never gains Æmber, so the remaining turns would be
	// identical empty ones ground out until maxTurns.
	for turn := 0; turn < maxTurns && g.Winner() < 0 && !d.done(); turn++ {
		player := turn % 2
		g.StartTurn(player)
		if e := g.InvariantError(); e != nil {
			return g, fmt.Errorf(
				"invariant violated after StartTurn (turn %d, player %d): %w",
				turn,
				player,
				e,
			)
		}
		if g.Winner() >= 0 {
			break
		}
		if e := playTurn(g, player, houses[player], d); e != nil {
			return g, e
		}
		g.EndPlayPhase(player)
		if e := g.InvariantError(); e != nil {
			return g, fmt.Errorf(
				"invariant violated after EndPlayPhase (turn %d, player %d): %w",
				turn,
				player,
				e,
			)
		}
	}
	return g, nil
}

// playTurn chooses the active house, then applies up to maxDecisionsPerTurn script-
// driven actions, checking invariants after each. It stops early when the script
// runs out (which winds the game down toward termination), the player chooses to
// stop, no legal action remains, or the game ends.
func playTurn(g *engine.Game, player int, houses []engine.House, d *decoder) error {
	if len(houses) > 0 {
		_ = g.ChooseHouse(player, houses[int(d.byte())%len(houses)])
	}
	for step := 0; step < maxDecisionsPerTurn; step++ {
		if g.Winner() >= 0 || d.done() {
			return nil
		}
		acted, err := doAction(g, player, d)
		if err != nil {
			return err
		}
		if !acted {
			return nil
		}
		if e := g.InvariantError(); e != nil {
			return fmt.Errorf(
				"invariant violated after action (player %d, step %d): %w",
				player,
				step,
				e,
			)
		}
	}
	return nil
}

// doAction gathers the currently legal plays, discards, and creature uses, lets the
// script pick one (or choose to end the turn), and applies it. It returns false
// when the player stops or there is nothing legal to do.
//
// An error means the engine refused a move its own CanPlay/CanUse had just
// declared legal — a contradiction worth failing the simulation over. Errors from
// the reap/fight/use-action split are not that: CanUse promises some use is legal,
// not the one the script picked, so those fall back rather than fail.
func doAction(g *engine.Game, player int, d *decoder) (bool, error) {
	var playable, discardable []engine.LocalID
	for _, id := range g.Hand(player) {
		if g.CanPlay(player, id) == nil {
			playable = append(playable, id)
		}
		if g.CanDiscard(player, id) == nil {
			discardable = append(discardable, id)
		}
	}
	var usable []engine.LocalID
	for _, id := range g.Battleline(player) {
		if g.CanUse(player, id) == nil {
			usable = append(usable, id)
		}
	}

	total := len(playable) + len(usable)
	if total == 0 {
		// Out of plays and uses, the turn is over. Ending it here rather than leaving
		// it to the rare stop roll also stops a hand of unplayable active-house cards
		// from keeping a dead turn alive as a run of discards.
		return false, nil
	}
	choice := d.pick(total, stopOneIn)
	if choice < 0 {
		// The rare branch splits between the two ways a player walks away from moves
		// they still have. Discarding shares it rather than sitting in the uniform
		// pool because nearly every active-house card in hand is discardable, so an
		// even slot would spend most hands on the discard pile instead of the board.
		if len(discardable) > 0 && d.bool() {
			return true, discardCard(g, player, discardable[int(d.byte())%len(discardable)])
		}
		return false, nil
	}
	if choice < len(playable) {
		return true, playCard(g, player, playable[choice], d)
	}
	return true, useCreature(g, player, usable[choice-len(playable)], d)
}

// discardCard discards the given hand card, which CanDiscard has already vetted.
func discardCard(g *engine.Game, player int, id engine.LocalID) error {
	idx := indexOf(g.Hand(player), id)
	if idx < 0 {
		return fmt.Errorf("card %d vanished from P%d's hand before discarding it", id, player)
	}
	if err := g.DiscardFromHand(player, idx); err != nil {
		return fmt.Errorf("CanDiscard allowed %s but discarding it failed: %w", g.Name(id), err)
	}
	return nil
}

// playCard plays the given hand card, dispatching on its type. flankLeft and any
// target choices come from the script through the chooser.
func playCard(g *engine.Game, player int, id engine.LocalID, d *decoder) error {
	idx := indexOf(g.Hand(player), id)
	if idx < 0 {
		return fmt.Errorf("card %d vanished from P%d's hand before playing it", id, player)
	}
	var err error
	switch g.Def(id).Type {
	case engine.Creature:
		_, err = g.PlayCreature(player, idx, d.bool())
	case engine.Artifact:
		_, err = g.PlayArtifact(player, idx)
	case engine.Tactic:
		err = g.PlayAction(player, idx)
	case engine.Upgrade:
		_, err = g.PlayUpgrade(player, idx)
	}
	if err != nil {
		return fmt.Errorf("CanPlay allowed %s but playing it failed: %w", g.Name(id), err)
	}
	return nil
}

// useCreature reaps, fights, or uses an action ability with the creature, starting
// from the use the script chose and falling through to the others when that one is
// not legal for it — a creature with no action ability has none to use, none can
// fight an empty battleline, and Tireless Crocag cannot reap. CanUse vouches only
// that *some* use is legal, so all three being refused means it contradicted itself.
func useCreature(g *engine.Game, player int, id engine.LocalID, d *decoder) error {
	start, target := int(d.byte()%3), int(d.byte())
	var err error
	for i := range 3 {
		switch (start + i) % 3 {
		case 0:
			err = g.Reap(player, id)
		case 1:
			enemies := g.Battleline(1 - player)
			if len(enemies) == 0 {
				err = errNoEnemies
				break
			}
			// Try every defender from the script's pick onward: taunt makes some of
			// them illegal without making fighting itself illegal.
			for j := range enemies {
				if err = g.Fight(player, id, enemies[(target+j)%len(enemies)]); err == nil {
					break
				}
			}
		default:
			err = g.UseAction(player, id)
		}
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("CanUse allowed %s but no use of it was legal: %w", g.Name(id), err)
}

// errNoEnemies stands in for the engine's refusal when the script picks a fight
// with an empty enemy battleline, so useCreature has one fallback path.
var errNoEnemies = errors.New("no enemy creature to fight")

// indexOf returns the position of id in ids, or -1 if absent.
func indexOf(ids []engine.LocalID, id engine.LocalID) int {
	for i, x := range ids {
		if x == id {
			return i
		}
	}
	return -1
}
