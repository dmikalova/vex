package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Chief Engineer Walls
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human
//
//	Elusive.
//	Play/Fight/Reap: You may put an upgrade or Robot card from your discard pile into your hand.
var ChiefEngineerWalls = card.New(
	"Chief Engineer Walls",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "293"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.PlayFightReap, card.May{
		Do: card.PutFromDiscard{
			Type:        card.Type.Upgrade,
			OrTrait:     card.Traits.Robot,
			Destination: card.To.Hand,
		},
	}),
)
