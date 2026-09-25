package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Forum of Giants
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	At the start of your turn, choose the most powerful creature. Its controller gains 1 Æmber.
func TestForumOfGiants(t *testing.T) {
	t.Run("the controller of the single most powerful creature gains", func(t *testing.T) {
		var mine, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ForumOfGiants,
					ct.Bind(&mine, ct.Creature(ct.Power(7))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&theirs, ct.Creature(ct.Power(4))))},
		})

		h.P1.EndTurn()
		h.P2.EndTurn() // back to P1: Forum fires at the start of the turn

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(0)
	})

	t.Run("a tie prompts the active player to pick the controller", func(t *testing.T) {
		var mine, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ForumOfGiants,
					ct.Bind(&mine, ct.Creature(ct.Power(5))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&theirs, ct.Creature(ct.Power(5))))},
		})

		h.P1.EndTurn()
		h.P2.EndTurn()         // back to P1: Forum fires and the tie prompts
		h.P1.ClickCard(theirs) // active player awards it to the opponent's creature

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(1)
	})
}
