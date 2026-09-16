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
				WithPower(1), WithFriendlyEntersPlayReady(EntersReadyGrant{Type: Creature})),
			0,
		)
		g.putIntoPlay(granter, 0)

		newbie := g.AddToHand(NewCard("newbie", Untamed, Creature, Common, WithPower(3)), 0)
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, newbie), false); err != nil {
			t.Fatalf("Playcreature: %v", err)
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
				WithPower(1), WithFriendlyEntersPlayReady(EntersReadyGrant{Type: Artifact})),
			0,
		)
		g.putIntoPlay(granter, 0)

		art := g.AddToHand(NewCard("art", Logos, Artifact, Common), 0)
		if _, err := g.PlayArtifact(0, handIdxByID(g, 0, art)); err != nil {
			t.Fatalf("Playartifact: %v", err)
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
				WithPower(1), WithFriendlyEntersPlayReady(EntersReadyGrant{Type: Artifact})),
			1,
		)
		g.putIntoPlay(granter, 1)

		art := g.AddToHand(NewCard("art", Logos, Artifact, Common), 0)
		if _, err := g.PlayArtifact(0, handIdxByID(g, 0, art)); err != nil {
			t.Fatalf("Playartifact: %v", err)
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
			t.Fatalf("Playcreature: %v", err)
		}
		if !g.State.Cards[newbie].Exhausted {
			t.Error("creature should enter play exhausted without a grant")
		}
	})

	t.Run(
		"a pool-gated, house-filtered grant readies only off-house cards while rich",
		func(t *testing.T) {
			// Fandangle: while you have 4+ Æmber, your non-Untamed creatures enter ready.
			grant := EntersReadyGrant{Type: Creature, MinAember: 4, ExceptHouse: Untamed}
			g := NewGame("A", "B", 1)
			g.StartTurn(0)
			granter := g.AddToHand(
				NewCard("Fandangle", Untamed, Creature, Common,
					WithPower(3), WithFriendlyEntersPlayReady(grant)),
				0,
			)
			g.putIntoPlay(granter, 0)

			// Too little Æmber: the grant is dormant even for an off-house creature.
			g.SetAember(0, 3)
			if g.entersPlayReady(0, Creature, Logos) {
				t.Error("grant should be dormant below the Æmber threshold")
			}

			// Enough Æmber, off-house creature: readied.
			g.SetAember(0, 4)
			if !g.entersPlayReady(0, Creature, Logos) {
				t.Error("a non-Untamed creature should be readied at 4 Æmber")
			}

			// Enough Æmber, but an Untamed creature is excluded.
			if g.entersPlayReady(0, Creature, Untamed) {
				t.Error("an Untamed creature should be excluded by the house filter")
			}
		},
	)

	t.Run("renders the printed line for each grant", func(t *testing.T) {
		creat := &CardDefinition{
			Name:             "Duskwitch",
			EntersReadyGrant: EntersReadyGrant{Type: Creature},
		}
		if rules := cardRules(creat, false); len(rules) != 1 ||
			rules[0] != "Your creatures enter play ready." {
			t.Errorf("creature rules = %v", rules)
		}
		art := &CardDefinition{
			Name:             "The Curator",
			EntersReadyGrant: EntersReadyGrant{Type: Artifact},
		}
		if rules := cardRules(art, false); len(rules) != 1 ||
			rules[0] != "Friendly artifacts enter play ready." {
			t.Errorf("artifact rules = %v", rules)
		}
		fan := &CardDefinition{Name: "Fandangle", EntersReadyGrant: EntersReadyGrant{
			Type: Creature, MinAember: 4, ExceptHouse: Untamed,
		}}
		if rules := cardRules(fan, false); len(rules) != 1 ||
			rules[0] != "While you have 4 or more Æmber, your non-Untamed creatures enter play ready." {
			t.Errorf("gated rules = %v", rules)
		}
		filtered := &CardDefinition{Name: "Filtered", EntersReadyGrant: EntersReadyGrant{
			Type: Creature, ExceptHouse: Untamed,
		}}
		if rules := cardRules(filtered, false); len(rules) != 1 ||
			rules[0] != "Your non-Untamed creatures enter play ready." {
			t.Errorf("filtered rules = %v", rules)
		}
	})
}
