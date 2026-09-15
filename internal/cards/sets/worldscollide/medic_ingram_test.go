package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Medic Ingram
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: You may choose a creature - heal 3 damage from it, and ward it.
func TestMedicIngram(t *testing.T) {
	t.Run("heals a creature 3 and wards it when played", func(t *testing.T) {
		var wounded ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(MedicIngram),
				InPlay: ct.Cards(
					ct.Bind(
						&wounded,
						ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(5)),
					),
				),
			},
		})
		wounded.Damaged(3)

		h.P1.Play(MedicIngram)
		h.P1.ClickCard(wounded)

		h.Expect(wounded).Damage(0)
		if !h.Game().Warded(wounded.ID()) {
			t.Error("the healed creature should be warded")
		}
	})
}
