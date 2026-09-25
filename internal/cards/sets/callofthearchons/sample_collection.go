package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Sample Collection
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each forged key your opponent has, put an enemy creature into your archives.
var SampleCollection = set.New(
	"Sample Collection",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "175"),
	card.WithAbility(
		card.Trigger.Play, card.ForEach{
			Times: card.ForgedKeys{Player: card.Opponent},
			Do: card.PutFromPlay{
				Target:      card.Target.EnemyCreature,
				Destination: card.To.Archives.Yours(),
			},
		}),
)
