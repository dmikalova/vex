package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Dark Æmber Vault
//
//	House:  None
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Each friendly Mutant creature gains +2 power.
//	After you play a Mutant creature, draw a card.
func TestDarkAemberVault(t *testing.T) {
	t.Run("draws and buffs when you play a Mutant creature", func(t *testing.T) {
		var mutant, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(DarkAemberVault),
				Hand: ct.Cards(ct.Bind(&mutant, ct.Creature(
					ct.OfHouse(card.House.Sanctum),
					ct.Traits(card.Traits.Mutant),
					ct.Power(3),
				))),
				Deck: ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
		})

		h.P1.Play(mutant)

		h.Expect(top).At(ct.Hand)
		if got := mutant.Power(); got != 5 {
			t.Errorf("friendly Mutant power = %d, want 5 (3 base + 2)", got)
		}
	})

	t.Run(
		"draws when you play a Mutant creature with treachery into the enemy line",
		func(t *testing.T) {
			var mutant, top ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Sanctum,
					InPlay: ct.Cards(DarkAemberVault),
					Hand: ct.Cards(ct.Bind(&mutant, ct.Creature(
						ct.OfHouse(card.House.Sanctum),
						ct.Traits(card.Traits.Mutant),
						ct.Keywords(card.Keyword.Treachery),
						ct.Power(3),
					))),
					Deck: ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
				},
			})

			h.P1.Play(mutant)

			// Treachery hands the creature to your opponent, but you played it, so the
			// "after you play a Mutant creature" reaction still draws.
			if got := h.Game().Controller(mutant.ID()); got != 1 {
				t.Fatalf("controller = %d, want 1 (your opponent, via treachery)", got)
			}
			h.Expect(top).At(ct.Hand)
		},
	)

	t.Run("no draw when you play a non-Mutant creature", func(t *testing.T) {
		var plain, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(DarkAemberVault),
				Hand: ct.Cards(ct.Bind(&plain, ct.Creature(
					ct.OfHouse(card.House.Sanctum),
					ct.Traits(card.Traits.Knight),
					ct.Power(3),
				))),
				Deck: ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
		})

		h.P1.Play(plain)

		h.Expect(top).At(ct.Deck)
		if got := plain.Power(); got != 3 {
			t.Errorf("friendly non-Mutant power = %d, want 3 (no bonus)", got)
		}
	})

	t.Run("does not buff an enemy Mutant creature", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(DarkAemberVault),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(
					ct.OfHouse(card.House.Brobnar),
					ct.Traits(card.Traits.Mutant),
					ct.Power(4),
				))),
			},
		})
		_ = h

		if got := enemy.Power(); got != 4 {
			t.Errorf("enemy Mutant power = %d, want 4 (not buffed)", got)
		}
	})
}

// TestIsMutantCreature checks the deck-generation predicate Dark Æmber Vault uses
// to guarantee its deck a pool of Mutant creatures (card.PullsMatching).
func TestIsMutantCreature(t *testing.T) {
	cases := []struct {
		name string
		def  card.Definition
		want bool
	}{
		{
			"mutant creature",
			card.Definition{
				Type:   card.Type.Creature,
				Traits: []card.Trait{card.Traits.Mutant},
			},
			true,
		},
		{
			"non-mutant creature",
			card.Definition{
				Type:   card.Type.Creature,
				Traits: []card.Trait{card.Traits.Knight},
			},
			false,
		},
		{
			"mutant artifact",
			card.Definition{
				Type:   card.Type.Artifact,
				Traits: []card.Trait{card.Traits.Mutant},
			},
			false,
		},
	}
	for i := range cases {
		tc := &cases[i]
		if got := isMutantCreature(tc.def); got != tc.want {
			t.Errorf("isMutantCreature(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}
