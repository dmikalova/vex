//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MushroomMan
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Fungus • Human
//
//	Mushroom Man gets +3 power for each unforged key you have.
var MushroomMan = card.New(
	"Mushroom Man",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 362),
	card.WithPower(2),
	card.WithTraits("Fungus", "Human"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
