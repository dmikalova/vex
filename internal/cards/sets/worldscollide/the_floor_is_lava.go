package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Floor is Lava
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	At the start of your turn, deal 1 damage to a friendly Creature, and deal 1 damage to an enemy Creature.
var TheFloorIsLava = set.New(
	"The Floor is Lava",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "51"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.StartOfTurn, card.Sequence{Effects: []card.Effect{
			card.DealDamage{
				Target: card.Target.FriendlyCreature,
				Amount: 1,
			},
			card.DealDamage{
				Target: card.Target.EnemyCreature,
				Amount: 1,
			},
		}}),
)
