package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Hapless Cadet
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Alien
//
//	Taunt.
//	Destroyed: Your opponent loses 3 Æmber.
var HaplessCadet = set.New(
	"Hapless Cadet",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "345"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Alien),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Destroyed, card.LoseAember{
			Player: card.Opponent,
			Amount: 3,
		}),
)
