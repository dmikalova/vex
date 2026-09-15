//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SaurianEgg
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Armor:  5
//	Æmber:  1
//	Traits: Dinosaur • Egg
//
//	Saurian Egg cannot fight or reap.
//	Omni: Discard the top 2 cards of your deck. If you discard any Saurian creatures this way, put them into play ready, give them three +1 power counters, and destroy Saurian Egg.
var SaurianEgg = set.New(
	"Saurian Egg",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "210"),
	card.WithPower(1),
	card.WithArmor(5),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Egg),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
