package engine

import (
	"errors"
	"testing"
)

// TestMayPlayOrUseText covers the rendered clause for every axis combination the
// out-of-house permission family folds into MayPlayOrUse.
func TestMayPlayOrUseText(t *testing.T) {
	cases := []struct {
		name string
		e    MayPlayOrUse
		want string
	}{
		{
			"fight chosen house (Brothers in Battle)",
			MayPlayOrUse{Houses: HouseSelector{Match: chosenHouse}, Grant: GrantFight},
			"for the remainder of the turn, each friendly creature of the chosen house may fight",
		},
		{
			"fight named house (Signal Fire)",
			MayPlayOrUse{
				Houses: HouseSelector{Match: namedHouse(Brobnar)},
				Grant:  GrantFight,
			},
			"for the remainder of the turn, each friendly Brobnar creature may fight",
		},
		{
			"fight any house (Follow the Leader)",
			MayPlayOrUse{Houses: HouseSelector{Match: anyHouse}, Grant: GrantFight},
			"for the remainder of the turn, each friendly creature may fight",
		},
		{
			"use named house (Ritual of the Hunt)",
			MayPlayOrUse{Houses: HouseSelector{Match: namedHouse(Sanctum)}, Grant: GrantUse},
			"for the remainder of the turn, you may use friendly Sanctum creatures",
		},
		{
			"play or use named house (House Ambassador)",
			MayPlayOrUse{
				Houses: HouseSelector{Match: namedHouse(Mars)},
				Grant:  GrantPlay | GrantUse,
			},
			"for the remainder of the turn, you may play or use a Mars card",
		},
		{
			"play named house only",
			MayPlayOrUse{Houses: HouseSelector{Match: namedHouse(Mars)}, Grant: GrantPlay},
			"for the remainder of the turn, you may play a Mars card",
		},
		{
			"use artifacts any house (Scientifical Hack)",
			MayPlayOrUse{
				Houses: HouseSelector{Match: anyHouse},
				Grant:  GrantUse,
				Types:  CardTypesOf(Artifact),
			},
			"for the remainder of the turn, you may use friendly artifacts as if they belonged to the active house",
		},
		{
			"exclusion non-creature (Com. Officer Kirby)",
			MayPlayOrUse{
				Houses: HouseSelector{Match: exceptHouse(StarAlliance)},
				Grant:  GrantPlay,
				Types:  CardTypesOf(Artifact, Upgrade, Tactic),
				Cards:  1,
			},
			"you may play a non-Star Alliance artifact, upgrade, or tactic this turn",
		},
		{
			"exclusion play or use (CXO Taber)",
			MayPlayOrUse{
				Houses: HouseSelector{Match: exceptHouse(StarAlliance)},
				Grant:  GrantPlay | GrantUse,
				Cards:  1,
			},
			"you may play or use one non-Star Alliance card this turn",
		},
		{
			"exclusion no excluded house",
			MayPlayOrUse{
				Houses: HouseSelector{Match: HouseMatcher{Kind: MatchExceptHouse}},
				Grant:  GrantPlay,
				Cards:  1,
			},
			"you may play one card this turn",
		},
		{
			"controlled (United Action)",
			MayPlayOrUse{Houses: HouseSelector{Controlled: true}, Grant: GrantPlay},
			"for the remainder of the turn, you may play cards from any house for which you have a card in play",
		},
	}
	for _, c := range cases {
		if got := c.e.Text(); got != c.want {
			t.Errorf("%s: Text() = %q, want %q", c.name, got, c.want)
		}
	}
}

// TestMayPlayOrUseValidate covers each reason a grant is rejected. There is no
// "no houses" case: HouseSelector embeds a HouseMatcher whose zero value is
// MatchAnyHouse, so an unset selector is a valid any-house grant (Follow the
// Leader), not an error.
func TestMayPlayOrUseValidate(t *testing.T) {
	if (MayPlayOrUse{Houses: HouseSelector{Match: anyHouse}}).validate() == nil {
		t.Error("a grant with no verb should be invalid")
	}
	if (MayPlayOrUse{
		Houses: HouseSelector{Match: HouseMatcher{Kind: MatchExceptHouse}},
		Grant:  GrantPlay,
		Cards:  -1,
	}).validate() == nil {
		t.Error("a negative count should be invalid")
	}
	if (MayPlayOrUse{Houses: HouseSelector{Match: anyHouse}, Grant: GrantFight}).validate() != nil {
		t.Error("a set house and grant should be valid")
	}
}

