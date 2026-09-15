package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Master the Theory
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: If there are no friendly creatures in play, for each enemy creature in play, you may archive a card from your hand.
func TestMasterTheTheory(t *testing.T) {
	t.Run("does nothing while you control a creature", func(t *testing.T) {
		var spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Creature(ct.Power(3))),
				Hand: ct.Cards(
					MasterTheTheory,
					ct.Bind(&spare, ct.Creature(ct.Power(2))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.Power(4))),
			},
		})

		h.P1.Play(MasterTheTheory)

		h.Expect(spare).At(ct.Hand)
	})
}
