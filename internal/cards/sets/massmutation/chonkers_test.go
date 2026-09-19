package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chonkers
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Skirmish.
//	After a creature is destroyed in a fight with Chonkers, give Chonkers +1 power counters equal to the number of +1 power counters on Chonkers.
//	Play: Give Chonkers a +1 power counter.
func TestChonkers(t *testing.T) {
	t.Run("Play gives Chonkers a +1 power counter", func(t *testing.T) {
		var chonkers ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&chonkers, Chonkers)),
			},
		})
		h.P1.Play(chonkers)
		h.Expect(chonkers).Power(2)
	})

	t.Run("doubles its power counters when an enemy dies fighting it", func(t *testing.T) {
		var chonkers, prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&chonkers, Chonkers)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&prey, ct.Creature(ct.Power(1))))},
		})
		h.Game().AddPowerCounter(chonkers.ID(), 2) // power 3

		h.P1.Fight(chonkers, prey)

		h.Expect(prey).At(ct.Discard)
		h.Expect(chonkers).Power(5) // counters doubled 2 -> 4, over base 1
	})
}
