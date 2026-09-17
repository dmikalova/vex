package engine

import "testing"

// banner returns a creature whose constant ability gives every friendly creature
// the given power bonus.
func banner(bonus int) CardDefinition {
	return NewCard("Banner", Untamed, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachFriendlyCreature},
			PowerBonus: bonus,
		}))
}

// TestZeroPowerIsDestroyed checks a creature left at 0 power by the loss of a
// buff is destroyed, without anything dealing it damage.
func TestZeroPowerIsDestroyed(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	victim := g.AddToBattleline(NewCard("Sapling", Untamed, Creature, Common), 0)

	if got := g.Power(victim); got != 2 {
		t.Fatalf("power = %d, want 2 while the banner is in play", got)
	}
	g.putIntoHand(src)
	g.settleDestroyed(0) // the resolution boundary settles the buff loss (ADR 0029)

	if g.inPlay(victim) {
		t.Errorf("a creature at 0 power should have been destroyed")
	}
}

// TestArtifactSelfDestroysWhenNoCreatures covers Doom Sigil: an artifact carrying
// a DestroyedWhen condition holds while a creature is in play and destroys itself
// once the board empties. Only artifacts in play with a met condition qualify.
func TestArtifactSelfDestroysWhenNoCreatures(t *testing.T) {
	g := started(t)
	sigil := g.AddArtifact(
		NewCard("Doom Sigil", Shadows, Artifact, Rare,
			WithDestroyedWhen(InPlay{Player: EachPlayer, Type: Creature, None: true})), 0)
	creature := g.AddToBattleline(
		NewCard("Sapling", Untamed, Creature, Common, WithPower(2)), 0)

	// A creature is in play, so the artifact holds.
	g.settleDestroyed(0)
	if !g.inPlay(sigil) {
		t.Fatal("Doom Sigil should survive while a creature is in play")
	}
	// A plain artifact with no DestroyedWhen never self-destroys, and a creature is
	// never an artifact self-destroy candidate.
	plain := g.AddArtifact(NewCard("Plain", Shadows, Artifact, Common), 0)
	if g.artifactShouldSelfDestroy(plain) {
		t.Error("an artifact with no DestroyedWhen should not self-destroy")
	}
	if g.artifactShouldSelfDestroy(creature) {
		t.Error("a creature is not an artifact self-destroy candidate")
	}

	// The board empties: the artifact destroys itself, and the plain one survives.
	g.putIntoHand(creature)
	g.settleDestroyed(0)
	if g.inPlay(sigil) {
		t.Error("Doom Sigil should self-destroy once no creatures remain")
	}
	if !g.inPlay(plain) {
		t.Error("a plain artifact should survive an empty board")
	}
	// An out-of-play artifact is not a candidate.
	if g.artifactShouldSelfDestroy(sigil) {
		t.Error("a destroyed artifact should not be a self-destroy candidate")
	}
}

// TestBuffLossKillsADamagedCreature checks a damaged creature is destroyed once a
// lost buff leaves its damage at or above its remaining power.
func TestBuffLossKillsADamagedCreature(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	victim := g.AddToBattleline(NewCard("Oak", Untamed, Creature, Common, WithPower(3)), 0)
	g.applyRawDamage(DamageTarget{ID: victim, Amount: 4, IgnoreArmor: true})

	if !g.inPlay(victim) {
		t.Fatal("4 damage should not destroy a 5-power creature")
	}
	g.putIntoHand(src)
	g.settleDestroyed(0) // the resolution boundary settles the buff loss (ADR 0029)

	if g.inPlay(victim) {
		t.Errorf("damage at or above the remaining power should destroy the creature")
	}
}

// TestSettleCascades checks the sweep repeats: the creature it destroys was
// itself buffing another, which then dies in the same settling.
func TestSettleCascades(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(banner(2), 0)
	middle := g.AddToBattleline(banner(2), 0)
	last := g.AddToBattleline(NewCard("Sprout", Untamed, Creature, Common), 0)

	// middle sits at 3 printed + 2 from src + 2 from itself; last is 0 + 4.
	g.applyRawDamage(DamageTarget{ID: middle, Amount: 6, IgnoreArmor: true})
	g.putIntoHand(src)
	g.settleDestroyed(0) // the resolution boundary settles the buff loss (ADR 0029)

	if g.inPlay(middle) {
		t.Errorf("the damaged buffer should have died once its own buff was gone")
	}
	if g.inPlay(last) {
		t.Errorf("the creature the dead buffer was propping up should have died too")
	}
}

