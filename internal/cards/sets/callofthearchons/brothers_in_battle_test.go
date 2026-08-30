package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Brothers in Battle
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose a house - for the remainder of the turn, each friendly creature of the chosen house may fight.
func TestBrothersInBattle(t *testing.T) {
	t.Run("lets friendly creatures of the chosen house fight out of the active house", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(BrothersInBattle),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3)))),
			},
		})

		h.P1.Play(BrothersInBattle)
		h.P1.ExpectPrompt("Choose a house").Source("Brothers in Battle")
		h.P1.ClickOption("Untamed")

		// The Untamed ally may now fight though Brobnar is the active house.
		h.P1.Fight(ally, foe)

		h.Expect(foe).At(ct.Discard)         // 5 power destroys the 3-power foe
		h.Expect(ally).Damage(3).Exhausted() // fought, so it took return damage
	})
}
