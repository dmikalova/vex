package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Groke
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Fight: Your opponent loses 1 Æmber.
var MegaGroke = card.New(
	"Mega Groke",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "57"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Fight, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
