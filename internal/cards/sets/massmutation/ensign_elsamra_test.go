package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ensign El-Samra
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Mutant
//
//	Action: Reveal a card from your hand. Resolve that card's bonus icons.
//	Enhance Draw Draw Draw.
func TestEnsignElSamra(t *testing.T) {
	t.Run("action reveals a hand card and resolves its bonus icons", func(t *testing.T) {
		var elsamra, revealed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&elsamra, EnsignElSamra)),
				Hand:   ct.Cards(ct.Bind(&revealed, ct.Creature(ct.Bonus(card.Bonus.Aember)))),
			},
		})

		before := h.P1.Amber()
		elsamra.Ready()
		h.P1.UseAction(elsamra)

		// The revealed card's Æmber bonus resolves for a gain of 1.
		if got := h.P1.Amber() - before; got != 1 {
			t.Fatalf("aember gained = %d, want 1", got)
		}
		// The revealed card stays in hand.
		h.Expect(revealed).At(ct.Hand)
	})
}
