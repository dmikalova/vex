package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// The Flex
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly ready Brobnar creature - exhaust it, and gain Æmber equal to half its power, rounded down.
var TheFlex = card.New(
	"The Flex",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "31"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.FriendlyCreature.OfHouse(card.House.Self).Ready(),
			Then: card.Sequence{Effects: []card.Effect{
				card.Exhaust{Target: card.Target.Triggering},
				card.GainAemberEqualTo{
					Player: card.Controller,
					Count:  card.PowerOfChosen{Of: card.Half},
				},
			}},
		}),
)
