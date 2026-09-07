//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TechivorePulpate
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Jelly
//
//	After a player chooses an active house, destroy each artifact of that house.
var TechivorePulpate = card.New(
	"Techivore Pulpate",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 341),
	card.WithPower(5),
	card.WithTraits(card.Traits.Jelly),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
