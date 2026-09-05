package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lamindra
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Deploy, Elusive.
//	Each neighboring creature gains elusive.
func TestLamindra(t *testing.T) {
	t.Run("grants elusive to a neighbor so its first attack deals no damage", func(t *testing.T) {
		var neighbor, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(
				Lamindra,
				ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
			)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars)
		h.P2.Fight(foe, neighbor)
		h.Expect(neighbor).At(ct.PlayArea).Damage(0) // elusive granted by Lamindra
	})
}
