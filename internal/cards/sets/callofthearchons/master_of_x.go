package callofthearchons

import (
	"fmt"
	"math/rand"

	"github.com/dmikalova/vactrol/internal/card"
)

// masterOf builds a Master of X variant — a power-4 Demon whose Reap may destroy
// a creature of a chosen power. n == 0 is the template face, whose power renders
// "X"; n in 1..5 fixes a concrete variant's power.
func masterOf(n int) []card.Option {
	power := card.Target.Creature.PowerVariable()
	if n > 0 {
		power = card.Target.Creature.PowerExactly(n)
	}
	return []card.Option{
		card.WithPower(4),
		card.WithTraits("Demon"),
		card.WithAbility(card.Trigger.Reap, card.May{Do: card.Destroy{Target: power}}),
	}
}

// Master of X
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power X.
var MasterOfX = card.New(
	"Master of X",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	append(masterOf(0),
		card.Provenance(card.CotA, 89),
		card.Provenance(card.CotA, 90),
		card.Provenance(card.CotA, 91),
		card.Template(func(ctx card.SlotContext, r *rand.Rand) card.Definition {
			n := r.Intn(5) + 1
			return card.Build(
				fmt.Sprintf("Master of %d", n),
				ctx.House,
				card.Type.Creature,
				card.Rarity.Rare,
				masterOf(n)...,
			)
		}),
	)...,
)
