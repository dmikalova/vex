package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nepeta Gigantica
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Choose one:
//	- Stun a creature with power 5 or higher
//	- Stun a Giant creature.
func TestNepetaGigantica(t *testing.T) {
	t.Run("stuns a creature with power 5 or higher", func(t *testing.T) {
		var nepeta, bruiser ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&nepeta, NepetaGigantica)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&bruiser, ct.Creature(ct.Power(6))))},
		})

		h.P1.UseAction(NepetaGigantica)
		h.P1.ClickOption("stun a creature with power 5 or higher")

		h.Expect(bruiser).Stunned(true)
	})

	t.Run("stuns a Giant creature", func(t *testing.T) {
		var nepeta, giant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&nepeta, NepetaGigantica)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&giant, ct.Creature(ct.Power(2), ct.Traits(card.Traits.Giant))),
			)},
		})

		h.P1.UseAction(NepetaGigantica)
		h.P1.ClickOption("stun a Giant creature")

		h.Expect(giant).Stunned(true)
	})
}
