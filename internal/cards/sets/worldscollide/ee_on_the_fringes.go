package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// EeOnTheFringes
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Imp
//
//	Elusive.
//	During your turn, after you discard a Dis card from your hand, you may purge a Dis card from a discard pile. If you do, steal 1A.
var EeOnTheFringes = card.New(
	"E'e on the Fringes",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "088"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.AfterDiscardFromHand, card.Conditional{
		Cond: card.ItIs{House: card.House.Self},
		Then: card.May{
			Do: card.Then{
				First:  card.PurgeCard{Zone: card.Discard, House: card.House.Self},
				Result: card.StealAember{Amount: 1},
			},
		},
	}),
)
