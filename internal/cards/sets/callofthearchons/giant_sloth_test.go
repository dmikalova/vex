package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Giant Sloth
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast
//
//	You cannot use this card unless you have discarded an Untamed card from your hand this turn.
//	Action: Gain 3 Æmber.
func TestGiantSloth(t *testing.T) {
	setup := func(t *testing.T, sloth *ct.Card) *ct.Harness {
		t.Helper()
		return ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(NaturesCall),
				InPlay: ct.Cards(
					ct.Bind(sloth, GiantSloth),
				),
			},
		})
	}

	t.Run("cannot be used before an Untamed card is discarded", func(t *testing.T) {
		var sloth ct.Card
		h := setup(t, &sloth)

		h.P1.ExpectCannotUse(sloth)
		h.P1.ExpectAmber(0)
	})

	t.Run("can be used once an Untamed card is discarded from hand", func(t *testing.T) {
		var sloth ct.Card
		h := setup(t, &sloth)

		h.P1.Discard(NaturesCall)
		h.P1.UseAction(sloth)

		h.P1.ExpectAmber(3)
	})
}
