package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Sting
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Vehicle
//
//	You skip your "forge a key" step.
//	You gain all Æmber your opponent spends when forging a key.
//	Action: Destroy The Sting.
var TheSting = card.New(
	"The Sting",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 295),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Vehicle),
	card.WithRestrictions(card.Restrictions{SkipForge: true}),
	card.WithGainsForgeAember(),
	card.WithAbility(card.Trigger.Action, card.Destroy{Target: card.Target.This}),
)
