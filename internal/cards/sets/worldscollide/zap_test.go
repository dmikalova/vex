package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Zap
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For each house represented among creatures in play, deal 1 damage to a creature.
func TestZap(t *testing.T) {
	t.Run("deals 1 damage per house represented among creatures in play", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(Zap),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.StarAlliance)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
			)},
		})

		h.P1.Play(Zap)
		h.P1.ClickCard(foe)

		// Star Alliance, Logos, Brobnar = 3 houses in play.
		h.Expect(foe).At(ct.PlayArea).Damage(3)
	})
}
