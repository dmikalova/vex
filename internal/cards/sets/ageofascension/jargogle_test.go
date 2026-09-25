package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Jargogle
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Beast • Mutant
//
//	Elusive.
//	Play: Put a card from your hand facedown under Jargogle.
//	Destroyed: If it is your turn, play the card under Jargogle. Otherwise, archive the card under Jargogle.
func TestJargogle(t *testing.T) {
	t.Run("play puts a card from hand facedown under it", func(t *testing.T) {
		var jargogle, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					ct.Bind(&jargogle, Jargogle),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(jargogle)

		h.Expect(jargogle).At(ct.PlayArea)
		h.Expect(buried).At(ct.Under)
	})

	t.Run("destroyed on your turn plays the card under it", func(t *testing.T) {
		var jargogle, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					ct.Bind(&jargogle, Jargogle),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})
		h.P1.Play(jargogle)

		h.Game().DestroyEach(0, []engine.LocalID{jargogle.ID()})

		h.Expect(jargogle).At(ct.Discard)
		h.Expect(buried).At(ct.PlayArea)
	})

	t.Run("destroyed on the opponent's turn archives the card under it", func(t *testing.T) {
		var jargogle, buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					ct.Bind(&jargogle, Jargogle),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})
		h.P1.Play(jargogle)
		h.Game().State.ActivePlayer = 1

		h.Game().DestroyEach(1, []engine.LocalID{jargogle.ID()})

		h.Expect(jargogle).At(ct.Discard)
		h.Expect(buried).At(ct.Archives)
	})
}
