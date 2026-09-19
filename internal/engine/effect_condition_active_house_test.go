package engine

import "testing"

// TestActiveHouseMatchesNoCardsInPlay covers the condition behind Sci. Officer
// Qincan: it is met only when no card in play — either player's creatures, their
// upgrades, or artifacts — belongs to the chosen active house.
func TestActiveHouseMatchesNoCardsInPlay(t *testing.T) {
	c := ActiveHouseMatchesNoCardsInPlay{}
	if got := c.CondText(); got != "which matches no cards in play" {
		t.Errorf("CondText = %q", got)
	}

	// The ability folds into a single readable sentence under its trigger.
	a := Ability{Trigger: TriggerAfterChooseHouse, EachPlayer: true, Effect: Conditional{
		Cond: c,
		Then: StealAember{Amount: 1},
	}}
	const want = "After a player chooses an active house which matches no cards in play, steal 1 Æmber."
	if got := RenderAbility(a); got != want {
		t.Errorf("render = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(NewCard("marine", StarAlliance, Creature, Common, WithPower(2)), 0)
	attachUpgrade(g, host, NewCard("charm", Logos, Upgrade, Common))
	g.AddArtifact(NewCard("relic", Shadows, Artifact, Common), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	// A creature of the active house is a match, so the condition is not met.
	g.State.ActiveHouse = StarAlliance
	if c.Met(ctx) {
		t.Error("a matching creature should make the condition not met")
	}

	// An upgrade of the active house (on a creature of another house) still matches.
	g.State.ActiveHouse = Logos
	if c.Met(ctx) {
		t.Error("a matching upgrade should make the condition not met")
	}

	// An enemy artifact of the active house matches.
	g.State.ActiveHouse = Shadows
	if c.Met(ctx) {
		t.Error("a matching artifact should make the condition not met")
	}

	// A house that nothing in play shares leaves the condition met.
	g.State.ActiveHouse = Brobnar
	if !c.Met(ctx) {
		t.Error("a house matching no card should make the condition met")
	}
}
