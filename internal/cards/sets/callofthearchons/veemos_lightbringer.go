package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Veemos Lightbringer
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Angel • Spirit
//
//	Play: Destroy each elusive creature.
var VeemosLightbringer = card.New(
	"Veemos Lightbringer",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 262),
	card.WithPower(6),
	card.WithTraits("Angel", "Spirit"),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Keyword(card.Keyword.Elusive),
		}),
)
