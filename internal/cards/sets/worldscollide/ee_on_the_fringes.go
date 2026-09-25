package worldscollide

import "github.com/dmikalova/vex/internal/card"

// E'e on the Fringes
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Imp
//
//	Elusive.
//	After you discard a Dis card, you may purge a Dis card from a discard pile -> steal 1 Æmber.
var EeOnTheFringes = set.New(
	"E'e on the Fringes",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "088"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.AfterDiscardFromHand, card.Conditional{
		Cond: card.ItIs{House: card.Houses.Named(card.House.Self)},
		Then: card.May{
			Do: card.Then{
				First: card.PurgeCard{
					Zones:     []card.Zone{card.Discard},
					Player:    card.ChosenPlayer,
					Selection: card.Chosen{House: card.Houses.Named(card.House.Self)},
				},
				Result: card.StealAember{Amount: 1},
			},
		},
	}),
)
