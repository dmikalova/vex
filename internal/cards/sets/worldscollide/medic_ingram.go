package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Medic Ingram
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Choose a creature. Heal 3 damage from it. Ward it.
var MedicIngram = set.New(
	"Medic Ingram",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "301"),
	card.InCluster(card.Pulled(ingramsBlasterCluster, 1, 1.25)),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.PlayFightReap, card.ChooseCreatureThen{
		Target: card.Target.Creature,
		Then: card.Sequence{Effects: []card.Effect{
			card.Heal{
				Amount: 3,
				Target: card.Target.Triggering,
			},
			card.Ward{Target: card.Target.Triggering},
		}},
	}),
)
