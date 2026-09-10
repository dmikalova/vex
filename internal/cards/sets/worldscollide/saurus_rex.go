package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Saurus Rex
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Dinosaur • Leader
//
//	Fight/Reap: If Saurus Rex is in the center of your battleline, you may exalt Saurus Rex -> search your deck for a Saurian card, reveal it, and put it into your hand. Shuffle your deck.
var SaurusRex = card.New(
	"Saurus Rex",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "227"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Leader),
	card.WithFightOrReap(card.Conditional{
		Cond: card.SourceInCenterOfBattleline{},
		Then: card.May{Do: card.Then{
			First: card.Exalt{Target: card.Target.This, Amount: 1},
			Result: card.Sentences{
				Effects: []card.Effect{
					card.SearchDeck{House: card.House.Self},
					card.ShuffleDeck{},
				},
			},
		}},
	}),
)
