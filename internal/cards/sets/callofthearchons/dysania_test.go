package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dysania
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	Play: For each card in your opponent's archives, gain 1 Æmber, and your opponent discards each of their archived cards.
func TestDysania(t *testing.T) {
	t.Run(
		"gains 1 Æmber per opponent-archived card, then discards those archives",
		func(t *testing.T) {
			var a, b ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(Dysania)},
				P2: ct.Side{
					Archives: ct.Cards(
						ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
						ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
					),
				},
			})

			h.P1.Play(Dysania)

			h.P1.ExpectAmber(2) // 1 per each of 2 opponent archived cards
			h.Expect(a).At(ct.Discard)
			h.Expect(b).At(ct.Discard)
		},
	)
}
