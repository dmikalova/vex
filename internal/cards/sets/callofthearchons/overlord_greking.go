package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Overlord Greking
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	After a creature is destroyed in a fight with Overlord Greking, put it into play under your control.
var OverlordGreking = set.New(
	"Overlord Greking",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "87"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.PutIntoPlay{
			Target:  card.Target.Triggering,
			Control: card.Yours,
		}),
)
