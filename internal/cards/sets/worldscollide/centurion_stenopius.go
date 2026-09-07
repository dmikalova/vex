//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CenturionStenopius
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Armor:  2
//	Traits: Dinosaur • Soldier
//
//	Centurion Stenopius gets +3 power for each A on it.
//	Play/Fight/Reap: You may exalt Centurion Stenopius.
var CenturionStenopius = card.New(
	"Centurion Stenopius",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 199),
	card.WithPower(3),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
