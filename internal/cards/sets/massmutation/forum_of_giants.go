//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ForumOfGiants
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	At the start of your turn, the player who controls the most powerful creature gains 1A.
var ForumOfGiants = set.New(
	"Forum of Giants",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "219"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
