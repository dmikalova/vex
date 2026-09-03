package engine

import "testing"

// TestDamageSkipsACardOutOfPlay checks damage aimed at a creature an earlier step
// already removed is dropped rather than marked on the card in its discard pile.
func TestDamageSkipsACardOutOfPlay(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(NewCard("Oak", Untamed, Creature, Common, WithPower(5)), 0)
	g.putIntoHand(victim)

	g.applyRawDamage(victim, 3, true)

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
	g.applyRawDamage(marauder, 3, true)

	if !g.inPlay(marauder) {
		t.Fatal("3 damage should not destroy it at 4 power")
	}
	g.addAmberOn(marauder, -1)

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
	g.applyRawDamage(id, 6, true)

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
