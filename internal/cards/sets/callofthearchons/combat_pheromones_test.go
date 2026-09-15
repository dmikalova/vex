package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Combat Pheromones
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Combat Pheromones. Use 2 other Mars cards, one at a time.
func TestCombatPheromones(t *testing.T) {
	t.Run("destroys itself and uses two other ready Mars creatures", func(t *testing.T) {
		var pheromones, first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&pheromones, CombatPheromones),
					ct.Bind(&first, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&second, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(pheromones)
		h.P1.ClickCard(second)

		h.Expect(pheromones).At(ct.Discard)
		h.Expect(first).Exhausted()
		h.Expect(second).Exhausted()
		h.P1.ExpectAmber(2)
	})

	t.Run("does not offer non-Mars cards", func(t *testing.T) {
		var pheromones, mars, logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&pheromones, CombatPheromones),
					ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(pheromones)

		h.Expect(pheromones).At(ct.Discard)
		h.Expect(mars).Exhausted()
		h.Expect(logos).Ready()
		h.P1.ExpectAmber(1)
	})
}
