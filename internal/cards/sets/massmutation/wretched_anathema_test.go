package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wretched Anathema
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  10
//	Traits: Demon
//
//	While there are no other friendly creatures in play, Wretched Anathema gains, "Action: Gain 4 Æmber."
//	Play/Reap: Destroy 2 other creatures.
func TestWretchedAnathema(t *testing.T) {
	t.Run("play destroys two other creatures", func(t *testing.T) {
		var wa, friendly, spared, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&wa, WretchedAnathema), card.GiganticArt(WretchedAnathema)),
				InPlay: ct.Cards(
					ct.Bind(&friendly, ct.Creature(ct.Power(3))),
					ct.Bind(&spared, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Play(wa)
		h.P1.ClickCard(friendly)
		h.P1.ClickCard(enemy)

		h.Expect(friendly).At(ct.Discard)
		h.Expect(enemy).At(ct.Discard)
		h.Expect(spared).Power(3)
	})

	t.Run("alone it gains an Action to gain 4 Æmber", func(t *testing.T) {
		var wa ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&wa, WretchedAnathema)),
			},
		})

		h.P1.UseAction(wa)

		h.P1.ExpectAmber(4)
	})
}
