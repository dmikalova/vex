package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Ganger Chieftain
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Play: Ready and fight with a neighboring creature.
var MegaGangerChieftain = set.New(
	"Mega Ganger Chieftain",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "56"),
	card.InCluster(card.Pulled(chieftainsBrewCluster, 1, 1.25)),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.OnChooseCreature{
			Target: card.Target.Creature.Neighboring(),
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
