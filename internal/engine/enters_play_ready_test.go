package engine

import "testing"

// A card with GrantsEntersReady makes friendly cards of that type enter play
// ready instead of exhausted while it is in play — Duskwitch for creatures, The
// Curator for artifacts.
func TestEntersPlayReady(t *testing.T) {
	t.Run("creature grant readies played creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		granter := g.AddToHand(
			NewCard("granter", Untamed, Creature, Common,
				WithPower(1), WithFriendlyEntersPlayReady(Creature)),
			0,
		)
		g.putIntoPlay(granter, 0)

		newbie := g.AddToHand(NewCard("newbie", Untamed, Creature, Common, WithPower(3)), 0)
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, newbie), false); err != nil {
			t.Fatalf("PlayCreature: %v", err)
		}
		if g.State.Cards[newbie].Exhausted {
			t.Error("creature should enter play ready under a creature grant")
		}
	})

	t.Run("artifact grant readies played artifacts", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		granter := g.AddToHand(
			NewCard("granter", Logos, Creature, Common,
				WithPower(1), WithFriendlyEntersPlayReady(Artifact)),
			0,
		)
		g.putIntoPlay(granter, 0)

		art := g.AddToHand(NewCard("art", Logos, Artifact, Common), 0)
		if _, err := g.PlayArtifact(0, handIdxByID(g, 0, art)); err != nil {
			t.Fatalf("PlayArtifact: %v", err)
		}
		if g.State.Cards[art].Exhausted {
			t.Error("artifact should enter play ready under an artifact grant")
		}
	})

	t.Run("an opponent's grant does not ready your card", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		granter := g.AddToHand(
			NewCard("granter", Logos, Creature, Common,
				WithPower(1), WithFriendlyEntersPlayReady(Artifact)),
			1,
		)
		g.putIntoPlay(granter, 1)

		art := g.AddToHand(NewCard("art", Logos, Artifact, Common), 0)
		if _, err := g.PlayArtifact(0, handIdxByID(g, 0, art)); err != nil {
			t.Fatalf("PlayArtifact: %v", err)
		}
		if !g.State.Cards[art].Exhausted {
			t.Error("artifact should enter exhausted: the grant belongs to the opponent")
		}
	})

	t.Run("without a grant cards enter play exhausted", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		newbie := g.AddToHand(NewCard("newbie", Untamed, Creature, Common, WithPower(3)), 0)
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, newbie), false); err != nil {
			t.Fatalf("PlayCreature: %v", err)
		}
		if !g.State.Cards[newbie].Exhausted {
			t.Error("creature should enter play exhausted without a grant")
		}
	})

	t.Run("renders the printed line for each type", func(t *testing.T) {
		creat := &CardDefinition{Name: "Duskwitch", GrantsEntersReady: Creature}
		if rules := cardRules(creat, false); len(rules) != 1 ||
			rules[0] != "Your Creatures enter play ready." {
			t.Errorf("creature rules = %v", rules)
		}
		art := &CardDefinition{Name: "The Curator", GrantsEntersReady: Artifact}
		if rules := cardRules(art, false); len(rules) != 1 ||
			rules[0] != "Friendly Artifacts enter play ready." {
			t.Errorf("artifact rules = %v", rules)
		}
	})
}
