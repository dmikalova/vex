package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Experimental Therapy
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains versatile.
//	Play: Stun and exhaust this creature.
func TestExperimentalTherapy(t *testing.T) {
	t.Run("stuns and exhausts its host when played", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(ExperimentalTherapy),
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(ExperimentalTherapy)

		h.Expect(host).Stunned(true).Exhausted()
	})
}
