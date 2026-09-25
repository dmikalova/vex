package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Glyxl Proliferator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Reap: If Glyxl Proliferator is on a flank, archive a Mars card from your discard pile.
func TestGlyxlProliferator(t *testing.T) {
	t.Run("reaps to archive a Mars card from discard while on a flank", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(GlyxlProliferator),
				Discard: ct.Cards(
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Reap(GlyxlProliferator)

		h.Expect(buried).At(ct.Archives)
	})

	t.Run("archives nothing while not on a flank", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
					GlyxlProliferator,
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
				),
				Discard: ct.Cards(
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.Reap(GlyxlProliferator)

		h.Expect(buried).At(ct.Discard)
	})
}
