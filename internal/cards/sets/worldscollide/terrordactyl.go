package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Terrordactyl
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  12
//	Traits: Beast
//
//	Terrordactyl deals 4 damage when fighting.
//	Terrordactyl enters play stunned.
//	Before Fight: Deal 4 damage to each neighbor of the creature Terrordactyl fights.
var Terrordactyl = set.New(
	"Terrordactyl",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "211"),
	card.WithPower(12),
	card.WithTraits(card.Traits.Beast),
	card.WithEntersPlay(card.Stun{Target: card.Target.This}),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 4,
		Fixed:  true,
	}),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Amount: 4,
			Target: card.Target.CreatureFought.NeighborsOf(),
		}),
)
