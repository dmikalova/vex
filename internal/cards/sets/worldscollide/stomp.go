package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Stomp
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Deal 5 damage to a creature. If this damage destroys that creature, exalt a friendly creature.
var Stomp = set.New(
	"Stomp",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "210"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 5,
			After:  card.IfDestroyed,
			Target: card.Target.Creature,
			Then: card.Exalt{
				Target: card.Target.FriendlyCreature,
				Amount: 1,
			},
		}),
)
