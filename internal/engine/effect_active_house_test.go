package engine

import "testing"

// TestChangeActiveHouseSetsActiveHouse covers the ChangeActiveHouse effect: the
// house of the card in context becomes the active player's active house.
func TestChangeActiveHouseSetsActiveHouse(t *testing.T) {
	if got := (ChangeActiveHouse{To: TheContextualHouse}).Text(); got != "its house becomes your active house" {
		t.Errorf("Text = %q", got)
	}

	g := started(t) // Brobnar is the active house
	it := g.AddToDeck(NewCard("scout", Logos, Creature, Common, WithPower(2)), 0)
	ChangeActiveHouse{
		To: TheContextualHouse,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, It: it, HasIt: true},
	)
	if g.State.ActiveHouse != Logos {
		t.Errorf("active house = %v, want Logos", g.State.ActiveHouse)
	}
}

// TestChangeActiveHouseNoContext covers the empty-context branch: with no card in
// context the active house is left untouched.
func TestChangeActiveHouseNoContext(t *testing.T) {
	g := started(t)
	ChangeActiveHouse{To: TheContextualHouse}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.ActiveHouse != Brobnar {
		t.Errorf("active house = %v, want Brobnar (unchanged)", g.State.ActiveHouse)
	}
}

// TestChangeActiveHouseValidateAndText covers the unset-source validation error and
// the chosen-house text rendering.
func TestChangeActiveHouseValidateAndText(t *testing.T) {
	if (ChangeActiveHouse{}).validate() == nil {
		t.Error("unset To should fail validation")
	}
	if (ChangeActiveHouse{To: TheContextualHouse}).validate() != nil {
		t.Error("TheContextualHouse should validate")
	}
	if got := (ChangeActiveHouse{To: TheChosenHouse}).Text(); got != "the chosen house becomes your active house" {
		t.Errorf("Text = %q", got)
	}
}
