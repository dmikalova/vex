package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dreadbone Decimus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Dinosaur • Assassin
//
//	Play/Fight: You may exalt Dreadbone Decimus -> destroy a creature with lower power than Dreadbone Decimus.
func TestDreadboneDecimus(t *testing.T) {
	t.Run("exalts itself and destroys a lower-power creature", func(t *testing.T) {
		var dreadbone, weak, strong ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&dreadbone, DreadboneDecimus)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&weak, ct.Creature(ct.Power(3))),
				ct.Bind(&strong, ct.Creature(ct.Power(6))),
			)},
		})

		h.P1.Play(dreadbone)
		h.P1.ClickCard(
			dreadbone,
		) // accept the optional exalt; only the lower-power creature is eligible

		h.Expect(dreadbone).AmberOn(1)
		h.Expect(weak).At(ct.Discard)
		h.Expect(strong).At(ct.PlayArea)
	})

	t.Run("declining the exalt destroys nothing", func(t *testing.T) {
		var dreadbone, weak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&dreadbone, DreadboneDecimus)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&weak, ct.Creature(ct.Power(3))))},
		})

		h.P1.Play(dreadbone)
		h.P1.ClickDone() // decline the exalt

		h.Expect(dreadbone).AmberOn(0)
		h.Expect(weak).At(ct.PlayArea)
	})
}