// TestDamageSkipsACardOutOfPlay checks damage aimed at a creature an earlier step
// already removed is dropped rather than marked on the card in its discard pile.
func TestDamageSkipsACardOutOfPlay(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(NewCard("Oak", Untamed, Creature, Common, WithPower(5)), 0)
	g.putIntoHand(victim)

	g.applyRawDamage(DamageTarget{ID: victim, Amount: 3, IgnoreArmor: true})

	if got := g.State.Cards[victim].Damage; got != 0 {
		t.Errorf("damage on a card in hand = %d, want 0", got)
	}
}

// TestCaptureSkipsACreatureOutOfPlay checks a creature destroyed by the fight it
// captured on keeps the Æmber in the opponent's pool instead of banking it on a
// card in a discard pile.
func TestCaptureSkipsACreatureOutOfPlay(t *testing.T) {
	g := started(t)
	captor := g.AddToBattleline(NewCard("Raider", Mars, Creature, Common, WithPower(3)), 0)
	g.State.Aember[1] = 3
	g.putIntoHand(captor)

	CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent}.
		Resolve(&EffectContext{Resolver: g, Source: captor, Controller: 0})

	if got := g.State.Aember[1]; got != 3 {
		t.Errorf("opponent pool = %d, want 3 left alone", got)
	}
	if got := g.State.Cards[captor].Amber; got != 0 {
		t.Errorf("Æmber on a card in hand = %d, want 0", got)
	}
}

// TestAemberOnACreatureFeedsItsPower checks the sweep runs when Æmber leaves a
// creature that draws its power from the Æmber sitting on it.
func TestAemberOnACreatureFeedsItsPower(t *testing.T) {
	g := started(t)
	marauder := g.AddToBattleline(
		NewCard("Marauder", Mars, Creature, Common, WithPower(2),
			WithConstantAbility(ConstantAbility{
				Target:     Target{Kind: TargetThisCreature},
				PowerBonus: 1,
				Per:        AemberOnThis{},
			})),
		0,
	)
	g.addAmberOn(marauder, 2)
	g.applyRawDamage(DamageTarget{ID: marauder, Amount: 3, IgnoreArmor: true})

	if !g.inPlay(marauder) {
		t.Fatal("3 damage should not destroy it at 4 power")
	}
	g.addAmberOn(marauder, -1)
	g.settleDestroyed(0) // the resolution boundary settles the loss (ADR 0029)

	if g.inPlay(marauder) {
		t.Errorf("losing the Æmber should drop it to 3 power and destroy it")
	}
}

// TestEnteringPlayUnderAPowerPenaltyIsDestroyed checks a creature that arrives
// into a constant ability leaving it at 0 power is destroyed on arrival.
func TestEnteringPlayUnderAPowerPenaltyIsDestroyed(t *testing.T) {
	g := started(t)
	g.AddToBattleline(
		NewCard("Crag King", Brobnar, Creature, Common, WithPower(8),
			WithConstantAbility(ConstantAbility{
				Target:     Target{Kind: TargetEachCreature},
				PowerBonus: -2,
			})),
		0,
	)
	id := g.AddToHand(NewCard("Sapling", Untamed, Creature, Common, WithPower(2)), 0)
	g.State.ActivePlayer = 0
	g.State.ActiveHouse = Untamed
	g.PlayFromHand(0, id)

	if g.inPlay(id) {
		t.Errorf("a 2-power creature entering a -2 power penalty should be destroyed")
	}
}

// TestForgingSettlesPowerFromUnforgedKeys checks the sweep runs on a forge, for
// creatures whose power counts the keys their controller has yet to forge.
func TestForgingSettlesPowerFromUnforgedKeys(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(
		NewCard("Fungus", Untamed, Creature, Common, WithPower(2),
			WithConstantAbility(ConstantAbility{
				Target:     Target{Kind: TargetThisCreature},
				PowerBonus: 3,
				Per:        UnforgedKeys{Player: Controller},
			})),
		0,
	)
	g.applyRawDamage(DamageTarget{ID: id, Amount: 6, IgnoreArmor: true})

	if !g.inPlay(id) {
		t.Fatal("6 damage should not destroy it while it holds 11 power")
	}
	g.forgeKeyFree(0)

	if !g.inPlay(id) {
		t.Errorf("forging should drop it to 8 power — still alive")
	}
	g.forgeKeyFree(0)
	if g.inPlay(id) {
		t.Errorf("a second forge should drop it to 5 power and destroy it")
	}
}

