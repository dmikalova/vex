package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Harvest Time
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a trait, then purge each card with that trait. Each player gains 1 Æmber for each card they controlled that was purged this way.
var HarvestTime = card.New(
	"Harvest Time",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "106"),
	card.WithAbility(card.Trigger.Play, card.PurgeEachOfChosenTrait{}),
)
