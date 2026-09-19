package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Deep Probe
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a house. Reveal your opponent's hand. Discard each creature of the chosen house from your opponent's hand.
var DeepProbe = set.New(
	"Deep Probe",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "162"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Sequence{Effects: []card.Effect{
				card.RevealHand{Player: card.Opponent},
				card.DiscardCard{
					Player: card.Opponent,
					Zones:  []card.Zone{card.Hand},
					Selection: card.Each{
						Type:  card.Type.Creature,
						House: card.Houses.Chosen,
					},
				},
			}},
		}),
)
