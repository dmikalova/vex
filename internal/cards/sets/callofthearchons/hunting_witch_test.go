package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hunting Witch
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Witch
//
//	After you play a creature, gain 1 Æmber.
func TestHuntingWitch(t *testing.T) {
	t.Run("gains 1 Æmber after you play another creature", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(HuntingWitch),
				Hand:   ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(ally)

		h.P1.ExpectAmber(1)
	})
}
