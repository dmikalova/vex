package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Praefectus Ludo
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	Each other friendly creature gains, "Destroyed: Move each Æmber on this creature to the common supply."
func TestPraefectusLudo(t *testing.T) {
	t.Run("a destroyed friendly creature moves its Æmber to the common supply", func(t *testing.T) {
		var ally, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					PraefectusLudo,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})
		h.Game().State.Cards[ally.ID()].Amber = 2

		h.P1.Fight(ally, enemy)

		h.Expect(ally).At(ct.Discard)
		h.P2.ExpectAmber(0) // Æmber went to the common supply, not the opponent
	})
}
