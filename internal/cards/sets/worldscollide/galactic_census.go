package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Galactic Census
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: If there are 3 or more houses represented among creatures in play, gain 1 Æmber. Gain 1 more if there are 5 or more. Gain 1 more if there are 6 or more.
var GalacticCensus = set.New(
	"Galactic Census",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "332"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Conditional{
				Cond: card.HousesRepresented{
					Among:  card.HousesAmong{Player: card.EachPlayer, Type: card.Type.Creature},
					Is:     card.AtLeast,
					Amount: 3,
				},
				Then: card.GainAember{Player: card.Controller, Amount: 1},
			},
			card.Conditional{
				Cond: card.HousesRepresented{
					Among:  card.HousesAmong{Player: card.EachPlayer, Type: card.Type.Creature},
					Is:     card.AtLeast,
					Amount: 5,
				},
				Then: card.GainAember{Player: card.Controller, Amount: 1},
			},
			card.Conditional{
				Cond: card.HousesRepresented{
					Among:  card.HousesAmong{Player: card.EachPlayer, Type: card.Type.Creature},
					Is:     card.AtLeast,
					Amount: 6,
				},
				Then: card.GainAember{Player: card.Controller, Amount: 1},
			},
		}}),
)
