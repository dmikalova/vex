package massmutation

import "github.com/dmikalova/vex/internal/card"

// J.O.N. Cargo
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	Reap: Discard the top card of your deck. Reveal your hand. Archive each card of that card's house from your hand.
var JONCargo = set.New(
	"J.O.N. Cargo",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "347"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Reap, card.Sequence{Effects: []card.Effect{
			card.DiscardTop{
				Amount: 1,
				Player: card.Controller,
			},
			card.RevealHand{Player: card.Controller},
			card.ArchiveCard{
				Zone:      card.Hand,
				Selection: card.Each{House: card.Houses.Contextual},
			},
		}},
	),
)
