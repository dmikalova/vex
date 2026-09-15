package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Carlo Phantom
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive, Skirmish.
//	After you play an artifact, steal 1 Æmber.
func TestCarloPhantom(t *testing.T) {
	t.Run("steals 1 Æmber each time its controller plays an artifact", func(t *testing.T) {
		var relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(CarloPhantom),
				Hand:   ct.Cards(ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Shadows)))),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(relic)

		h.P1.ExpectAmber(1) // stole on artifact play
		h.P2.ExpectAmber(2)
	})
}
