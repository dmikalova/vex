package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Begone!
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose one:
//	- Destroy each Dis creature
//	- Gain 1 Æmber.
func TestBegone(t *testing.T) {
	t.Run("first option destroys each Dis creature", func(t *testing.T) {
		var disGuy, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(Begone)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&disGuy, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Play(Begone)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("Dis creature")

		h.Expect(disGuy).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea)
	})
}
