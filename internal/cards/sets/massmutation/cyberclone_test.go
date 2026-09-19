package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cyber-Clone
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Play: Purge another creature. Cyber-Clone has power equal to the same creature's printed power and gains its printed armor, keywords, and traits.
func TestCyberClone(t *testing.T) {
	t.Run("purges another creature and copies its printed power and armor", func(t *testing.T) {
		var model ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(CyberClone),
				InPlay: ct.Cards(
					ct.Bind(&model, ct.Creature(ct.Power(6), ct.Armor(3))),
				),
			},
		})

		h.P1.Play(CyberClone)

		h.Expect(model).At(ct.Purge)
		// Cyber-Clone's own power (1) and armor (0) give way to the purged
		// creature's printed 6 power and 3 armor.
		h.Expect(CyberClone).At(ct.PlayArea).Power(6).Armor(3)
	})

	t.Run("copies nothing when there is no other creature to purge", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(CyberClone),
			},
		})

		h.P1.Play(CyberClone)

		h.Expect(CyberClone).At(ct.PlayArea).Power(1).Armor(0)
	})
}
