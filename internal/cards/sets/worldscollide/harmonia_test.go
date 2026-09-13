package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Harmonia
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Witch
//
//	Elusive.
//	After you play a Creature, if you are overwhelmed, gain 1 Æmber.
func TestHarmonia(t *testing.T) {
	t.Run("gains 1 Æmber after you play a creature while overwhelmed", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Harmonia),
				Hand:   ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(),
				ct.Creature(),
				ct.Creature(),
			)},
		})

		h.P1.Play(ally)

		h.P1.ExpectAmber(1)
	})

	t.Run("gains nothing when not overwhelmed", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Harmonia),
				Hand:   ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(ally)

		h.P1.ExpectAmber(0)
	})

	t.Run("gains 1 Æmber from its own entrance while overwhelmed", func(t *testing.T) {
		var harmonia ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&harmonia, Harmonia)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Creature(),
				ct.Creature(),
				ct.Creature(),
			)},
		})

		h.P1.Play(harmonia)

		h.P1.ExpectAmber(1)
	})
}
