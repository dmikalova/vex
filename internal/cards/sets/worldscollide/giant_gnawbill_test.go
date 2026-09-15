package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Giant Gnawbill
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Beast
//
//	After a player chooses an active house, that player destroys an artifact of that house.
func TestGiantGnawbill(t *testing.T) {
	t.Run("the chooser destroys an artifact of the chosen house", func(t *testing.T) {
		var relic, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(GiantGnawbill)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Untamed))),
				ct.Bind(&other, ct.Artifact(ct.OfHouse(card.House.Logos))),
			)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Untamed)

		h.Expect(relic).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea)
	})

	t.Run("does nothing when no artifact of the chosen house is in play", func(t *testing.T) {
		var relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(GiantGnawbill)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Brobnar)),
				ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Untamed))),
			)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.Expect(relic).At(ct.PlayArea)
	})
}
