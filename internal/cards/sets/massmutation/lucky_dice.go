//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LuckyDice
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Omni: Destroy Lucky Dice. During your opponent's next turn, friendly creatures cannot be dealt damage.
var LuckyDice = set.New(
	"Lucky Dice",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "267"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
