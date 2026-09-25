package worldscollide

import "github.com/dmikalova/vex/internal/card"

// The Floor is Lava
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	At the start of your turn, deal 1 damage to a friendly creature. Deal 1 damage to an enemy creature.
var TheFloorIsLava = set.New(
	"The Floor is Lava",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "51"),
	card.WithBonus(card.Bonus.Aember),
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
