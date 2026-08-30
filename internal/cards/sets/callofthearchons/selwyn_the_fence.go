//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SelwynTheFence
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Thief
//
//	Fight/Reap: Move 1 Aember from one of your cards to your pool.
var SelwynTheFence = card.New(
	"Selwyn the Fence",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 309),
	card.WithPower(3),
	card.WithTraits("Elf", "Thief"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
