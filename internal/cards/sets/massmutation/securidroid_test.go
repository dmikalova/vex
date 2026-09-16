package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Securi-Droid
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Robot
//
//	Taunt.
//	Securi-Droid may be played as an upgrade instead of a creature, with the text: "This creature gains taunt."
func TestSecuriDroid(t *testing.T) {
	t.Run("plays as a creature", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.StarAlliance, Hand: ct.Cards(SecuriDroid)},
		})

		h.P1.Play(SecuriDroid)

		h.Expect(SecuriDroid).At(ct.PlayArea)
	})
}
