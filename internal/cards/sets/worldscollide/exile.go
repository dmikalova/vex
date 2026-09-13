package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Exile
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent gains control of a friendly Creature.
var Exile = card.New(
	"Exile",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "202"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.TakeControl{
			Target:     card.Target.FriendlyCreature,
			ToOpponent: true,
			Duration:   card.Duration.Forever,
		}),
)
