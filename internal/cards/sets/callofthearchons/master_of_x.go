package callofthearchons

import (
	"fmt"

	"github.com/dmikalova/vactrol/internal/card"
)

// master builds one Master of N: a power-4 Demon whose Reap may destroy a
// creature of power n. Each variant registers as its own pool card at one-fifth
// draft weight, so all five together draft about as often as one ordinary Rare
// card. KeyForge prints Master of 1/2/3 (#89/#90/#91); Vactrol extends the family
// to 5. The name is built from n; mage generateComments resolves it through the
// wrapper to document each variant. opts carries a variant's provenance; the home
// set is stamped by set.New, so a variant with no provenance passes no opts.
func master(n int, opts ...card.Option) card.Definition {
	return set.New(fmt.Sprintf("Master of %d", n),
		card.House.Dis, card.Type.Creature, card.Rarity.Rare,
		append(opts,
			card.WithPower(4),
			card.WithTraits(card.Traits.Demon),
			card.WithAbility(card.Trigger.Reap, card.May{
				Do: card.Destroy{Target: card.Target.Creature.PowerExactly(n)},
			}),
			card.RarityWeight(0.2),
		)...)
}

// Master of 1
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a Creature with power 1.
var MasterOf1 = master(1, card.Provenance(card.CotA, "89"))

// Master of 2
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a Creature with power 2.
var MasterOf2 = master(2, card.Provenance(card.CotA, "90"))

// Master of 3
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a Creature with power 3.
var MasterOf3 = master(3, card.Provenance(card.CotA, "91"))

// Master of 4
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a Creature with power 4.
var MasterOf4 = master(4)

// Master of 5
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a Creature with power 5.
var MasterOf5 = master(5)
