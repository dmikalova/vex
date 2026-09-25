package ageofascension

import "github.com/dmikalova/vex/internal/card"

// The Flex
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly ready Brobnar creature. Exhaust it. Gain Æmber equal to half its power, rounded down.
var TheFlex = set.New(
	"The Flex",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "31"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.FriendlyCreature.House(card.Houses.Named(card.House.Self)).Ready(),
			Then: card.Sequence{Effects: []card.Effect{
				card.Exhaust{Target: card.Target.Triggering},
				card.GainAember{
					Player:  card.Controller,
					EqualTo: card.PowerOfChosen{Of: card.HalfRoundedDown},
				},
			}},
		}),
)
