package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Snag
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Fight: Your opponent must choose the house of the Creature Snag fights as their active house on their next turn.
var Snag = set.New(
	"Snag",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "096"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Fight,
		card.OpponentMustChooseHouse{Source: card.FoughtActiveHouse},
	),
)
