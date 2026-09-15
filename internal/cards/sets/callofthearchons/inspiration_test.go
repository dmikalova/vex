package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Inspiration
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ready and use a friendly creature.
func TestInspiration(t *testing.T) {
	t.Run("readies and uses a friendly creature to reap", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Hand:   ct.Cards(Inspiration),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
		})

		ally.Exhaust()
		h.P1.Play(Inspiration)

		h.P1.ExpectAmber(1) // the readied creature reaped
	})
}
