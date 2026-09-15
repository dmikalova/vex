//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// CyberClone
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Play: Purge another creature. Until Cyber-Clone leaves play, it has power equal to the purged creature's power, and gains that creature's armor, keywords, and traits.
var CyberClone = set.New(
	"Cyber-Clone",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "102"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
