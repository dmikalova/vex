//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// AngryMob
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Before Fight: You may discard cards from the top of your deck until you discard an Angry Mob or run out of cards. If you discard an Angry Mob this way, put it into your hand.
var AngryMob = set.New(
	"Angry Mob",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "143"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
