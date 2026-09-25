package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Vault's Blessing
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: For each friendly Mutant creature, each player gains 1 Æmber.
func TestVaultsBlessing(t *testing.T) {
	t.Run("each player gains 1 Æmber per Mutant creature they control", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(VaultsBlessing),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Mutant)),
					ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Mutant)),
					ct.Creature(ct.OfHouse(card.House.Untamed)), // non-Mutant, not counted
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Traits(card.Traits.Mutant)),
				),
			},
		})

		h.P1.Play(VaultsBlessing)

		// P1 controls 2 Mutants -> gains 2; plus the card's own 1 Æmber bonus = 3.
		h.P1.ExpectAmber(3)
		// P2 controls 1 Mutant -> gains 1.
		h.P2.ExpectAmber(1)
	})

	t.Run("a player with no Mutant creatures gains nothing", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(VaultsBlessing),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Traits(card.Traits.Mutant)),
				),
			},
		})

		h.P1.Play(VaultsBlessing)

		// P1 controls 0 Mutants -> only the card's 1 Æmber bonus.
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
