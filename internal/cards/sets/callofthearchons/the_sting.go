package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Sting
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Vehicle
//
//	You skip your "forge a key" phase.
//	You gain all Æmber your opponent spends when forging a key.
//	Action: Destroy The Sting.
var TheSting = set.New(
	"The Sting",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "295"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Vehicle),
	card.WithRestrictions(card.Restrictions{SkipForge: true}),
	card.WithGainsForgeAember(),
	card.WithAbility(card.Trigger.Action, card.Destroy{Target: card.Target.This}),
)
