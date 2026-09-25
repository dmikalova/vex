package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Master of the Grey
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Human • Monk
//
//	Your opponent cannot resolve bonus icons on cards they play.
func TestMasterOfTheGrey(t *testing.T) {
	t.Run("the opponent's played card resolves none of its bonus icons", func(t *testing.T) {
		var bearer ct.Card
		h := ct.Play(t, ct.Setup{
			// P1 controls Master of the Grey; P2 is the barred opponent.
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(MasterOfTheGrey),
			},
			P2: ct.Side{
				House: card.House.Untamed,
				Hand: ct.Cards(
					ct.Bind(
						&bearer,
						ct.Creature(ct.OfHouse(card.House.Untamed), ct.AemberBonus(2)),
					),
				),
			},
		})

		// hand P2 the turn so they can play their creature.
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Untamed)
		h.P2.Play(bearer)

		h.P2.ExpectAmber(0) // the two Æmber bonus icons did not resolve
	})

	t.Run("your own played card still resolves its bonus icons", func(t *testing.T) {
		var bearer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(MasterOfTheGrey),
				Hand: ct.Cards(
					ct.Bind(
						&bearer,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.AemberBonus(2)),
					),
				),
			},
		})

		h.P1.Play(bearer)

		h.P1.ExpectAmber(2) // the controller is not barred
	})
}
