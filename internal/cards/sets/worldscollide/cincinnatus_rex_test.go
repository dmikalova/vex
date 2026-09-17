package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cincinnatus Rex
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  4
//	Traits: Dinosaur • Soldier
//
//	If there are no enemy creatures in play, destroy Cincinnatus Rex.
//	Fight: You may exalt Cincinnatus Rex. Ready each other friendly card.
func TestCincinnatusRex(t *testing.T) {
	t.Run("fighting may exalt itself and ready each other friendly card", func(t *testing.T) {
		var rex, ally, foe, spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&rex, CincinnatusRex),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				// A second enemy keeps the self-destroy condition false.
				ct.Bind(&spare, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})
		ally.Exhaust()

		h.P1.Fight(rex, foe)
		h.P1.ClickCard(rex)

		h.Expect(foe).At(ct.Discard)
		h.Expect(rex).AmberOn(1)
		h.Expect(ally).Ready()
	})

	t.Run("dies when the last enemy creature leaves play", func(t *testing.T) {
		var rex, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&rex, CincinnatusRex)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		h.P1.Fight(rex, foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(rex).At(ct.Discard)
	})
}
