package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Shattered Throne
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	After a creature is used to fight, it captures 1 Æmber from its opponent.
var ShatteredThrone = set.New(
	"Shattered Throne",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "28"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.AfterCreatureFights, card.CaptureAember{
			Amount: 1,
			Target: card.Target.Triggering,
			Source: card.ItsOpponent,
		}),
)
