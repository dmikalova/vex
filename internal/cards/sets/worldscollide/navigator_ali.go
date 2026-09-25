package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Navigator Ali
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Look at the top 3 cards of your deck and put them back in any order.
var NavigatorAli = set.New(
	"Navigator Ali",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "314"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.PlayFightReap, card.LookAtTopOfDeck{
		Amount: 3,
		Then:   []card.TopAct{card.ReorderRest{}},
	}),
)
