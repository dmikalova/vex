package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Martian Hounds
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: For each damaged creature in play, give a creature 2 +1 power counters.
var MartianHounds = card.New(
	"Martian Hounds",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 167),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.Creature,
			Amount: 2,
			Per: card.InPlay{
				Player:  card.EachPlayer,
				Type:    card.Type.Creature,
				Damaged: true,
			},
		}),
)
