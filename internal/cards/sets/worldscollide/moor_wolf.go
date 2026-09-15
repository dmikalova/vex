package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Moor Wolf
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast • Wolf
//
//	Skirmish.
//	Play: Ready each other friendly Wolf creature.
var MoorWolf = set.New(
	"Moor Wolf",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "393"),
	card.InCluster(card.Cluster{
		Name:     "Moor Wolf",
		Strategy: card.ClusterStrategy.SelfPull,
		Trigger:  card.ClusterTrigger.ByAnyMember,
		Min:      3,
		Mean:     5,
	}),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Wolf),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Play, card.Ready{
			Target: card.Target.EachOtherFriendlyCreature.WithTrait(card.Traits.Wolf),
		}),
)
