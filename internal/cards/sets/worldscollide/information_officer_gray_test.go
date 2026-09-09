package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Information Officer Gray
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Play/Fight/Reap: You may reveal a non-Star Alliance card from your hand and archive it.
func TestInformationOfficerGray(t *testing.T) {
	t.Run("reaping reveals and archives a non-Star Alliance card", func(t *testing.T) {
		var gray, outsider ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&gray, InformationOfficerGray)),
				Hand: ct.Cards(
					ct.Bind(&outsider, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Reap(gray)
		h.P1.ClickCard(outsider)

		h.Expect(outsider).At(ct.Archives)
	})

	t.Run("a Star Alliance card in hand cannot be revealed", func(t *testing.T) {
		var gray, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&gray, InformationOfficerGray)),
				Hand: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				),
			},
		})

		h.P1.Reap(gray)

		h.Expect(ally).At(ct.Hand)
	})
}
