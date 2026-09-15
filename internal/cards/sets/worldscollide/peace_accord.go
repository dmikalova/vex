package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Peace Accord
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	Play: Each player gains 2 Æmber.
//	After a creature is used to fight, its controller loses 4 Æmber. Destroy Peace Accord.
var PeaceAccord = set.New(
	"Peace Accord",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "335"),
	card.WithTraits(card.Traits.Law),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.EachPlayer,
			Amount: 2,
		}),
	card.WithAbility(
		card.Trigger.AfterCreatureFights, card.Sentences{Effects: []card.Effect{
			card.LoseAember{
				Player: card.ItsOwner,
				Amount: 4,
			},
			card.Destroy{Target: card.Target.This},
		}}),
)
