package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mindfire
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent discards a random card from their hand. For each bonus icon on the discarded card, steal 1 Æmber.
func TestMindfire(t *testing.T) {
	t.Run("steals 1 Æmber per bonus icon on the discarded card", func(t *testing.T) {
		var pitched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Mindfire),
			},
			P2: ct.Side{
				Amber: 5,
				Hand:  ct.Cards(ct.Bind(&pitched, ct.Tactic(ct.AemberBonus(2)))),
			},
		})

		h.P1.Play(Mindfire)

		h.Expect(pitched).At(ct.Discard)
		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(3)
	})

	t.Run("a discarded card with no bonus icons steals nothing", func(t *testing.T) {
		var pitched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Mindfire),
			},
			P2: ct.Side{
				Amber: 5,
				Hand:  ct.Cards(ct.Bind(&pitched, ct.Tactic())),
			},
		})

		h.P1.Play(Mindfire)

		h.Expect(pitched).At(ct.Discard)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(5)
	})

	t.Run("an empty opponent hand steals nothing", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Mindfire),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(Mindfire)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(5)
	})
}