// recordInPlay is a test effect that records whether a given card is still in play
// at the moment it resolves.
type recordInPlay struct {
	id  LocalID
	got *bool
}

func (recordInPlay) Text() string { return "record in play" }
func (e recordInPlay) Resolve(ctx *EffectContext) {
	*e.got = ctx.Resolver.InPlay(e.id)
}

// TestArrivalKillsFlankNeighborBeforeAfterPlay pins the timing: a creature that
// loses its "+2 while on a flank" bonus because an arriving creature pushes it
// interior dies from that state-based check the instant it is placed — before the
// played creature's own after-play ability resolves, so that ability never sees
// the doomed neighbor still in play (Buzzle purging a neighbor it just doomed).
func TestArrivalKillsFlankNeighborBeforeAfterPlay(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("left", 3), 0)
	neighbor := g.AddToBattleline(
		testCreature("flank guard", 4,
			WithConstantAbility(ConstantAbility{
				Target:     Target{Kind: TargetThisCreature}.OnFlank(),
				PowerBonus: 2,
			})),
		0,
	)
	g.State.Cards[neighbor].Damage = 5 // lethal at 4 power, survives at 6 on a flank

	var neighborStillInPlay bool
	id := g.AddToHand(
		testCreature("arriving", 3,
			WithAbility(TriggerAfterPlay, recordInPlay{id: neighbor, got: &neighborStillInPlay})),
		0,
	)
	g.State.ActivePlayer = 0
	g.State.ActiveHouse = Brobnar
	g.SetChooser(0, optionPicker{idx: 1}) // land the arrival on the right flank
	g.PlayFromHand(0, id)

	if neighborStillInPlay {
		t.Error("the after-play ability saw the doomed neighbor still in play; " +
			"it should have died on arrival, before the after-play window")
	}
	if g.inPlay(neighbor) {
		t.Error("the neighbor should be destroyed once it is pushed off its flank")
	}
}

// TestTreacheryHandoffSettlesTheNewController pins that handing a Treachery
// creature to the opponent settles the board it lands in: the active player lands
// the seized creature on the new controller's right flank, pushing that side's
// flank creature interior, so a creature that loses its "+2 while on a flank" bonus
// and drops to or below its damage must be destroyed in the same action, not left
// lingering.
func TestTreacheryHandoffSettlesTheNewController(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("left", 3), 0)
	guard := g.AddToBattleline(
		testCreature("flank guard", 4,
			WithConstantAbility(ConstantAbility{
				Target:     Target{Kind: TargetThisCreature}.OnFlank(),
				PowerBonus: 2,
			})),
		0,
	)
	g.State.Cards[guard].Damage = 5 // lethal at 4 power, survives at 6 on a flank

	seized := g.AddToHand(testCreature("treachery", 3, WithKeywords(Treachery)), 1)
	g.State.ActivePlayer = 1
	g.State.ActiveHouse = Brobnar
	g.SetChooser(1, optionPicker{idx: 1}) // the active player lands it on the right flank
	g.PlayFromHand(1, seized)

	if g.controller(seized) != 0 {
		t.Fatal("a Treachery creature should enter under the opponent's control")
	}
	if g.inPlay(guard) {
		t.Error("the guard should be destroyed once the handoff pushes it off its flank")
	}
}

// TestPlaceAemberOnACardOutOfPlayLandsOnNothing checks an ability that places
// Æmber on its own source after that source has left play banks nothing on the
// card in its discard pile (Strange Gizmo forging a key mid-window).
func TestPlaceAemberOnACardOutOfPlayLandsOnNothing(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(NewCard("Gizmo", Logos, Creature, Common, WithPower(3)), 0)
	g.putIntoHand(src)

	PlaceAemberOnThis{Amount: 2}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

	if got := g.State.Cards[src].Amber; got != 0 {
		t.Errorf("Æmber on a card in hand = %d, want 0", got)
	}
}
