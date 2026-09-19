package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fangs of Gizelhart
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Purge the most powerful creature.
func TestFangsOfGizelhart(t *testing.T) {
	t.Run("purges the most powerful creature", func(t *testing.T) {
		var strong, weak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(FangsOfGizelhart),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&strong, ct.Creature(ct.Power(6))),
				ct.Bind(&weak, ct.Creature(ct.Power(2))),
			)},
		})

		h.P1.Play(FangsOfGizelhart)

		h.Expect(strong).At(ct.Purge)
		h.Expect(weak).At(ct.PlayArea)
	})
}
