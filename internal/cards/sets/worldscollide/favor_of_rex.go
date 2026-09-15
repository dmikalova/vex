package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Favor of Rex
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Trigger the play effect of a creature.
var FavorOfRex = set.New(
	"Favor of Rex",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "219"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.TriggerAbility{
			Trigger: card.Trigger.Play,
			Target:  card.Target.Creature,
		}),
)