// TestMayPlayOrUseResolveFight resolves the fight grants onto their state slots
// and confirms an off-house creature can then fight, reading the chosen house for
// the SelectHouse-with-no-house form.
func TestMayPlayOrUseResolveFight(t *testing.T) {
	g := NewGame("A", "B", 1)
	MayPlayOrUse{Houses: HouseSelector{Match: chosenHouse}, Grant: GrantFight}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Untamed},
	)
	if g.State.MayFightHouse[0] != Untamed {
		t.Errorf("MayFightHouse[0] = %v, want Untamed", g.State.MayFightHouse[0])
	}

	MayPlayOrUse{
		Houses: HouseSelector{Match: namedHouse(Brobnar)},
		Grant:  GrantFight,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 1},
	)
	if g.State.MayFightHouse[1] != Brobnar {
		t.Errorf("MayFightHouse[1] = %v, want Brobnar", g.State.MayFightHouse[1])
	}

	g2 := NewGame("A", "B", 1)
	g2.StartTurn(0)
	if err := g2.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	outsider := g2.AddToBattleline(NewCard("outsider", Logos, Creature, Common, WithPower(5)), 0)
	enemy := g2.AddToBattleline(NewCard("enemy", Dis, Creature, Common, WithPower(2)), 1)
	if err := g2.Fight(0, outsider, enemy); !errors.Is(err, ErrWrongHouse) {
		t.Fatalf("Fight before the grant = %v, want ErrWrongHouse", err)
	}
	MayPlayOrUse{Houses: HouseSelector{Match: anyHouse}, Grant: GrantFight}.Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
	if !g2.State.MayFightAny[0] {
		t.Error("MayFightAny[0] should be set")
	}
	if err := g2.Fight(0, outsider, enemy); err != nil {
		t.Fatalf("Fight after the grant = %v, want nil", err)
	}
}

// TestMayPlayOrUseResolveUsePlay resolves the named-house use/play grant and the
// artifacts-any-house grant, confirming each frees the off-house card its gate
// controls.
func TestMayPlayOrUseResolveUsePlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Sanctum); err != nil {
		t.Fatal(err)
	}
	off := NewCard("marauder", Mars, Creature, Common, WithPower(3))
	if g.mayPlayFromHand(0, &off) {
		t.Fatal("an off-house card should not be playable before the grant")
	}
	c := g.AddToBattleline(NewCard("cleric", Mars, Creature, Common, WithPower(3)), 0)
	if g.usableInActiveHouse(c) {
		t.Fatal("an off-house creature should not be usable before the grant")
	}
	MayPlayOrUse{
		Houses: HouseSelector{Match: namedHouse(Mars)},
		Grant:  GrantPlay | GrantUse,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.MayPlayHouse[0] != Mars || g.State.MayUseHouse[0] != Mars {
		t.Fatalf("grant should record Mars for play and use")
	}
	if !g.mayPlayFromHand(0, &off) {
		t.Error("the granted-house card should be playable")
	}
	if !g.usableInActiveHouse(c) {
		t.Error("the granted-house creature should be usable")
	}

	g2 := NewGame("A", "B", 1)
	g2.StartTurn(0)
	if err := g2.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	relic := g2.AddArtifact(NewCard("relic", Sanctum, Artifact, Common,
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1})), 0)
	if g2.usableInActiveHouse(relic) {
		t.Fatal("an off-house artifact should not be usable before the grant")
	}
	MayPlayOrUse{
		Houses: HouseSelector{Match: anyHouse},
		Grant:  GrantUse,
		Types:  CardTypesOf(Artifact),
	}.Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
	if !g2.State.MayUseArtifactsAnyHouse[0] {
		t.Error("MayUseArtifactsAnyHouse[0] should be set")
	}
	if !g2.usableInActiveHouse(relic) {
		t.Error("the friendly artifact should be usable after the grant")
	}
}

