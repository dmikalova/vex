package engine

import "testing"

func TestBlankEnemyTextEffect(t *testing.T) {
	if got := (BlankEnemyText{}).Text(); got !=
		"until your next turn, enemy creatures' text boxes are "+
			"considered blank (except for traits)" {
		t.Errorf("text = %q", got)
	}

	g := started(t)
	// Player 1's creature has taunt and combat keywords; blank it via Shadow of Dis.
	foe := g.AddToBattleline(
		testCreature("foe", 5, WithKeywords(Taunt),
			WithAssault(2), WithHazardous(3), WithSplashAttack(1)), 1)
	if !g.hasKeyword(foe, Taunt) {
		t.Fatal("precondition: foe should have taunt before blanking")
	}

	BlankEnemyText{}.Resolve(&EffectContext{Resolver: g, Source: foe, Controller: 0})

	if !g.textBlanked(foe) {
		t.Error("foe should be blanked")
	}
	if g.hasKeyword(foe, Taunt) {
		t.Error("blanked creature should lose its printed keyword")
	}
	if g.assault(foe) != 0 {
		t.Errorf("blanked creature assault = %d, want 0", g.assault(foe))
	}
	if g.hazardous(foe) != 0 {
		t.Errorf("blanked creature hazardous = %d, want 0", g.hazardous(foe))
	}
	if g.splashAttack(foe) != 0 {
		t.Errorf("blanked creature splash = %d, want 0", g.splashAttack(foe))
	}
}

func TestBlankEnemyTextSuppressesAbilitiesAndConstants(t *testing.T) {
	g := started(t)

	// A blanked creature's own triggered ability does not fire.
	fired := false
	spy := recordingEffect{fn: func() { fired = true }}
	foe := g.AddToBattleline(
		NewCard("chatterbox", Dis, Creature, Common, WithPower(3),
			WithAbility(TriggerAfterReap, spy)), 1)
	g.BlankEnemyText(0)
	g.triggerAbilities(foe, TriggerAfterReap, 0, false)
	if fired {
		t.Error("blanked creature's own ability should not fire")
	}

	// A blanked creature's constant ability grants nothing.
	granter := g.AddToBattleline(
		NewCard("granter", Dis, Creature, Common, WithPower(4),
			WithConstantAbility(ConstantAbility{
				PowerBonus: 2,
				Target:     Target{Kind: TargetEachFriendlyCreature},
			})), 1)
	ally := g.AddToBattleline(testCreature("ally", 3), 1)
	g.BlankEnemyText(0) // player 1's creatures, including granter, are blank
	if got := g.Power(ally); got != 3 {
		t.Errorf("ally power with blanked granter = %d, want 3", got)
	}
	if g.constantActive(granter, g.cat.def(granter).ConstantAbilities[0]) {
		t.Error("a blanked creature's constant ability should be inactive")
	}
}

// Blanking a creature that buffs its own power with a constant ability drops it to
// its base power, which can meet the damage already on it — the resolution
// boundary must then destroy it rather than leave a lethal creature in play
// (Mushroom Man, Marmo Swarm under Shadow of Dis).
func TestBlankEnemyTextSettlesLethal(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(
		NewCard("selfbuff", Dis, Creature, Common, WithPower(2),
			WithConstantAbility(ConstantAbility{
				PowerBonus: 4,
				Target:     Target{Kind: TargetThisCreature},
			})), 1)
	g.DealDamage(0, []DamageTarget{{ID: foe, Amount: 5}})
	if !g.inPlay(foe) {
		t.Fatal("5 damage should not destroy it at 6 power")
	}

	g.BlankEnemyText(0)
	g.settleDestroyed(0) // the resolution boundary settles the blank (ADR 0029)

	if g.inPlay(foe) {
		t.Errorf("blanking to base 2 power under 5 damage should destroy it")
	}
}

func TestBlankEnemyTextIgnoresNonCreatures(t *testing.T) {
	g := started(t)
	art := g.AddArtifact(NewCard("relic", Dis, Artifact, Common), 1)
	g.BlankEnemyText(0)
	if g.textBlanked(art) {
		t.Error("an artifact is not a creature and is never blanked")
	}
}

func TestBlankEnemyTextLiftsAfterOpponentTurn(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("foe", 5, WithKeywords(Taunt)), 1)
	g.BlankEnemyText(0)

	// Still blank through the controller's turn.
	g.EndPlayPhase(0)
	if !g.textBlanked(foe) {
		t.Error("blank should persist into the opponent's turn")
	}
	// Blank lifts once the opponent's own turn ends (their ready phase clears it).
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	g.EndPlayPhase(1)
	if g.textBlanked(foe) {
		t.Error("blank should lift after the opponent's turn")
	}
}

// A blank can suppress an enemy constant power-reducer (Shadow of Dis blanking
// King of the Crag). When the blank lifts at the very end of the turn, the
// reducer reactivates and can drop a creature to zero power — that creature must
// be settled before the turn hands over, not left in play at 0 power.
func TestBlankLiftingSettlesLethalConstant(t *testing.T) {
	g := started(t)
	// Player 0's victim sits at base power 2.
	victim := g.AddToBattleline(testCreature("victim", 2), 0)
	// Player 0 blanks player 1's creatures before the reducer arrives, so the
	// reducer's -2 is suppressed and the victim keeps its 2 power.
	g.BlankEnemyText(0)
	g.AddToBattleline(
		NewCard("reducer", Dis, Creature, Common, WithPower(4),
			WithConstantAbility(ConstantAbility{
				PowerBonus: -2,
				Target:     Target{Kind: TargetEachEnemyCreature},
			})), 1)
	if !g.inPlay(victim) {
		t.Fatal("precondition: victim should survive while the reducer is blanked")
	}

	// The blank persists through the opponent's turn, then lifts at its end.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	g.EndPlayPhase(1)

	if g.inPlay(victim) {
		t.Error("victim dropped to 0 power when the blank lifted; it must be destroyed")
	}
}

// recordingEffect is a minimal effect that runs fn when it resolves, used to
// observe whether an ability fired.
type recordingEffect struct{ fn func() }

func (recordingEffect) Text() string             { return "record" }
func (e recordingEffect) Resolve(*EffectContext) { e.fn() }
