package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Stealthster
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot
//
//	Elusive.
//	Stealthster may be played as an upgrade instead of a creature, with the text: "This creature gains elusive."
func TestStealthster(t *testing.T) {
	t.Run("plays as a creature", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(Stealthster),
			},
		})

		h.P1.Play(Stealthster)

		h.Expect(Stealthster).At(ct.PlayArea)
	})
}
