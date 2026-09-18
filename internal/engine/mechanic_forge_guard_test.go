package engine

import (
	"strings"
	"testing"
)

// keyforgeryCard is Keyforgery: an artifact that interrupts the opponent's key
// forges through a "when your opponent would forge a key" ability.
func keyforgeryCard() CardDefinition {
	return NewCard("Keyforgery", Shadows, Artifact, Rare, WithAbility(
		TriggerBeforeOpponentForgesKey,
		Sentences{Effects: []Effect{
			OpponentNamesHouse{},
			RevealRandomFromHand{},
			Conditional{
				Cond: ItIsNotOfNamedHouse{Noun: ThatCard},
				Then: Sequence{Effects: []Effect{
					Destroy{Target: Target{Kind: TargetThisCreature}},
					CancelForge{},
				}},
			},
		}},
	))
}

func TestForgeGuardText(t *testing.T) {
	def := keyforgeryCard()
	want := "When your opponent would forge a key, they name a house. Reveal a " +
		"random card from your hand. If that card is not of the named house, " +
		"destroy Keyforgery, and they do not forge that key."
	if !strings.Contains(RenderCardRules(&def), want) {
		t.Errorf(
			"card rules should render the before-forge ability, got:\n%s",
			RenderCardRules(&def),
		)
	}
}

// TestForgeGuardPreventsForge covers a wrong guess: the opponent names a house the
// revealed card is not, so the guard is destroyed and the forge is prevented with
// no Æmber spent.
func TestForgeGuardPreventsForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddArtifact(NewCard("plain", Dis, Artifact, Common), 1) // a non-guard card is skipped
	kf := g.AddArtifact(keyforgeryCard(), 1)
	g.AddToHand(NewCard("Logos Card", Logos, Creature, Common, WithPower(3)), 1)
	g.State.Aember[0] = 6
	g.SetChooser(0, optionPicker{idx: 0}) // name Brobnar; the reveal is Logos

	keysBefore := g.Keys(0)
	g.forgePhase(0)

	if g.Keys(0) != keysBefore {
		t.Error("a prevented forge should not forge a key")
	}
	if g.State.Aember[0] != 6 {
		t.Errorf("a prevented forge should spend no Æmber, pool = %d", g.State.Aember[0])
	}
	if g.inPlay(kf) {
		t.Error("a wrong guess should destroy Keyforgery")
	}
}

// TestForgeGuardAllowsForge covers a correct guess: the opponent names the house
// of the revealed card, so the guard stays in play and the forge proceeds.
func TestForgeGuardAllowsForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	kf := g.AddArtifact(keyforgeryCard(), 1)
	g.AddToHand(NewCard("Logos Card", Logos, Creature, Common, WithPower(3)), 1)
	g.State.Aember[0] = 6
	g.SetChooser(0, optionPicker{idx: 2}) // name Logos; the reveal is Logos

	g.forgePhase(0)

	if g.Keys(0) != 1 {
		t.Error("a correct guess should let the forge proceed")
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("a completed forge should spend the Æmber, pool = %d", g.State.Aember[0])
	}
	if !g.inPlay(kf) {
		t.Error("a correct guess should leave Keyforgery in play")
	}
}

// TestForgeGuardEmptyHand covers a guard whose controller has no card to reveal:
// it cannot act, so the forge proceeds and the guard survives.
func TestForgeGuardEmptyHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	kf := g.AddArtifact(keyforgeryCard(), 1)
	g.State.Aember[0] = 6
	g.SetChooser(0, optionPicker{idx: 0})

	g.forgePhase(0)

	if g.Keys(0) != 1 {
		t.Error("a guard with an empty hand should not prevent the forge")
	}
	if !g.inPlay(kf) {
		t.Error("a guard that cannot reveal should stay in play")
	}
}

// TestForgeGuardOnlyOpponent covers that the guard interrupts the opponent's
// forges, never its own controller's.
func TestForgeGuardOnlyOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	kf := g.AddArtifact(keyforgeryCard(), 0)
	g.AddToHand(NewCard("Logos Card", Logos, Creature, Common, WithPower(3)), 0)
	g.State.Aember[0] = 6
	g.SetChooser(0, optionPicker{idx: 0})

	g.forgePhase(0)

	if g.Keys(0) != 1 {
		t.Error("a controller's own Keyforgery should not prevent their forge")
	}
	if !g.inPlay(kf) {
		t.Error("a controller's own Keyforgery should stay in play")
	}
}

// TestForgeGuardOnFreeForge covers a guard interrupting a free (unpaid) forge:
// Keyforgery still prevents it on a wrong guess.
func TestForgeGuardOnFreeForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	kf := g.AddArtifact(keyforgeryCard(), 1)
	g.AddToHand(NewCard("Logos Card", Logos, Creature, Common, WithPower(3)), 1)
	g.SetChooser(0, optionPicker{idx: 0}) // name Brobnar; the reveal is Logos

	g.forgeKeyFree(0)

	if g.Keys(0) != 0 {
		t.Error("a prevented free forge should not forge a key")
	}
	if g.inPlay(kf) {
		t.Error("a wrong guess should destroy Keyforgery on a free forge")
	}
}
