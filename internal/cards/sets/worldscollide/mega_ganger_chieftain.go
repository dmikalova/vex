package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Play: You may ready and fight with a neighboring creature.
var MegaGangerChieftain = card.New(
	"Mega Ganger Chieftain",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "56"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.May{
			Do: card.OnChooseCreature{
				Target: card.Target.Creature.Neighboring(),
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
			},
		}),
)
