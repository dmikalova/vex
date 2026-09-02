package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Psychic Network
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each friendly ready Mars creature, steal 1 Æmber.
var PsychicNetwork = card.New(
	"Psychic Network",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 174),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
				Ready:  true,
			},
		}),
)
