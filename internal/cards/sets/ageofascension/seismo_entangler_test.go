package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Seismo-entangler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Choose a house. Your opponent cannot use creatures of the chosen house to reap during their next turn.
func TestSeismoEntangler(t *testing.T) {
	t.Run("bars the chosen house from reaping on the opponent's next turn", func(t *testing.T) {
		var seismo, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&seismo, SeismoEntangler)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(seismo)
		h.P1.ClickOption("Brobnar")
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.ExpectCannotUseTo(foe, engine.ReapUse)
	})

	t.Run("lifts once that turn is over", func(t *testing.T) {
		var seismo, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&seismo, SeismoEntangler)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(seismo)
		h.P1.ClickOption("Brobnar")
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.Reap(foe)
	})

	t.Run("leaves other houses free to reap", func(t *testing.T) {
		var seismo, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&seismo, SeismoEntangler)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(seismo)
		h.P1.ClickOption("Brobnar")
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Untamed)

		h.P2.Reap(foe)
	})
}
