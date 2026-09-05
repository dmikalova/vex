package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Masterplan
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Versatile.
//	Play: Put a card from your hand facedown under Masterplan.
//	Action: Play the card under Masterplan, and destroy Masterplan.
func TestMasterplan(t *testing.T) {
	t.Run("play puts a card from hand facedown under it", func(t *testing.T) {
		var masterplan, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand: ct.Cards(
					ct.Bind(&masterplan, Masterplan),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
				),
			},
		})

		h.P1.Play(masterplan)

		h.Expect(masterplan).At(ct.PlayArea)
		h.Expect(buried).At(ct.Under)
	})

	t.Run("action plays the buried card facedown and destroys itself", func(t *testing.T) {
		var masterplan, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand: ct.Cards(
					ct.Bind(&masterplan, Masterplan),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
				),
			},
			P2: ct.Side{House: card.House.Brobnar},
		})
		h.P1.Play(masterplan)
		h.P1.EndTurn()
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Shadows)

		h.P1.UseAction(masterplan)

		h.Expect(masterplan).At(ct.Discard)
		h.Expect(buried).At(ct.PlayArea)
	})

	t.Run("play does nothing with an empty hand", func(t *testing.T) {
		var masterplan ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&masterplan, Masterplan)),
			},
		})

		h.P1.Play(masterplan)

		h.Expect(masterplan).At(ct.PlayArea)
	})

	t.Run("action does nothing with no card buried", func(t *testing.T) {
		var masterplan ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&masterplan, Masterplan)),
			},
		})

		h.P1.UseAction(masterplan)

		h.Expect(masterplan).At(ct.Discard)
	})
}
