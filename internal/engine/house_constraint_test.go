package engine

import (
	"errors"
	"slices"
	"testing"
)

// TestNoActiveHouseWhenAllForbidden covers the empty allowed set: when every
// choosable house is barred, the player has no active house and No House
// (HouseNone) is the only legal choice (ADR 0035).
func TestNoActiveHouseWhenAllForbidden(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Mars})
	g.State.ActivePlayer = 0
	g.State.HouseConstraints[0][0] = HouseConstraint{
		Kind:  constraintCannotHouse,
		House: Mars,
	}
	g.State.HouseConstraintCount[0] = 1
	if got := g.AllowedHouses(0); len(got) != 0 {
		t.Fatalf("allowed = %v, want none", got)
	}
	if err := g.ChooseHouse(0, Mars); !errors.Is(err, ErrHouseNotAllowed) {
		t.Errorf("choosing a barred house = %v, want ErrHouseNotAllowed", err)
	}
	if err := g.ChooseHouse(0, HouseNone); err != nil {
		t.Errorf("No House = %v, want nil", err)
	}
}

// TestNoHouseRejectedWhenHouseAvailable covers the other side: No House is illegal
// while any house is still choosable.
func TestNoHouseRejectedWhenHouseAvailable(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Mars, Logos, Untamed})
	g.State.ActivePlayer = 0
	if err := g.ChooseHouse(0, HouseNone); !errors.Is(err, ErrHouseNotAllowed) {
		t.Errorf("No House with houses available = %v, want ErrHouseNotAllowed", err)
	}
}

// TestChoosableHousesIncludeControlledCards covers the widened choosable set: a
// house on a card the player controls in play in its own right — a creature or an
// artifact — is choosable even when it is off their identity (rulebook line 454).
// A controlled card duplicating an identity house exercises the dedup path.
func TestChoosableHousesIncludeControlledCards(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Mars})
	g.State.ActivePlayer = 0
	g.AddToBattleline(NewCard("marsC", Mars, Creature, Common, WithPower(1)), 0)
	g.AddArtifact(NewCard("logosA", Logos, Artifact, Common), 0)
	got := g.AllowedHouses(0)
	if !slices.Contains(got, Mars) || !slices.Contains(got, Logos) {
		t.Fatalf("allowed = %v, want Mars and Logos", got)
	}
	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Errorf("choosing the controlled off-house = %v, want nil", err)
	}
}

// TestHouseConstraintsStack covers the stacking the single slots could not express:
// must Mars, must Logos, cannot Mars leaves exactly {Logos}.
func TestHouseConstraintsStack(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Mars, Logos, Untamed})
	g.State.ActivePlayer = 0
	g.State.HouseConstraints[0][0] = HouseConstraint{
		Kind:  constraintMustHouse,
		House: Mars,
	}
	g.State.HouseConstraints[0][1] = HouseConstraint{
		Kind:  constraintMustHouse,
		House: Logos,
	}
	g.State.HouseConstraints[0][2] = HouseConstraint{
		Kind:  constraintCannotHouse,
		House: Mars,
	}
	g.State.HouseConstraintCount[0] = 3
	got := g.AllowedHouses(0)
	if len(got) != 1 || got[0] != Logos {
		t.Fatalf("allowed = %v, want [Logos] (must Mars, must Logos, cannot Mars)", got)
	}
}

// TestMustChooseFoughtCreatureHouse covers the must-creature constraint: the
// house is read dynamically from the referenced creature at choice time, not
// frozen when the constraint was armed (ADR 0035).
func TestMustChooseFoughtCreatureHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Mars, Logos, Untamed})
	g.State.ActivePlayer = 0
	id := g.AddToBattleline(NewCard("logosC", Logos, Creature, Common, WithPower(1)), 0)
	g.State.HouseConstraints[0][0] = HouseConstraint{
		Kind:     constraintMustCreature,
		Creature: id,
	}
	g.State.HouseConstraintCount[0] = 1
	got := g.AllowedHouses(0)
	if len(got) != 1 || got[0] != Logos {
		t.Fatalf("allowed = %v, want [Logos] (must the fought creature's house)", got)
	}
}

// TestHouseConstraintTableFull covers the table cap: entries armed past
// maxHouseConstraints are dropped rather than growing the flat array.
func TestHouseConstraintTableFull(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i <= maxHouseConstraints; i++ {
		g.CannotChooseHouseNextTurn(0, Mars, 0)
	}
	if got := int(g.State.HouseConstraintCountNext[0]); got != maxHouseConstraints {
		t.Errorf("count = %d, want %d (overflow dropped)", got, maxHouseConstraints)
	}
}
