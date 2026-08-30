//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// QyxxlyxPlagueMaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Martian • Scientist
//
//	Fight/Reap: Deal 3 Damage to each Human creature. This damage cannot be prevented by armor.
var QyxxlyxPlagueMaster = card.New(
	"Qyxxlyx Plague Master",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 198),
	card.WithPower(3),
	card.WithTraits("Martian", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
