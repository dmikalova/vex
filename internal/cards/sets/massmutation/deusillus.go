package massmutation

import "github.com/dmikalova/vex/internal/card"

// Deusillus
//
//	House:  Saurian
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  20
//	Traits: Mutant
//
//	Play: Deusillus captures all your opponent's Æmber. Deal 5 damage to an enemy creature.
//	Fight/Reap: Move 1 Æmber from Deusillus to the common supply. Deal 2 damage to each enemy creature.
var Deusillus = set.Gigantic(
	"Deusillus",
	card.House.Saurian,
	card.Rarity.Rare,
	card.Provenance(card.MM, "244"),
	card.WithPower(20),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.CaptureAember{
				All:    true,
				Source: card.Opponent,
				Target: card.Target.This,
			},
			card.DealDamage{
				Amount: 5,
				Target: card.Target.EnemyCreature,
			},
		}}),
	card.WithAbility(
		card.Trigger.FightReap, card.Sequence{Effects: []card.Effect{
			card.MoveAemberToSupply{
				Amount: 1,
				Target: card.Target.This,
			},
			card.DealDamage{
				Amount: 2,
				Target: card.Target.EachEnemyCreature,
			},
		}}),
)
