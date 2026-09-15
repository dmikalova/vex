package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Instrument of Silence
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains skirmish.
//	This creature gains, "Fight: Gain 1 Æmber."
func TestInstrumentOfSilence(t *testing.T) {
	t.Run("its host gains skirmish and gains 1 Æmber when it fights", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
						InstrumentOfSilence,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(host, foe)

		h.P1.ExpectAmber(1)
		h.Expect(host).Damage(0) // skirmish: no retaliation damage
	})
}
