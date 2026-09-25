package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Senator Quintina
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	After a creature reaps, exalt it.
func TestSenatorQuintina(t *testing.T) {
	t.Run("exalts a friendly creature after it reaps", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					SenatorQuintina,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(ally)

		h.Expect(ally).AmberOn(1)
	})

	t.Run("exalts an enemy creature after it reaps", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(SenatorQuintina),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3))))},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Reap(enemy)

		h.Expect(enemy).AmberOn(1)
	})
}
