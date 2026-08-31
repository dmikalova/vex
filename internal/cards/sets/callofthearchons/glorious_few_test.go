package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Glorious Few
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each creature your opponent controls in excess of you, gain 1 Æmber.
func TestGloriousFew(t *testing.T) {
	t.Run("gains 1 Æmber for each excess creature the opponent controls", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Hand:   ct.Cards(GloriousFew),
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
			)},
		})

		h.P1.Play(GloriousFew)

		h.P1.ExpectAmber(2) // opponent controls 3, you control 1 → 2 excess
	})
}
