package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Mega Groke
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Fight: Your opponent loses 1 Æmber.
var MegaGroke = set.New(
	"Mega Groke",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "57"),
	card.InCluster(card.Pulled(grokesBrewCluster, 1, 1.25)),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Fight, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
