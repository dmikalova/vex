package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stomp
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 5 damage to a creature. If this damage destroys that creature, exalt a friendly creature.
var Stomp = card.New(
	"Stomp",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 210),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThenIfDestroyed{
			Amount: 5,
			Target: card.Target.Creature,
			Then: card.Exalt{
				Target: card.Target.FriendlyCreature,
				Amount: 1,
			},
		}),
)
