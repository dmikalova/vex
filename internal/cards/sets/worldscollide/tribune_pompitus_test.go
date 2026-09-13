package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tribune Pompitus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Dinosaur • Politician
//
//	Each friendly Creature gains +2 power for each Æmber on it.
//	Before Fight: You may exalt Tribune Pompitus.
func TestTribunePompitus(t *testing.T) {
	t.Run("each friendly creature is buffed by the Æmber on it", func(t *testing.T) {
		var tribune, friend, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&tribune, TribunePompitus),
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})

		// With no Æmber anywhere the bonus scales to zero for both.
		h.Expect(tribune).Power(4)
		h.Expect(friend).Power(3)

		// Before Fight exalts Tribune, placing one Æmber on it; only Tribune grows.
		h.P1.Fight(tribune, foe)
		h.P1.ClickCard(tribune)

		h.Expect(tribune).AmberOn(1)
		h.Expect(tribune).Power(6)
		h.Expect(friend).Power(3)
	})

	t.Run("declining the exalt leaves power at base", func(t *testing.T) {
		var tribune, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&tribune, TribunePompitus)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})

		h.P1.Fight(tribune, foe)
		h.P1.ClickDone()

		h.Expect(tribune).AmberOn(0)
		h.Expect(tribune).Power(4)
	})
}
