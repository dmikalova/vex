package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bigtwig
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Bigtwig can only fight stunned creatures.
//	Reap: Stun and exhaust a creature.
var Bigtwig = card.New(
	"Bigtwig",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 346),
	card.WithPower(7),
	card.WithTraits(card.Traits.Beast),
	card.WithFightRestriction(card.Stunned),
	card.WithAbility(
		card.Trigger.Reap, card.OnChooseCreature{
			Target: card.Target.Creature,
			Verbs:  []card.CreatureVerb{card.StunVerb{}, card.ExhaustVerb{}},
		}),
)
