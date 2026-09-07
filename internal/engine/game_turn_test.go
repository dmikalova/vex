package engine

import "testing"

// TestChooseHouseSkipsCardsRemovedMidWindow covers the AfterChooseHouse window
// when an earlier card's ability removes a later card from play: the snapshot is
// taken once, so the removed card must drop its own AfterChooseHouse trigger
// (ADR 0013) rather than fire it from the discard pile and strand in-play state
// (Æmber placed on a card no longer in play) — the invariant that caught it.
func TestChooseHouseSkipsCardsRemovedMidWindow(t *testing.T) {
	destroyer := NewCard("Purge Kin", Logos, Creature, Rare, WithPower(3),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Logos},
			Then: Destroy{Target: Target{Kind: TargetEachOtherFriendlyCreature}},
		}))
	hoarder := NewCard("Aember Sink", Logos, Creature, Rare, WithPower(3),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Logos},
			Then: PlaceAemberOnThis{Amount: 1},
		}))

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	// Order matters: the destroyer resolves first and removes the hoarder before
	// the loop reaches it.
	g.AddToBattleline(destroyer, 0)
	sink := g.AddToBattleline(hoarder, 0)

	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.inPlay(sink) {
		t.Fatalf("hoarder should have been destroyed")
	}
	if core := g.State.Cards[sink]; core != (CardCore{}) {
		t.Errorf("destroyed card carries in-play state: %+v", core)
	}
	if err := g.InvariantError(); err != nil {
		t.Errorf("invariant violated: %v", err)
	}
}
