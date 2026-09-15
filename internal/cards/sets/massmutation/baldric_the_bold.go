//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BaldricTheBold
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Before Fight: If the creature Baldric the Bold fights is the most powerful enemy creature, gain 2A.
var BaldricTheBold = set.New(
	"Baldric the Bold",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "144"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
