package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Dark Harbinger
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant • Witch
//
//	After you play an Untamed tactic, ready Dark Harbinger.
var DarkHarbinger = set.New(
	"Dark Harbinger",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "381"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Witch),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{
			House: card.Houses.Named(card.House.Self),
			Type:  card.Type.Tactic,
		},
		Then: card.Ready{Target: card.Target.This},
	}),
)
