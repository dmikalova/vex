package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Painmail
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "After a player chooses Dis as their active house, archive Painmail, and destroy this creature."
func TestPainmail(t *testing.T) {
	t.Run("choosing Dis archives Painmail and destroys the host", func(t *testing.T) {
		var pain, host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(4))),
				),
				Hand: ct.Cards(ct.Bind(&pain, Painmail)),
			},
		})

		h.P1.Play(pain) // the lone host auto-attaches
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)

		h.Expect(pain).At(ct.Archives)
		h.Expect(host).At(ct.Discard)
	})

	t.Run("choosing another house leaves Painmail attached", func(t *testing.T) {
		var pain, host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(4))),
				),
				Hand: ct.Cards(ct.Bind(&pain, Painmail)),
			},
		})

		h.P1.Play(pain)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)

		h.Expect(pain).At(ct.Attached)
		h.Expect(host).At(ct.PlayArea)
	})
}
