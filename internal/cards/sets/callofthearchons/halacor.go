package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Halacor
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Beast
//
//	Each friendly flank creature gains skirmish.
var Halacor = card.New(
	"Halacor",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 355),
	card.WithPower(4),
	card.WithTraits("Beast"),
	card.WithConstantAbility(card.ConstantAbility{
		Keywords: card.Keywords(card.Keyword.Skirmish),
		Target:   card.Target.EachFriendlyCreature.OnFlank(),
	}),
)
