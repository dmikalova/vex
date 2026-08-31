package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Harland Mindlock
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Cyborg • Scientist
//
//	Play: Take control of an enemy flank creature until Harland Mindlock leaves play.
var HarlandMindlock = card.New(
	"Harland Mindlock",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 143),
	card.WithPower(1),
	card.WithTraits("Cyborg", "Scientist"),
	card.WithAbility(card.Trigger.Play, card.TakeControl{
		Target:   card.Target.EnemyCreature.OnFlank(),
		Duration: card.Duration.UntilThisLeavesPlay,
	}),
)
