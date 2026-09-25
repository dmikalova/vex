package web

import (
	"reflect"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestReplayReproducesLiveGame is the load-bearing invariant of the
// event-sourced command log: replaying the recorded inputs from a fresh deal
// must reproduce the live game exactly. It drives ordinary play through the real
// client handlers — mulligans, a house choice, creature plays (with their flank
// and position choices), and an end of turn — then replays the log the client
// recorded and asserts the rebuilt game's state and narrated log match the live
// one bit for bit.
func TestReplayReproducesLiveGame(t *testing.T) {
	c := newClient(t)

	// Play several turns of ordinary play so the log exercises the full range of
	// root actions and their choosers: house choices, creature plays (with their
	// flank and position choices), reaps, and end-of-turn windows.
	for turn := 0; turn < 6 && c.g.g.Winner() < 0; turn++ {
		h := c.startTurn()
		// Reap with every creature that can, before playing new ones.
		for _, id := range c.board() {
			if c.g.g.CanUseTo(c.g.active(), id, engine.ReapUse) == nil {
				c.g.selectBoardID(c.ctx, id)
				c.do(c.g.reap)
			}
		}
		for _, id := range c.hand() {
			def := c.g.g.Def(id)
			if def.Type == engine.Creature && def.House == h {
				c.playFromHand(id)
			}
		}
		c.pass()
	}

	live := c.g.g
	replayed, err := replayGame(c.g.seed, c.g.setNames, c.g.inputs, c.g.defByName)
	if err != nil {
		t.Fatalf("replayGame: %v", err)
	}

	if replayed.State != live.State {
		t.Errorf("replayed state does not match the live game")
	}
	if got, want := replayed.LogText(), live.LogText(); !reflect.DeepEqual(got, want) {
		t.Errorf("replayed log does not match the live game:\n got %v\nwant %v", got, want)
	}
}
