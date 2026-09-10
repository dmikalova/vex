package engine

import "testing"

// TestMakeItsHouseActiveSetsActiveHouse covers the MakeItsHouseActive effect: the
// house of the card in context becomes the active player's active house.
func TestMakeItsHouseActiveSetsActiveHouse(t *testing.T) {
	if got := (MakeItsHouseActive{}).Text(); got != "its house becomes your active house" {
		t.Errorf("Text = %q", got)
	}

	g := started(t) // Brobnar is the active house
	it := g.AddToDeck(NewCard("scout", Logos, Creature, Common, WithPower(2)), 0)
	MakeItsHouseActive{}.Resolve(&EffectContext{Resolver: g, Controller: 0, It: it, HasIt: true})
	if g.State.ActiveHouse != Logos {
		t.Errorf("active house = %v, want Logos", g.State.ActiveHouse)
	}
}

// TestMakeItsHouseActiveNoContext covers the empty-context branch: with no card in
// context the active house is left untouched.
func TestMakeItsHouseActiveNoContext(t *testing.T) {
	g := started(t)
	MakeItsHouseActive{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.ActiveHouse != Brobnar {
		t.Errorf("active house = %v, want Brobnar (unchanged)", g.State.ActiveHouse)
	}
}
