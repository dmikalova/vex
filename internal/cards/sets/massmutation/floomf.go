package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Floomf
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast • Cat
//
//	Skirmish.
//	Fight: Give a Beast creature two +1 power counters.
var Floomf = set.New(
	"Floomf",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "366"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast, card.Traits.Cat),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.AddPowerCounter{
			Amount: 2,
			Target: card.Target.Creature.WithTrait(card.Traits.Beast),
		}),
)
