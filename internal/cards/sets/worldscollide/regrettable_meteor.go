//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// RegrettableMeteor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy each Dinosaur creature and each creature with power 6 or higher.
var RegrettableMeteor = card.New(
	"Regrettable Meteor",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 208),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
