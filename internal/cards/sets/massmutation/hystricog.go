package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Hystricog
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant
//
//	Action: Destroy a damaged creature.
//	Enhance Damage Damage Damage.
var Hystricog = set.New(
	"Hystricog",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "024"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithEnhance(card.Bonus.Damage, card.Bonus.Damage, card.Bonus.Damage),
	card.WithAbility(
		card.Trigger.Action, card.Destroy{Target: card.Target.Creature.Damaged()}),
)
