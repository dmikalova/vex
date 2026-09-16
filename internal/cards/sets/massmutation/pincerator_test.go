package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pincerator
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	At the end of each player's turn, deal 1 damage to each flank creature.
func TestPincerator(t *testing.T) {
	t.Run("damages every flank creature on both sides but not the middle", func(t *testing.T) {
		var p1Left, p1Mid, p1Right ct.Card
		var p2Left, p2Mid, p2Right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					Pincerator,
					ct.Bind(&p1Left, ct.Creature(ct.Power(3))),
					ct.Bind(&p1Mid, ct.Creature(ct.Power(3))),
					ct.Bind(&p1Right, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&p2Left, ct.Creature(ct.Power(3))),
					ct.Bind(&p2Mid, ct.Creature(ct.Power(3))),
					ct.Bind(&p2Right, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.EndTurn() // Pincerator fires at the end of P1's turn.

		h.Expect(p1Left).Damage(1)
		h.Expect(p1Mid).Damage(0)
		h.Expect(p1Right).Damage(1)
		h.Expect(p2Left).Damage(1)
		h.Expect(p2Mid).Damage(0)
		h.Expect(p2Right).Damage(1)
	})
}
