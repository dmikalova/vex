package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Gargantodon
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  16
//	Traits: Beast
//
//	Gargantodon deals 4 damage when fighting.
//	Each Æmber that would be stolen is captured by a creature controlled by the active player instead.
//	Gargantodon enters play stunned.
func TestGargantodon(t *testing.T) {
	t.Run("enters play stunned", func(t *testing.T) {
		var garg ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&garg, Gargantodon)),
			},
		})

		h.P1.Play(Gargantodon)

		h.Expect(garg).Stunned(true)
	})

	t.Run("deals exactly 4 fight damage despite its 16 power", func(t *testing.T) {
		var garg, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&garg, Gargantodon)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.Power(6))),
			)},
		})

		h.P1.Fight(garg, foe)

		h.Expect(foe).At(ct.PlayArea).Damage(4)
		h.Expect(garg).At(ct.PlayArea).Damage(6)
	})

	t.Run("a steal is redirected into a capture onto a friendly creature", func(t *testing.T) {
		var garg ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&garg, Gargantodon)),
			},
			P2: ct.Side{Amber: 5},
		})
		g := h.Game()

		engine.StealAember{Amount: 2}.Resolve(&engine.EffectContext{
			Resolver:   g,
			Controller: 0,
		})

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
		h.Expect(garg).AmberOn(2)
	})
}
