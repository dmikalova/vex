package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Exile
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Your opponent gains control of a friendly creature.
var Exile = set.New(
	"Exile",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "202"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.TakeControl{
			Target:     card.Target.FriendlyCreature,
			ToOpponent: true,
			Duration:   card.Duration.UntilCardLeavesPlay,
		}),
)
