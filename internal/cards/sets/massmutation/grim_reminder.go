package massmutation

import "github.com/dmikalova/vex/internal/card"

// Grim Reminder
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a house. Archive each creature of the chosen house from your discard pile. Gain 1 chain.
var GrimReminder = set.New(
	"Grim Reminder",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "023"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Sequence{
				Effects: []card.Effect{
					card.ArchiveCard{
						Zone: card.Discard,
						Selection: card.Each{
							Type:  card.Type.Creature,
							House: card.Houses.Chosen,
						},
					},
					card.GainChains{Amount: 1},
				},
			},
		}),
)
