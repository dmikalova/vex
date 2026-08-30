//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ZyzzixTheMany
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Fight/Reap: You may reveal a creature from your hand. If you do, archive it and Zyzzix the Many gets three +1 power counters.
var ZyzzixTheMany = card.New(
	"Zyzzix the Many",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 207),
	card.WithPower(3),
	card.WithTraits("Martian", "Soldier"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
