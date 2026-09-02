package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Smiling Ruth
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: If you forged a key this turn, take control of an enemy flank creature.
var SmilingRuth = card.New(
	"Smiling Ruth",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 312),
	card.WithPower(1),
	card.WithTraits("Elf", "Thief"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.Conditional{
			Cond: card.ForgedKey{Player: card.Controller},
			Then: card.TakeControl{
				Target:   card.Target.EnemyCreature.OnFlank(),
				Duration: card.Duration.Forever,
			},
		}),
)
