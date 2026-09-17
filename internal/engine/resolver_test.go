package engine

import "testing"

// TestSetAemberClampsAtZero exercises the Resolver's Æmber setter, including its
// floor at zero (a pool can never go negative).
func TestSetAemberClampsAtZero(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetAember(0, 5)
	if g.Aember(0) != 5 {
		t.Errorf("Aember = %d, want 5", g.Aember(0))
	}
	g.SetAember(0, -3) // negative clamps to zero
	if g.Aember(0) != 0 {
		t.Errorf("Aember after negative set = %d, want 0", g.Aember(0))
	}
}

// TestNoteAemberStolenFromClamps covers the running steal tally: it ignores
// non-positive amounts and clamps at the int8 max the history holds.
func TestNoteAemberStolenFromClamps(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.NoteAemberStolenFrom(0, 0) // ignored
	g.NoteAemberStolenFrom(0, 100)
	g.NoteAemberStolenFrom(0, 100) // 200 clamps to 127
	if got := g.State.TurnHistory[0][AemberStolenFromThisTurn]; got != 127 {
		t.Errorf("tally = %d, want 127 (clamped)", got)
	}
}

// An ability keeps resolving after its target has left play, so every writer of
// in-play state has to land on nothing rather than leave counters or a stun on a
// card sitting in the discard pile.
func TestInPlayStateWritesSkipCardsOutOfPlay(t *testing.T) {
	g := started(t)
	gone := g.Register(testCreature("gone", 3), 0)
	g.State.Discard[0].add(gone)

	g.SetDamage(gone, 2)
	g.SetStunned(gone, true)
	g.SetWarded(gone, true)
	g.SetDamageImmune(gone, RemainderOfPlayerTurn)
	g.SetExhausted(gone, true)
	g.BelongToHouseForRemainderOfTurn(gone, Dis)
	g.SetLastingHouse(gone, Dis)
	g.SetNamedHouse(gone, Dis)
	g.AddPowerCounter(gone, 3)

	if core := g.State.Cards[gone]; core != (CardCore{}) {
		t.Errorf("core = %+v, want zero: a card out of play takes no in-play state", core)
	}
}
