package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// United Action
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For the remainder of the turn, you may play cards from any house for which you have a card in play. You cannot use any cards for the remainder of the turn.
func TestUnitedAction(t *testing.T) {
	t.Run("frees plays from controlled houses and bars using cards", func(t *testing.T) {
		var united ct.Card
		marsInPlay := ct.Creature(ct.OfHouse(card.House.Mars))
		marsInHand := ct.Creature(ct.OfHouse(card.House.Mars))
		disInHand := ct.Creature(ct.OfHouse(card.House.Dis))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(marsInPlay),
				Hand:   ct.Cards(ct.Bind(&united, UnitedAction), marsInHand, disInHand),
			},
		})

		h.P1.Play(united)

		// Mars is a house P1 has a card in play for, so its cards may be played;
		// Dis is not.
		h.P1.Play(marsInHand)
		h.P1.ExpectCannotPlay(disInHand)
		// And P1 cannot use any card this turn.
		h.P1.ExpectCannotUse(marsInPlay)
	})
}
