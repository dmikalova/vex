package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Skoll
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Assault 3.
//	After a creature is destroyed by Skoll's assault damage, give a friendly creature a +1 power counter.
var Skoll = set.New(
	"Skoll",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "29"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithAssault(3),
	card.WithAbility(
		card.Trigger.AfterAssaultDestroys, card.AddPowerCounter{
			Target: card.Target.FriendlyCreature,
			Amount: 1,
		}),
)
