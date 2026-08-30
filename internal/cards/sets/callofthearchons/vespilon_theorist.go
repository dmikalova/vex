//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// VespilonTheorist
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: Choose a house. Reveal the top card of your deck. If it is of that house, archive it and gain 1 Aember. Otherwise, discard it.
var VespilonTheorist = card.New(
	"Vespilon Theorist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 155),
	card.WithPower(2),
	card.WithTraits("Cyborg", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
