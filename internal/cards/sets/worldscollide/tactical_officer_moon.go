package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tactical Officer Moon
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Assault 2.
//	Play: You may rearrange the Creatures in a player's battleline.
var TacticalOfficerMoon = set.New(
	"Tactical Officer Moon",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "320"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithAssault(2),
	card.WithAbility(card.Trigger.Play, card.May{Do: card.RearrangeBattleline{}}),
)
