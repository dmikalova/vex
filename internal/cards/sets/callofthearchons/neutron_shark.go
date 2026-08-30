//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// NeutronShark
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Mutant
//
//	Play/Fight/Reap: Destroy an enemy creature or artifact and a friendly creature or artifact. Discard the top card of your deck. If that card is not a Logos card, trigger this effect again.
var NeutronShark = card.New(
	"Neutron Shark",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 146),
	card.WithPower(1),
	card.WithTraits("Beast", "Mutant"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
