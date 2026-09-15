//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// NewFrontiers
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Choose a house. Reveal the top 3 cards of your deck. Archive each card of the chosen house and discard the others.
var NewFrontiers = set.New(
	"New Frontiers",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "326"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
