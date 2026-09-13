package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Feathered Shaman
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human • Witch
//
//	Elusive.
//	Fight/Reap: Ward each neighboring Creature.
var TheFeatheredShaman = card.New(
	"The Feathered Shaman",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "383"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.FightReap,
		card.Ward{Target: card.Target.EachCreature.Neighboring()},
	),
)
