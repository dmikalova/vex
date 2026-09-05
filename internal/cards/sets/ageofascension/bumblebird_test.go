package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bumblebird
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Beast • Insect
//
//	Alpha.
//	Play: Give each other friendly Untamed creature 2 +1 power counters.
func TestBumblebird(t *testing.T) {
	t.Run("adds 2 power counters to each other friendly creature", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Bumblebird),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				),
			},
		})

		h.P1.Play(Bumblebird)

		h.Expect(ally).Power(6)
	})
}
