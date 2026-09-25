package worldscollide

import "github.com/dmikalova/vex/internal/card"

// General Order 24
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	At the start of each player's turn, if there are no friendly creatures in play, destroy General Order 24. Otherwise, choose a friendly creature. Destroy each creature of that card's house.
var GeneralOrder24 = set.New(
	"General Order 24",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "333"),
	card.WithTraits(card.Traits.Law),
	card.WithEachPlayerAbility(
		card.Trigger.StartOfTurn, card.Conditional{
			Cond: card.CardsInPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				None:   true,
			},
			Then: card.Destroy{Target: card.Target.This},
			Else: card.ChooseCreatureThen{
				Target: card.Target.FriendlyCreature,
				Then: card.Destroy{
					Target: card.Target.EachCreature.House(card.Houses.Contextual),
				},
			},
		}),
)
