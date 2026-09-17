package engine

import "testing"

// TestConstantBlankTextBlanksArtifacts checks a while-in-play BlankText constant
// ability blanks the artifacts its Target reaches, suppressing their abilities.
func TestConstantBlankTextBlanksArtifacts(t *testing.T) {
	g := started(t)
	// An enemy artifact grants every creature +2 power via a constant ability.
	buffer := g.AddArtifact(NewCard("buffer", Untamed, Artifact, Rare,
		WithConstantAbility(ConstantAbility{
			PowerBonus: 2,
			Target:     Target{Kind: TargetEachCreature},
		})), 1)
	victim := g.AddToBattleline(testCreature("victim", 3), 1)
	if g.Power(victim) != 5 {
		t.Fatalf("precondition: buffer should give +2 power, got %d", g.Power(victim))
	}

	drake := g.AddToBattleline(NewCard("drake", Untamed, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:    Target{Kind: TargetEachArtifact},
			BlankText: true,
		})), 0)

	if !g.textBlanked(buffer) {
		t.Error("the artifact's text box should be blanked")
	}
	if g.Power(victim) != 3 {
		t.Errorf("the blanked artifact's buff should vanish, power = %d", g.Power(victim))
	}

	// Blanking the drake itself (registry blank) lifts its constant blank, so the
	// artifact's text returns.
	g.BlankEnemyText(1) // player 1 blanks player 0's creatures, including the drake
	if !g.textBlanked(drake) {
		t.Fatal("the drake should be blanked by the registry")
	}
	if g.textBlanked(buffer) {
		t.Error("a blanked drake grants no blank, so the artifact's text returns")
	}
}
