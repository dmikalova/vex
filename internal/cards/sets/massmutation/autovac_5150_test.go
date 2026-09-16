package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Auto-Vac 5150
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Discard a card from your archives -> keys cost +3 Æmber during your opponent's next turn. Otherwise, archive a card from your hand.
func TestAutoVac5150(t *testing.T) {
	t.Run("discards from archives to tax the opponent's next turn", func(t *testing.T) {
		var vac, stored ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:    card.House.Logos,
				InPlay:   ct.Cards(ct.Bind(&vac, AutoVac5150)),
				Archives: ct.Cards(ct.Bind(&stored, ct.Creature())),
			},
			P2: ct.Side{House: card.House.Logos, Amber: 8},
		})

		h.P1.UseAction(AutoVac5150)
		h.P1.ClickCard(stored) // discard from archives rather than decline
		h.Expect(stored).At(ct.Discard)

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)

		// 8 Æmber does not cover a key at 6 + 3.
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(8)
	})

	t.Run("declining archives a card from hand instead", func(t *testing.T) {
		var vac, stored, held ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:    card.House.Logos,
				InPlay:   ct.Cards(ct.Bind(&vac, AutoVac5150)),
				Archives: ct.Cards(ct.Bind(&stored, ct.Creature())),
				Hand:     ct.Cards(ct.Bind(&held, ct.Creature())),
			},
			P2: ct.Side{House: card.House.Logos, Amber: 6},
		})

		h.P1.UseAction(AutoVac5150)
		h.P1.ClickDone() // decline the archives discard

		h.Expect(stored).At(ct.Archives) // untouched
		h.Expect(held).At(ct.Archives)   // the hand card is archived

		// No tax applied: a key at 6 is affordable.
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)
		h.P2.ExpectKeys(1)
	})
}
