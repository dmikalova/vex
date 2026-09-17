package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Cera Severa
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	Fight/Reap: Cera Severa captures 1 Æmber from your opponent.
//	Destroyed: For each Æmber on Cera Severa, deal 1 damage to an enemy creature.
var CeraSevera = set.New(
	"Cera Severa",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "232"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.FightReap, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Destroyed, card.DealDamage{
			Amount: 1,
			Per:    card.AemberOnThis{},
			Target: card.Target.EnemyCreature,
		}),
)
