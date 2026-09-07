package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Navigator Ali
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Look at the top 3 cards of your deck and put them back in any order.
var NavigatorAli = card.New(
	"Navigator Ali",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "314"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithPlayFightReap(card.ReorderTop{Amount: 3}),
)
