package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Grey Rider
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Monk
//
//	Deploy.
//	Play/Fight/Reap: You may ready and fight with a neighboring creature.
func TestTheGreyRider(t *testing.T) {
	t.Run("gains aember when reaping and may decline the bonus effect", func(t *testing.T) {
		var rider, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&rider, TheGreyRider),
					ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Reap(rider)
		h.P1.ClickDone()

		h.P1.ExpectAmber(1)
	})
}