// TestMayPlayOrUseResolveTrait resolves the trait-scoped use grant and confirms it
// frees only friendly creatures of that trait, whatever their house.
func TestMayPlayOrUseResolveTrait(t *testing.T) {
	e := MayPlayOrUse{Trait: Mutant, Grant: GrantUse}
	if got, want := e.Text(),
		"for the remainder of the turn, you may use friendly Mutant creatures"; got != want {
		t.Errorf("Text() = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Sanctum); err != nil {
		t.Fatal(err)
	}
	mutant := g.AddToBattleline(
		NewCard("splicer", Mars, Creature, Common, WithPower(3), WithTraits(Mutant)), 0)
	beast := g.AddToBattleline(
		NewCard("marauder", Mars, Creature, Common, WithPower(3), WithTraits(Beast)), 0)
	if g.usableInActiveHouse(mutant) {
		t.Fatal("an off-house Mutant should not be usable before the grant")
	}

	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.MayUseTrait[0] != Mutant {
		t.Fatalf("grant should record Mutant, got %v", g.State.MayUseTrait[0])
	}
	if !g.usableInActiveHouse(mutant) {
		t.Error("the off-house Mutant should be usable after the grant")
	}
	if g.usableInActiveHouse(beast) {
		t.Error("a non-Mutant off-house creature should stay unusable")
	}
}

// TestMayPlayOrUseResolvePermit resolves the exclusion and controlled grants into
// stored off-house permits, tracking Remaining and clearing at end of turn.
func TestMayPlayOrUseResolvePermit(t *testing.T) {
	g := started(t)
	MayPlayOrUse{
		Houses: HouseSelector{Match: exceptHouse(StarAlliance)},
		Grant:  GrantPlay,
		Cards:  2,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.OffHousePermitCount[0] != 1 {
		t.Fatalf("permit count = %d, want 1", g.State.OffHousePermitCount[0])
	}
	if got := g.State.OffHousePermits[0][0]; got.Remaining != 2 || got.Except != StarAlliance {
		t.Errorf("bounded permit = %+v, want Remaining 2 Except StarAlliance", got)
	}

	MayPlayOrUse{Houses: HouseSelector{Controlled: true}, Grant: GrantPlay}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if got := g.State.OffHousePermits[0][1]; got.Remaining != permitUnlimited || !got.Controlled {
		t.Errorf("controlled permit = %+v, want unbounded and Controlled", got)
	}

	g.EndPlayPhase(0)
	g.StartTurn(0)
	if g.State.OffHousePermitCount[0] != 0 {
		t.Error("ready phase should clear off-house permits")
	}
}

// TestCardTypes covers the allowed-type bitset: the empty set admits everything, a
// narrowed set admits only its members, and list renders the rulebook phrase.
func TestCardTypes(t *testing.T) {
	if !(CardTypes(0)).all() {
		t.Error("the zero set should be all")
	}
	if !(CardTypes(0)).has(Creature) {
		t.Error("the empty set should admit every type")
	}
	if (CardTypes(0)).list() != "card" {
		t.Errorf("empty list = %q, want %q", CardTypes(0).list(), "card")
	}

	artifacts := CardTypesOf(Artifact)
	if !artifacts.has(Artifact) || artifacts.has(Creature) {
		t.Error("a single-type set should admit only its member")
	}
	if got := artifacts.list(); got != "artifact" {
		t.Errorf("single list = %q, want %q", got, "artifact")
	}

	if got := CardTypesOf(Artifact, Upgrade).list(); got != "artifact or upgrade" {
		t.Errorf("two-type list = %q, want %q", got, "artifact or upgrade")
	}
	if got := CardTypesOf(Artifact, Upgrade, Tactic).list(); got != "artifact, upgrade, or tactic" {
		t.Errorf("three-type list = %q, want %q", got, "artifact, upgrade, or tactic")
	}
}
