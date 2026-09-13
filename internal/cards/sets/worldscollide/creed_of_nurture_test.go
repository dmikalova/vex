package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Creed of Nurture
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Creed of Nurture. Reveal a Creature from your hand and choose a Creature in play - for the remainder of the turn, the chosen Creature gains the text box of the revealed Creature.
func TestCreedOfNurture(t *testing.T) {
	t.Run("lends a revealed creature's text box to a chosen creature", func(t *testing.T) {
		var recipient, loaner, defender ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					CreedOfNurture,
					ct.Bind(&recipient, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(6))),
				),
				Hand: ct.Cards(
					ct.Bind(&loaner, ct.Creature(ct.Keywords(card.Keyword.Skirmish))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&defender, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.UseAction(CreedOfNurture)
		h.P1.ClickCard(recipient) // the revealed loaner is the only hand creature

		// The action destroys Creed of Nurture.
		h.Expect(CreedOfNurture).At(ct.Discard)

		// The recipient gained the loaner's skirmish for the turn: its 6 fight
		// damage destroys the 5-power defender and skirmish spares its return
		// damage, so it survives undamaged.
		h.P1.Fight(recipient, defender)
		h.Expect(defender).At(ct.Discard)
		h.Expect(recipient).At(ct.PlayArea).Damage(0)
	})
}
