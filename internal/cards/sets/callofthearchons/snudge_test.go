package callofthearchons_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	cota "github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
)

// Snudge
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Fight/Reap: Put an artifact or flank creature into its owner's hand.
func TestSnudge(t *testing.T) {
	t.Run("returns an enemy artifact", func(t *testing.T) {
		var snudge, relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&snudge, cota.Snudge)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&relic, ct.Artifact()))},
		})

		h.P1.Reap(snudge)
		h.P1.ClickCard(relic)
		h.Expect(relic).At(ct.Hand)
	})

	t.Run("reaches a flank creature but not one in the middle", func(t *testing.T) {
		var snudge, flanker, middle ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&snudge, cota.Snudge)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&flanker, ct.Creature(ct.Power(3))),
					ct.Bind(&middle, ct.Creature(ct.Power(3))),
					ct.Creature(ct.Power(3)),
				),
			},
		})

		// Snudge itself is on a flank of its own line, so the choice is real.
		h.P1.Reap(snudge)
		h.P1.ClickCard(flanker)
		h.Expect(flanker).At(ct.Hand)
		h.Expect(middle).At(ct.PlayArea)
	})
}
