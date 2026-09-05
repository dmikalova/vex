package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Plague Rat
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Rat
//
//	Elusive.
//	Play: For each Rat trait creature in play, deal 1 damage to each non-Rat trait creature.
func TestPlagueRat(t *testing.T) {
	t.Run("deals 1 damage to each non-rat creature for each rat in play", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Hand:   ct.Cards(PlagueRat),
				InPlay: ct.Cards(ct.Creature(ct.Power(5), ct.Traits(card.Traits.Rat))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(10))))},
		})

		h.P1.Play(PlagueRat)

		// Two Rats in play (the ally Rat and Plague Rat itself), so 2 damage.
		h.Expect(foe).Damage(2)
	})
}
