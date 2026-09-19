package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Guji Dinosaur Hunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant • Hunter
//
//	Elusive.
//	Action: Choose a creature. If it is a Dinosaur creature or it has Æmber on it, deal 6 damage to it. Otherwise, deal 2 damage to it.
func TestGujiDinosaurHunter(t *testing.T) {
	t.Run("deals 2 to an ordinary creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(GujiDinosaurHunter)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		h.P1.UseAction(GujiDinosaurHunter)
		h.P1.ClickCard(foe)

		h.Expect(foe).Damage(2)
	})

	t.Run("deals 6 to a Dinosaur creature", func(t *testing.T) {
		var dino ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(GujiDinosaurHunter)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&dino,
					ct.Creature(ct.Power(7), ct.Traits(card.Traits.Dinosaur)))),
			},
		})

		h.P1.UseAction(GujiDinosaurHunter)
		h.P1.ClickCard(dino)

		h.Expect(dino).Damage(6)
	})

	t.Run("deals 6 to a creature with Æmber on it", func(t *testing.T) {
		var rich ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(GujiDinosaurHunter)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&rich, ct.Creature(ct.Power(7))))},
		})
		h.Game().State.Cards[rich.ID()].Amber = 1

		h.P1.UseAction(GujiDinosaurHunter)
		h.P1.ClickCard(rich)

		h.Expect(rich).Damage(6)
	})
}
