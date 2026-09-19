package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gizelhart's Wrath
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each Mutant creature.
func TestGizelhartsWrath(t *testing.T) {
	t.Run("destroys each Mutant creature", func(t *testing.T) {
		var mutant, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(GizelhartsWrath),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&mutant, ct.Creature(ct.Traits(card.Traits.Mutant))),
				ct.Bind(&other, ct.Creature(ct.Traits(card.Traits.Beast))),
			)},
		})

		h.P1.Play(GizelhartsWrath)

		h.Expect(mutant).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea)
	})
}
