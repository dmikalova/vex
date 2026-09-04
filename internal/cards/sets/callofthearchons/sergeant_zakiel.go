package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sergeant Zakiel
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Play: You may ready and fight with a neighboring creature.
var SergeantZakiel = card.New(
	"Sergeant Zakiel",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 258),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.May{
			Do: card.OnChooseCreature{
				Target: card.Target.Creature.Neighboring(),
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
			},
		}),
)
