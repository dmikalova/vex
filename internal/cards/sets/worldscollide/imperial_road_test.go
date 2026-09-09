package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Imperial Road
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Versatile.
//	Action: Play a Saurian creature -> stun it.
func TestImperialRoad(t *testing.T) {
	t.Run("plays a Saurian creature from hand and stuns it", func(t *testing.T) {
		var dino ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ImperialRoad),
				Hand: ct.Cards(
					ct.Bind(&dino, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
				),
			},
		})

		h.P1.UseAction(ImperialRoad)

		h.Expect(dino).At(ct.PlayArea).Stunned(true)
	})
}
