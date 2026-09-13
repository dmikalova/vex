package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Shattered Throne
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	After a Creature is used to fight, it captures 1 Æmber from its opponent.
var ShatteredThrone = card.New(
	"Shattered Throne",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "28"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.AfterCreatureFights, card.CaptureAember{
			Amount: 1,
			Target: card.Target.Triggering,
			Source: card.ItsOpponent,
		}),
)
