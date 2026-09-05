package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// One Stood Against Many
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Ready and fight with a friendly creature 3 times, each time against a different enemy creature. Resolve these fights one at a time.
var OneStoodAgainstMany = card.New(
	"One Stood Against Many",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 223),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.RepeatedFight{
			Times:  card.Fixed(3),
			Target: card.Target.FriendlyCreature,
		}),
)
