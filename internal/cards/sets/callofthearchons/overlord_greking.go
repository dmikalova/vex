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
//	After a creature is destroyed fighting Overlord Greking, put it into play under your control.
var OverlordGreking = card.New(
	"Overlord Greking",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 87),
	card.WithPower(7),
	card.WithTraits("Demon"),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.PutIntoPlay{
			Target:           card.Target.Triggering,
			UnderYourControl: true,
		}),
)
