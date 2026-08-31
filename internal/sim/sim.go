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
func Simulate(script []byte) (err error) {
	d := &decoder{script: script}
	seed := int64(d.uint64())
	g, houses := match.New("P0", "P1", seed)
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
		return fmt.Errorf("invariant violated at setup: %w", e)
	}

	for turn := 0; turn < maxTurns && g.Winner() < 0; turn++ {
		player := turn % 2
		g.BeginTurn(player)
		if e := g.InvariantError(); e != nil {
			return fmt.Errorf(
				"invariant violated after BeginTurn (turn %d, player %d): %w",
				turn,
				player,
				e,
			)
		}
		if g.Winner() >= 0 {
			break
		}
		if e := playTurn(g, player, houses[player], d); e != nil {
			return e
		}
		g.EndTurn(player)
		if e := g.InvariantError(); e != nil {
			return fmt.Errorf(
				"invariant violated after EndTurn (turn %d, player %d): %w",
				turn,
				player,
				e,
			)
		}
	}
	return nil
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
		if !doAction(g, player, d) {
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

// doAction gathers the currently legal plays and creature uses, lets the script
// pick one (or choose to stop), and applies it. It returns false when the player
// stops or there is nothing legal to do. Engine action errors are swallowed: the
// script is a heuristic move generator, not a legality oracle, and the invariant
// check after each action is what actually guards correctness.
func doAction(g *engine.Game, player int, d *decoder) bool {
	var playable []engine.LocalID
	for _, id := range g.Hand(player) {
		if g.CanPlay(player, id) == nil {
			playable = append(playable, id)
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
		return false
	}
	choice := int(d.byte()) % (total + 1) // the extra slot is "stop here"
	if choice == total {
		return false
	}
	if choice < len(playable) {
		playCard(g, player, playable[choice], d)
	} else {
		useCreature(g, player, usable[choice-len(playable)], d)
	}
	return true
}

// playCard plays the given hand card, dispatching on its type. flankLeft and any
// target choices come from the script through the chooser.
func playCard(g *engine.Game, player int, id engine.LocalID, d *decoder) {
	idx := indexOf(g.Hand(player), id)
	if idx < 0 {
		return
	}
	switch g.Def(id).Type {
	case engine.Creature:
		_, _ = g.PlayCreature(player, idx, d.bool())
	case engine.Artifact:
		_, _ = g.PlayArtifact(player, idx)
	case engine.Tactic:
		_ = g.PlayAction(player, idx)
	case engine.Upgrade:
		_, _ = g.PlayUpgrade(player, idx)
	}
}

// useCreature reaps, fights, or uses an action ability with the creature, chosen by
// the script. Fighting falls back to reaping when the enemy battleline is empty.
func useCreature(g *engine.Game, player int, id engine.LocalID, d *decoder) {
	switch d.byte() % 3 {
	case 0:
		_ = g.Reap(player, id)
	case 1:
		enemies := g.Battleline(1 - player)
		if len(enemies) == 0 {
			_ = g.Reap(player, id)
			return
		}
		_ = g.Fight(player, id, enemies[int(d.byte())%len(enemies)])
	default:
		if err := g.UseAction(player, id); err != nil {
			_ = g.Reap(player, id)
		}
	}
}

// indexOf returns the position of id in ids, or -1 if absent.
func indexOf(ids []engine.LocalID, id engine.LocalID) int {
	for i, x := range ids {
		if x == id {
			return i
		}
	}
	return -1
}
