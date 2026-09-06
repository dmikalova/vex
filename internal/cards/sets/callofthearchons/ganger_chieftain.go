package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Play: Ready and fight with a neighboring creature.
var GangerChieftain = card.New(
	"Ganger Chieftain",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 33),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.Creature.Neighboring(),
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
