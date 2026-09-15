package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hystricog
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant
//
//	Action: Destroy a damaged creature.
//	Enhance Damage Damage Damage.
func TestHystricog(t *testing.T) {
	var healthy, wounded ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Dis,
			InPlay: ct.Cards(Hystricog),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&healthy, ct.Creature(ct.Power(5))),
				ct.Bind(&wounded, ct.Creature(ct.Power(5))),
			),
		},
	})
	wounded.Damaged(2)

	h.P1.UseAction(Hystricog) // only the damaged creature is a legal target

	h.Expect(wounded).At(ct.Discard)
	h.Expect(healthy).At(ct.PlayArea)
}
